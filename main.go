package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cosi-project/runtime/pkg/controller/conformance"
	"github.com/cosi-project/runtime/pkg/controller/runtime"
	"github.com/cosi-project/runtime/pkg/controller/runtime/options"
	"github.com/cosi-project/runtime/pkg/logging"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/cosi-project/runtime/pkg/state/impl/inmem"
	"github.com/cosi-project/runtime/pkg/state/impl/namespaced"
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

	inmemState := state.WrapCore(namespaced.NewState(inmem.Build))
	logger := logging.DefaultLogger()

	controllerRuntime, err := runtime.NewRuntime(inmemState, logger, options.WithMetrics(true))
	if err != nil {
		return fmt.Errorf("error setting up controller runtime: %w", err)
	}

	var eg errgroup.Group

	eg.Go(func() error {
		ctrl := &conformance.IntToStrController{
			SourceNamespace: "default",
			TargetNamespace: "default",
		}

		if err := controllerRuntime.RegisterController(ctrl); err != nil {
			return fmt.Errorf("error registering controller: %w", err)
		}

		return controllerRuntime.Run(ctx)
	})

	eg.Go(func() error {
		return runController(ctx, inmemState)
	})

	logger.Info("waiting for <-ctx.Done()")
	<-ctx.Done()

	logger.Info("waiting for eg.Wait()")
	return eg.Wait()
}

func runController(ctx context.Context, st state.State) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	i := 1

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			intRes := conformance.NewIntResource("default", fmt.Sprintf("int-%d", i), i)

			i++

			if err := st.Create(ctx, intRes); err != nil {
				return fmt.Errorf("error creating resource: %w", err)
			}
		}
	}
}
