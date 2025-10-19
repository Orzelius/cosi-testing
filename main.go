package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Orzelius/cosi-testing/constants"
	"github.com/Orzelius/cosi-testing/controllers"
	"github.com/Orzelius/cosi-testing/myresource"
	"github.com/cosi-project/runtime/pkg/controller/runtime"
	"github.com/cosi-project/runtime/pkg/controller/runtime/options"
	"github.com/cosi-project/runtime/pkg/logging"
	cosistate "github.com/cosi-project/runtime/pkg/state"
	"github.com/cosi-project/runtime/pkg/state/impl/inmem"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	// fileStateCore := mystate.NewState()
	// go fileStateCore.StartFileWatcher(ctx)
	// defer fileStateCore.CloseFileWatcher()
	// fileState := cosistate.WrapCore(fileStateCore)

	state := cosistate.WrapCore(inmem.NewState(constants.NS))

	logger := logging.DefaultLogger()

	controllerRuntime, err := runtime.NewRuntime(state, logger, options.WithMetrics(true))
	if err != nil {
		return fmt.Errorf("error setting up controller runtime: %w", err)
	}

	var eg errgroup.Group

	eg.Go(func() error {
		// if err := controllerRuntime.RegisterQController(&controllers.QIntToStrController{}); err != nil {
		// 	return fmt.Errorf("error registering controller: %w", err)
		// }
		// if err := controllerRuntime.RegisterController(&controllers.IntController{}); err != nil {
		// 	return fmt.Errorf("error registering controller: %w", err)
		// }

		if err := controllerRuntime.RegisterController(&controllers.KubernetesInputController{}); err != nil {
			return fmt.Errorf("error registering controller: %w", err)
		}

		return controllerRuntime.Run(ctx)
	})

	eg.Go(func() error {
		return runCreateController(ctx, controllerRuntime.CachedState())
	})

	<-ctx.Done()
	return eg.Wait()
}

func runCreateController(ctx context.Context, r cosistate.CoreState) error {
	k8sManifests := `apiVersion: v1
kind: Namespace
metadata:
  name: test-lab
  labels:
    app.kubernetes.io/part-of: test-lab
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: test-lab
data:
  APP_MESSAGE: "hello from configmap"
  APP_PORT: "8080"
---
apiVersion: v1
kind: Secret
metadata:
  name: app-secret
  namespace: test-lab
type: Opaque
stringData:
  PASSWORD: "dummy-password"
`

	intBagSmall := myresource.NewKubernetesInput(strconv.Itoa(1), k8sManifests)
	if err := r.Create(ctx, intBagSmall); err != nil {
		return fmt.Errorf("error creating resource: %w", err)
	}

	return nil
}
