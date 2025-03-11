package controllers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Orzelius/cosi-testing/constants"
	"github.com/Orzelius/cosi-testing/myresource"
	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/siderolabs/gen/optional"
	"go.uber.org/zap"
)

// QIntToStrController converts IntResource to StrResource as a QController.
type QIntToStrController struct {
	ShutdownCalled bool
}

// Name implements controller.QController interface.
func (ctrl *QIntToStrController) Name() string {
	return "QIntToStrController"
}

// Settings implements controller.QController interface.
func (ctrl *QIntToStrController) Settings() controller.QSettings {
	failRunHook := true

	return controller.QSettings{
		Inputs: []controller.Input{
			{
				Namespace: constants.NS,
				Type:      myresource.IntResourceType,
				Kind:      controller.InputQPrimary,
			},
			{
				Namespace: constants.NS,
				Type:      myresource.StringResourceType,
				Kind:      controller.InputQMappedDestroyReady,
			},
		},
		Outputs: []controller.Output{
			{
				Type: myresource.StringResourceType,
				Kind: controller.OutputExclusive,
			},
		},
		Concurrency: optional.Some(uint(4)),
		RunHook: func(ctx context.Context, _ *zap.Logger, r controller.QRuntime) error {
			interval := time.NewTicker(time.Second)

			defer interval.Stop()

			if failRunHook {
				failRunHook = false

				return errors.New("oh no")
			}

			for {
				select {
				case <-ctx.Done():
					return nil
				case <-interval.C:
					// if err := safe.WriterModify(ctx, r, NewStrResource("hooks", "lastInvokation", ""), func(r *StrResource) error {
					// 	r.value.value = time.Now().Format(time.RFC3339)

					// 	return nil
					// }); err != nil {
					// 	return err
					// }
				}
			}
		},
		ShutdownHook: func() {
			ctrl.ShutdownCalled = true
		},
	}
}

// Reconcile implements controller.QController interface.
func (ctrl *QIntToStrController) Reconcile(ctx context.Context, _ *zap.Logger, r controller.QRuntime, ptr resource.Pointer) error {
	src, err := safe.ReaderGet[*myresource.IntResource](ctx, r, ptr)
	if err != nil {
		if state.IsNotFoundError(err) {
			return nil
		}

		return err
	}

	switch src.Metadata().Phase() {
	case resource.PhaseTearingDown:
		// cleanup destination resource as needed
		dst := myresource.NewStringResource(src.Metadata().ID(), "").Metadata()

		ready, err := r.Teardown(ctx, dst)
		if err != nil {
			if state.IsNotFoundError(err) {
				return r.RemoveFinalizer(ctx, ptr, ctrl.Name())
			}

			return err
		}

		if !ready {
			// not ready for teardown, wait
			return nil
		}

		if err := r.Destroy(ctx, dst); err != nil {
			return err
		}

		return r.RemoveFinalizer(ctx, ptr, ctrl.Name())
	case resource.PhaseRunning:
		if err := r.AddFinalizer(ctx, ptr, ctrl.Name()); err != nil {
			return err
		}

		strValue := strconv.Itoa(src.TypedSpec().Val)

		return safe.WriterModify(ctx, r, myresource.NewStringResource(src.Metadata().ID(), strValue), func(r *myresource.StringResource) error {
			r.TypedSpec().Val = strValue

			return nil
		})
	default:
		panic("unexpected phase")
	}
}

// MapInput implements controller.QController interface.
func (ctrl *QIntToStrController) MapInput(_ context.Context, _ *zap.Logger, _ controller.QRuntime, ptr resource.Pointer) ([]resource.Pointer, error) {
	switch {
	case ptr.Type() == myresource.StringResourceType:
		// remap output to input to recheck on finalizer removal
		return []resource.Pointer{resource.NewMetadata(constants.NS, myresource.IntResourceType, ptr.ID(), resource.VersionUndefined)}, nil
	default:
		return nil, fmt.Errorf("unexpected input %s", ptr)
	}
}
