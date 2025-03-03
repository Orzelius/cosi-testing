package controllers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Orzelius/cosi-testing/myresources"
	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"go.uber.org/zap"
)

const NS = "default"

// IntToStrController converts IntResource to StrResource.
type IntToStrController struct{}

// Name implements controller.Controller interface.
func (ctrl *IntToStrController) Name() string {
	return "IntToStrController"
}

// Inputs implements controller.Controller interface.
func (ctrl *IntToStrController) Inputs() []controller.Input {
	return []controller.Input{
		{
			Namespace: NS,
			Type:      myresources.IntResourceType,
			Kind:      controller.InputStrong,
		},
		{
			Namespace: NS,
			Type:      myresources.StrResourceType,
			Kind:      controller.InputDestroyReady,
		},
	}
}

// Outputs implements controller.Controller interface.
func (ctrl *IntToStrController) Outputs() []controller.Output {
	return []controller.Output{
		{
			Type: myresources.StrResourceType,
			Kind: controller.OutputExclusive,
		},
	}
}

// Run implements controller.Controller interface.
//
//nolint:gocognit
func (ctrl *IntToStrController) Run(ctx context.Context, r controller.Runtime, _ *zap.Logger) error {
	sourceMd := resource.NewMetadata(NS, myresources.IntResourceType, "", resource.VersionUndefined)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-r.EventCh():
		}

		intList, err := safe.ReaderList[interface {
			myresources.IntegerResource
			resource.Resource
		}](ctx, r, sourceMd)
		if err != nil {
			return fmt.Errorf("error listing objects: %w", err)
		}

		for intRes := range intList.All() {
			strRes := myresources.NewStrResource(NS, intRes.Metadata().ID(), "")

			switch intRes.Metadata().Phase() {
			case resource.PhaseRunning:
				if err = r.AddFinalizer(ctx, intRes.Metadata(), resource.String(strRes)); err != nil {
					return fmt.Errorf("error adding finalizer: %w", err)
				}

				if err = safe.WriterModify(ctx, r, strRes,
					func(r *myresources.StrResource) error {
						r.SetValue(strconv.Itoa(intRes.Value()))

						return nil
					}); err != nil {
					return fmt.Errorf("error updating objects: %w", err)
				}
			case resource.PhaseTearingDown:
				ready, err := r.Teardown(ctx, strRes.Metadata())
				if err != nil && !state.IsNotFoundError(err) {
					return fmt.Errorf("error tearing down: %w", err)
				}

				if err == nil && !ready {
					// reconcile other resources, reconcile loop will be triggered once resource is ready for teardown
					continue
				}

				if err = r.Destroy(ctx, strRes.Metadata()); err != nil && !state.IsNotFoundError(err) {
					return fmt.Errorf("error destroying: %w", err)
				}

				if err = r.RemoveFinalizer(ctx, intRes.Metadata(), resource.String(strRes)); err != nil {
					if !state.IsNotFoundError(err) {
						return fmt.Errorf("error removing finalizer (str controller): %w", err)
					}
				}
			}
		}

		r.ResetRestartBackoff()
	}
}
