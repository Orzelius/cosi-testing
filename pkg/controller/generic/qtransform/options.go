// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package qtransform

import (
	"context"
	"fmt"

	"github.com/siderolabs/gen/optional"
	"go.uber.org/zap"

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/controller/generic"
	"github.com/cosi-project/runtime/pkg/resource"
)

type namespaceType struct {
	Namespace resource.Namespace
	Type      resource.Type
}

type mapperFunc func(context.Context, *zap.Logger, controller.QRuntime, resource.Resource) ([]resource.Pointer, error)

// MapperFuncGeneric is a generic version of mapperFunc.
type MapperFuncGeneric[I generic.ResourceWithRD] func(context.Context, *zap.Logger, controller.QRuntime, I) ([]resource.Pointer, error)

// MapperSameID is a mapper that returns the same namespace ID as the input resource, but uses output resource type.
func MapperSameID[I generic.ResourceWithRD, O generic.ResourceWithRD]() MapperFuncGeneric[I] {
	var zeroOutput O

	outputType := zeroOutput.ResourceDefinition().Type

	return func(_ context.Context, _ *zap.Logger, _ controller.QRuntime, v I) ([]resource.Pointer, error) {
		return []resource.Pointer{resource.NewMetadata(v.Metadata().Namespace(), outputType, v.Metadata().ID(), resource.VersionUndefined)}, nil
	}
}

// MapperNone is a mapper that returns no pointers.
func MapperNone[I generic.ResourceWithRD]() MapperFuncGeneric[I] {
	return func(context.Context, *zap.Logger, controller.QRuntime, I) ([]resource.Pointer, error) {
		return nil, nil
	}
}

func mapperFuncFromGeneric[I generic.ResourceWithRD](generic MapperFuncGeneric[I]) mapperFunc {
	return func(ctx context.Context, logger *zap.Logger, r controller.QRuntime, res resource.Resource) ([]resource.Pointer, error) {
		v, ok := res.(I)
		if !ok {
			return nil, fmt.Errorf("unexpected resource type in mapFunc %T", res)
		}

		return generic(ctx, logger, r, v)
	}
}

// ControllerOptions configures QTransformController.
type ControllerOptions struct {
	mappers            map[namespaceType]mapperFunc
	leftoverFinalizers map[resource.Finalizer]struct{}
	extraInputs        []controller.Input
	extraOutputs       []controller.Output
	primaryOutputKind  controller.OutputKind
	concurrency        optional.Optional[uint]
}

// ControllerOption is an option for QTransformController.
type ControllerOption func(*ControllerOptions)

// WithExtraMappedInput adds an extra  mapped inputs to the controller.
func WithExtraMappedInput[I generic.ResourceWithRD](mapFunc MapperFuncGeneric[I]) ControllerOption {
	return func(o *ControllerOptions) {
		var zeroInput I

		if o.mappers == nil {
			o.mappers = map[namespaceType]mapperFunc{}
		}

		nsType := namespaceType{
			Namespace: zeroInput.ResourceDefinition().DefaultNamespace,
			Type:      zeroInput.ResourceDefinition().Type,
		}

		if _, present := o.mappers[nsType]; present {
			panic(fmt.Errorf("duplicate mapper for %q", nsType))
		}

		o.mappers[nsType] = mapperFuncFromGeneric(mapFunc)

		o.extraInputs = append(o.extraInputs, controller.Input{
			Namespace: zeroInput.ResourceDefinition().DefaultNamespace,
			Type:      zeroInput.ResourceDefinition().Type,
			Kind:      controller.InputQMapped,
		})
	}
}

// WithExtraOutputs adds extra outputs to the controller.
func WithExtraOutputs(outputs ...controller.Output) ControllerOption {
	return func(o *ControllerOptions) {
		o.extraOutputs = append(o.extraOutputs, outputs...)
	}
}

// WithOutputKind sets main output resource kind.
func WithOutputKind(kind controller.OutputKind) ControllerOption {
	return func(o *ControllerOptions) {
		o.primaryOutputKind = kind
	}
}

// WithConcurrency sets the maximum number of concurrent reconciles.
func WithConcurrency(n uint) ControllerOption {
	return func(o *ControllerOptions) {
		o.concurrency = optional.Some(n)
	}
}

// WithIgnoreTeardownUntil ignores input resource teardown until the input resource has only mentioned finalizers left.
//
// This allows to keep output resources not destroyed until other controllers remove their finalizers.
//
// Implicitly the controller will also ignore its own finalizer, so if the list is empty, the controller will wait
// to be the last one not done with the resource.
func WithIgnoreTeardownUntil(finalizers ...resource.Finalizer) ControllerOption {
	return func(o *ControllerOptions) {
		if o.leftoverFinalizers == nil {
			o.leftoverFinalizers = map[resource.Finalizer]struct{}{}
		}

		for _, fin := range finalizers {
			o.leftoverFinalizers[fin] = struct{}{}
		}
	}
}
