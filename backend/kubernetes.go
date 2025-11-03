package backend

import (
	"bytes"
	"context"
	"errors"

	"github.com/Orzelius/cosi-testing/log"
	k8sdiff "k8s.io/kubectl/pkg/cmd/diff"

	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/cli-utils/pkg/apply"
	"sigs.k8s.io/cli-utils/pkg/apply/event"
	"sigs.k8s.io/cli-utils/pkg/common"
	"sigs.k8s.io/cli-utils/pkg/inventory"
	"sigs.k8s.io/cli-utils/pkg/manifestreader"
	"sigs.k8s.io/cli-utils/pkg/object/validation"
)

var _ Backend = &Kubernetes{}

type Kubernetes struct {
	applier      *apply.Applier
	mapper       meta.RESTMapper
	inventory    inventory.Inventory
	l            *log.MyLogger
	factory      util.Factory
	fieldManager string
}

func (b *Kubernetes) Init() error {
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

	inventory, err := inventoryClient.NewInventory(inventory.NewSingleObjectInfo("inventory", types.NamespacedName{Namespace: "inventory", Name: "inventory"}))
	if err != nil {
		return err
	}

	b.inventory = inventory
	b.l = log.GetLogger()

	return nil
}

func (b *Kubernetes) Apply(ctx context.Context, data []byte) error {
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
		return err
	}

	b.l.Info("starting apply proccess")

	eventCh := b.applier.Run(ctx, b.inventory.Info(), inputResources, apply.ApplierOptions{
		// DryRunStrategy: common.DryRunServer,
		ServerSideOptions: common.ServerSideOptions{
			ServerSideApply: true,
			FieldManager:    b.fieldManager,

			// ForceConflicts overwrites the fields when applying if the field manager differs.
			ForceConflicts: true,
		},
		EmitStatusEvents:       true,
		InventoryPolicy:        inventory.PolicyAdoptIfNoInventory,
		ValidationPolicy:       validation.ExitEarly,
		PrunePropagationPolicy: v1.DeletePropagationBackground,
	})

	for {
		select {
		case <-ctx.Done():
			return nil
		case e, ok := <-eventCh:
			if e.Type == event.ErrorType {
				b.l.Error(e.String())
			} else {
				b.l.Debug(e.String())
			}
			if !ok {
				b.l.Info("finished applying")
				return nil
			}
		}
	}
}

func (b *Kubernetes) Diff(ctx context.Context, data []byte) error {
	ops := k8sdiff.DiffOptions{
		FieldManager:    b.fieldManager,
		Builder:         b.factory.NewBuilder(),
		ForceConflicts:  false,
		ServerSideApply: true, // might not work
		OpenAPIGetter:   b.factory,
	}

	dynamicClient, err := b.factory.DynamicClient()
	if err != nil {
		return err
	}
	ops.DynamicClient = dynamicClient

	ops.CmdNamespace, ops.EnforceNamespace, err = b.factory.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	err = ops.Run()
	if err != nil {
		return err
	}

	// give up: the code is too CLI specific and a pain to work around

	return errors.New("diffing is not supported with kubernetes backend")
}
