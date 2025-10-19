package controllers

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/Orzelius/cosi-testing/constants"
	"github.com/Orzelius/cosi-testing/myresource"
	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/cli-utils/pkg/manifestreader"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

type KubernetesInputController struct{}

// Name implements controller.Controller interface.
func (ctrl *KubernetesInputController) Name() string {
	return "KubernetesInputController"
}

// Inputs implements controller.Controller interface.
func (ctrl *KubernetesInputController) Inputs() []controller.Input {
	return []controller.Input{}
}

// Outputs implements controller.Controller interface.
func (ctrl *KubernetesInputController) Outputs() []controller.Output {
	return []controller.Output{
		{
			Type: myresource.KubernetesInputType,
			Kind: controller.OutputExclusive,
		},
	}
}

// Run implements controller.Controller interface.
//
//nolint:gocognit
func (ctrl *KubernetesInputController) Run(ctx context.Context, r controller.Runtime, l *zap.Logger) error {
	sourceMd := resource.NewMetadata(constants.NS, myresource.IntResourceType, constants.NS, resource.VersionUndefined)
	var eg errgroup.Group

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// reconcile one time
	err = reconcileK8sInputResources(ctx, r, l)
	if err != nil {
		return err
	}

	err = watcher.Add("./test/k8s")
	if err != nil {
		return err
	}

	defer watcher.Close()

	for {
		select {
		case <-ctx.Done():
			return eg.Wait()
		case <-r.EventCh():
			l.Info("received reconcile event for " + sourceMd.String())
		case e := <-watcher.Errors:
			return e
		case <-watcher.Events:
			err := reconcileK8sInputResources(ctx, r, l)
			if err != nil {
				return err
			}
		}
	}
}

func reconcileK8sInputResources(ctx context.Context, r controller.Runtime, l *zap.Logger) error {
	data, err := os.ReadFile("./test/k8s/1.yaml")
	if err != nil {
		return err
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	kubeconfig, err := clientConfig.ClientConfig()
	if err != nil {
		return err
	}

	client, err := rest.HTTPClientFor(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to create HTTP client for %q: %w", kubeconfig.Host, err)
	}

	mapper, err := apiutil.NewDynamicRESTMapper(kubeconfig, client)
	if err != nil {
		return err
	}

	reader := &manifestreader.StreamManifestReader{
		ReaderName: "inmem",
		Reader:     bytes.NewReader([]byte(data)),
		ReaderOptions: manifestreader.ReaderOptions{
			Validate:  true, // set true if you wire a RESTMapper
			Namespace: "",   // optional default namespace
			Mapper:    mapper,
		},
	}

	inputResources, err := reader.Read()
	if err != nil {
		return err
	}
	l.Info("reconciled " + strconv.Itoa(len(inputResources)) + " input resources")

	return nil
}
