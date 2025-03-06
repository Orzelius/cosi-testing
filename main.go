package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Orzelius/cosi-testing/myresource"
	"github.com/Orzelius/cosi-testing/mystate"
	"github.com/cosi-project/runtime/pkg/logging"
	cosistate "github.com/cosi-project/runtime/pkg/state"
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

	fileState := cosistate.WrapCore(mystate.NewState())
	logger := logging.DefaultLogger()

	// controllerRuntime, err := runtime.NewRuntime(fileState, logger, options.WithMetrics(true))
	// if err != nil {
	// 	return fmt.Errorf("error setting up controller runtime: %w", err)
	// }

	var eg errgroup.Group

	// eg.Go(func() error {
	// 	if err := controllerRuntime.RegisterController(&controllers.IntToStrController{}); err != nil {
	// 		return fmt.Errorf("error registering controller: %w", err)
	// 	}

	// 	return controllerRuntime.Run(ctx)
	// })

	eg.Go(func() error {
		return runCreateController(ctx, fileState)
	})

	logger.Info("waiting for <-ctx.Done()")
	<-ctx.Done()

	logger.Info("waiting for eg.Wait()")
	return eg.Wait()
}

func runCreateController(ctx context.Context, st cosistate.State) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	i := 1

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			intRes := myresource.NewIntResource("", strconv.Itoa(i), i)
			i++

			if err := st.Create(ctx, intRes); err != nil {
				return fmt.Errorf("error creating resource: %w", err)
			}
		}
	}
}
