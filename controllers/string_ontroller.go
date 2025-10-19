package controllers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Orzelius/cosi-testing/constants"
	"github.com/Orzelius/cosi-testing/myresource"
	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"go.uber.org/zap"
)

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
			Namespace: constants.NS,
			Type:      myresource.IntResourceType,
			Kind:      controller.InputStrong,
		},
		{
			Namespace: constants.NS,
			Type:      myresource.StringResourceType,
			Kind:      controller.InputDestroyReady,
		},
	}
}

// Outputs implements controller.Controller interface.
func (ctrl *IntToStrController) Outputs() []controller.Output {
	return []controller.Output{
		{
			Type: myresource.StringResourceType,
			Kind: controller.OutputExclusive,
		},
	}
}

// Run implements controller.Controller interface.
//
//nolint:gocognit
func (ctrl *IntToStrController) Run(ctx context.Context, r controller.Runtime, l *zap.Logger) error {
	sourceMd := resource.NewMetadata(constants.NS, myresource.IntResourceType, constants.NS, resource.VersionUndefined)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-r.EventCh():
		}

		l.Info("running reconcile")

		intList, err := safe.ReaderList[*myresource.IntResource](ctx, r, sourceMd)
		if err != nil {
			return fmt.Errorf("error listing objects: %w", err)
		}

		for intRes := range intList.All() {
			strRes := myresource.NewStringResource(intRes.Metadata().ID(), constants.NS)

			switch intRes.Metadata().Phase() {
			case resource.PhaseRunning:
				// AddFinalizer adds finalizer to resource metadata handling conflicts.
				// Finalizer is a free-form string which blocks resource destruction.
				// Resource can't be destroyed until all the finalizers are cleared.
				if err = r.AddFinalizer(ctx, intRes.Metadata(), resource.String(strRes)); err != nil {
					return fmt.Errorf("error adding finalizer: %w", err)
				}

				if err = safe.WriterModify(ctx, r, strRes,
					func(r *myresource.StringResource) error {
						r.TypedSpec().Val = strconv.Itoa(intRes.TypedSpec().Val)

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
