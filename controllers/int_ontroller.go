package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Orzelius/cosi-testing/constants"
	"github.com/Orzelius/cosi-testing/myresource"
	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/resource"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type IntController struct{}

// Name implements controller.Controller interface.
func (ctrl *IntController) Name() string {
	return "IntController"
}

// Inputs implements controller.Controller interface.
func (ctrl *IntController) Inputs() []controller.Input {
	return []controller.Input{
		{
			Namespace: constants.NS,
			Type:      myresource.IntResourceType,
			Kind:      controller.InputDestroyReady,
		},
		{
			Namespace: constants.NS,
			Type:      myresource.StringResourceType,
			Kind:      controller.InputStrong,
		},
	}
}

// Outputs implements controller.Controller interface.
func (ctrl *IntController) Outputs() []controller.Output {
	return []controller.Output{
		{
			Type: myresource.IntResourceType,
			Kind: controller.OutputExclusive,
		},
	}
}

// Run implements controller.Controller interface.
//
//nolint:gocognit
func (ctrl *IntController) Run(ctx context.Context, r controller.Runtime, l *zap.Logger) error {
	sourceMd := resource.NewMetadata(constants.NS, myresource.IntResourceType, constants.NS, resource.VersionUndefined)
	var eg errgroup.Group

	eg.Go(func() error {
		return runCreateController(ctx, r)
	})

	for {
		select {
		case <-ctx.Done():
			return eg.Wait()
		case <-r.EventCh():
			l.Info("received reconcile event for " + sourceMd.String())
		}
	}
}

func runCreateController(ctx context.Context, r controller.Runtime) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	i := 1
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			ints, err := r.List(ctx, resource.NewMetadata(constants.NS, myresource.IntResourceType, constants.NS, resource.VersionUndefined))
			if err != nil {
				panic(err)
			}
			if len(ints.Items) >= 6 {
				continue
			}

			intRes := myresource.NewIntResource(strconv.Itoa(i), rand.Intn(11))
			i++

			if err := r.Create(ctx, intRes); err != nil {
				return fmt.Errorf("error creating resource: %w", err)
			}
		}
	}
}
