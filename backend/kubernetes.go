package backend

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/Orzelius/cosi-testing/log"
	"github.com/fluxcd/cli-utils/pkg/kstatus/polling"
	"github.com/fluxcd/pkg/ssa"
	"github.com/fluxcd/pkg/ssa/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/cli-utils/pkg/apply"
	"sigs.k8s.io/cli-utils/pkg/apply/event"
	"sigs.k8s.io/cli-utils/pkg/apply/prune"
	"sigs.k8s.io/cli-utils/pkg/common"
	"sigs.k8s.io/cli-utils/pkg/inventory"
	"sigs.k8s.io/cli-utils/pkg/manifestreader"
	"sigs.k8s.io/cli-utils/pkg/object/validation"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ Backend = &Kubernetes{}

type Kubernetes struct {
	applier         *apply.Applier
	mapper          meta.RESTMapper
	inventoryClient inventory.Client
	inventoryInfo   *inventory.SingleObjectInfo
	log             *log.MyLogger
	factory         util.Factory
	fieldManager    string
	resourceManager *ssa.ResourceManager
}

func (b *Kubernetes) Init(ctx context.Context) error {
	b.log = log.GetLogger()

	flags := genericclioptions.NewConfigFlags(true)
	kubeconfig, err := flags.ToRESTConfig()
	if err != nil {
		return err
	}
	mapper, err := flags.ToRESTMapper()
	if err != nil {
		return err
	}
	b.mapper = mapper

	factory := util.NewFactory(flags)
	b.factory = factory

	inventoryClient, err := inventory.ConfigMapClientFactory{StatusEnabled: true}.NewClient(factory)
	if err != nil {
		return err
	}
	b.inventoryClient = inventoryClient

	builder := apply.NewApplierBuilder()

	builder.WithFactory(factory)
	builder.WithInventoryClient(inventoryClient)
	builder.WithRestMapper(mapper)
	builder.WithRestConfig(kubeconfig)

	applier, err := builder.Build()
	if err != nil {
		return err
	}
	b.applier = applier
	b.fieldManager = "my-test-manager"

	kubeClient, err := client.New(kubeconfig, client.Options{
		Mapper: mapper,
	})
	if err != nil {
		return err
	}

	poller := polling.NewStatusPoller(kubeClient, mapper, polling.Options{})

	b.resourceManager = ssa.NewResourceManager(kubeClient, poller, ssa.Owner{
		Field: "resource-manager",
		Group: "resource-manager.io",
	})

	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return err
	}

	inventoryInfo := inventory.NewSingleObjectInfo("inventory", types.NamespacedName{Namespace: "inventory", Name: "inventory"})
	b.inventoryInfo = inventoryInfo

	_, err = clientset.CoreV1().Namespaces().Get(ctx, "inventory", metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			b.log.Infof("creating 'inventory' namespace")
			ns := &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "inventory",
				},
			}
			_, err = clientset.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	_, err = inventoryClient.Get(ctx, inventoryInfo, inventory.GetOptions{})
	if apierrors.IsNotFound(err) {
		b.log.Infof("creating the in cluster inventory (configmap)")
		inv, err := inventoryClient.NewInventory(inventoryInfo)
		if err != nil {
			return err
		}
		err = inventoryClient.CreateOrUpdate(ctx, inv, inventory.UpdateOptions{})
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func (b *Kubernetes) Apply(ctx context.Context, data []byte, dryRun bool) error {
	inputResources, err := getInputResources(data, b)
	if err != nil {
		return err
	}

	b.log.Infof("starting apply process of %d item(s), dry-run=%t", len(inputResources), dryRun)

	applyOps := apply.ApplierOptions{
		ServerSideOptions: common.ServerSideOptions{
			ServerSideApply: true,
			FieldManager:    b.fieldManager,

			// ForceConflicts overwrites the fields when applying if the field manager differs.
			ForceConflicts: true,
		},
		ReconcileTimeout:       15 * time.Second,
		EmitStatusEvents:       true,
		InventoryPolicy:        inventory.PolicyAdoptIfNoInventory,
		ValidationPolicy:       validation.ExitEarly,
		PrunePropagationPolicy: metav1.DeletePropagationBackground,
	}

	if dryRun {
		applyOps.DryRunStrategy = common.DryRunServer
	}
	eventCh := b.applier.Run(ctx, b.inventoryInfo, inputResources, applyOps)

	for {
		select {
		case <-ctx.Done():
			return nil
		case e, ok := <-eventCh:
			if e.Type == event.ErrorType {
				b.log.Error(e.String())
			} else {
				b.log.Debug(e.String())
			}
			if e.WaitEvent.Status == event.ReconcileTimeout {
				b.log.Warnf("Reconcile of %s timed out", e.WaitEvent.Identifier)
			}
			if e.WaitEvent.Status == event.ReconcileFailed {
				b.log.Warnf("Reconcile of %s failed", e.WaitEvent.Identifier)
			}
			if !ok {
				b.log.Info("finished applying")
				return nil
			}
		}
	}
}

func (b *Kubernetes) Diff(ctx context.Context, data []byte) error {
	// Read using ssa utils for the diff as reading with kubernetes logic adds unnecessary junk that causes diff with the cluster
	inputResources, err := utils.ReadObjects(strings.NewReader(string(data)))
	if err != nil {
		return err
	}

	pruner, err := prune.NewPruner(b.factory, b.inventoryClient)
	if err != nil {
		return err
	}

	inventory, err := b.inventoryClient.Get(ctx, b.inventoryInfo, inventory.GetOptions{})
	if err != nil {
		return err
	}

	pruneObjs, err := pruner.GetPruneObjs(ctx, inventory, inputResources, prune.Options{})
	if err != nil {
		return err
	}
	for _, obj := range pruneObjs {
		b.log.Infof("%s/%s-%s to delete", obj.GetNamespace(), obj.GetKind(), obj.GetName())
	}

	for _, obj := range inputResources {
		changeSet, inclusterObj, inputObj, err := b.resourceManager.Diff(ctx, obj, ssa.DiffOptions{})
		if (err != nil && apierrors.IsNotFound(err)) || (err == nil && changeSet.Action == ssa.CreatedAction) {
			b.log.Infof("%s/%s-%s to create", obj.GetNamespace(), obj.GetKind(), obj.GetName())
			continue
		}
		if err != nil {
			return err
		}

		extraDiffDetails := ""
		if changeSet.Action == ssa.ConfiguredAction {
			extraDiffDetails, err = getHumanReadableDiff(inclusterObj, inputObj)
			if err != nil {
				return err
			}
		}

		b.log.Infof("%s %s %s%s", changeSet.Action.String(), changeSet.ObjMetadata.GroupKind, changeSet.ObjMetadata.Name, extraDiffDetails)
	}

	return nil
}

func getInputResources(data []byte, b *Kubernetes) ([]*unstructured.Unstructured, error) {
	reader := &manifestreader.StreamManifestReader{
		ReaderName: "inmem",
		Reader:     bytes.NewReader(data),
		ReaderOptions: manifestreader.ReaderOptions{
			Validate:  true, // set true if you wire a RESTMapper
			Namespace: "",   // optional default namespace
			Mapper:    b.mapper,
		},
	}

	inputResources, err := reader.Read()
	if err != nil {
		return nil, err
	}
	return inputResources, nil
}
