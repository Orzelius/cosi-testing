package backend

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Orzelius/cosi-testing/log"
	"github.com/fluxcd/cli-utils/pkg/kstatus/polling"
	"github.com/fluxcd/pkg/ssa"
	"github.com/fluxcd/pkg/ssa/utils"
	"github.com/gonvenience/ytbx"
	"github.com/homeport/dyff/pkg/dyff"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

var _ Backend = &FluxSSA{}

type FluxSSA struct {
	resourceManager *ssa.ResourceManager
	l               *log.MyLogger
}

func (b *FluxSSA) Init(ctx context.Context) error {
	flags := genericclioptions.NewConfigFlags(true)
	kubeconfig, err := flags.ToRESTConfig()
	if err != nil {
		return err
	}
	restMapper, err := flags.ToRESTMapper()
	if err != nil {
		return err
	}

	kubeClient, err := client.New(kubeconfig, client.Options{
		Mapper: restMapper,
	})
	if err != nil {
		return err
	}

	poller := polling.NewStatusPoller(kubeClient, restMapper, polling.Options{})

	b.resourceManager = ssa.NewResourceManager(kubeClient, poller, ssa.Owner{
		Field: "resource-manager",
		Group: "resource-manager.io",
	})

	b.l = log.GetLogger()

	return nil
}

func (b *FluxSSA) Apply(ctx context.Context, data []byte, dryRun bool) error {
	if dryRun {
		return errors.New("ssa does not supoprt dry run")
	}

	objects, err := utils.ReadObjects(strings.NewReader(string(data)))
	if err != nil {
		return err
	}

	b.l.Infof("starting apply process of %d item(s), dry-run=%t", len(objects), dryRun)

	changeSet, err := b.resourceManager.ApplyAll(ctx, objects, ssa.ApplyOptions{
		Force:             false,
		ExclusionSelector: nil,
		WaitInterval:      2 * time.Second,
		WaitTimeout:       60 * time.Second,
	})
	if err != nil {
		return err
	}

	for _, e := range changeSet.Entries {
		b.l.Debugf("%s %s %s", e.Action.String(), e.ObjMetadata.GroupKind, e.ObjMetadata.Name)
	}

	b.l.Infof("apply finished")

	return nil

}

func (b *FluxSSA) Diff(ctx context.Context, data []byte) error {
	return nil
}

func getHumanReadableDiff(inclusterObj *unstructured.Unstructured, inputObj *unstructured.Unstructured) (string, error) {
	// use dyff to diff in-cluster vs input object in a temporary directory
	tmpDir, err := os.MkdirTemp("", "dyff-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	liveYAML, _ := yaml.Marshal(inclusterObj)
	liveFile := filepath.Join(tmpDir, "live.yaml")
	if err := os.WriteFile(liveFile, liveYAML, 0o644); err != nil {
		return "", err
	}

	mergedYAML, _ := yaml.Marshal(inputObj)
	mergedFile := filepath.Join(tmpDir, "merged.yaml")
	if err := os.WriteFile(mergedFile, mergedYAML, 0o644); err != nil {
		return "", err
	}

	from, to, err := ytbx.LoadFiles(liveFile, mergedFile)
	if err != nil {
		return "", err
	}

	report, err := dyff.CompareInputFiles(from, to,
		dyff.IgnoreOrderChanges(false),
		dyff.KubernetesEntityDetection(true),
	)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	reportWriter := &dyff.HumanReport{Report: report, OmitHeader: true}
	if err := reportWriter.WriteReport(&buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
