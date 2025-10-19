package controllers

import (
	"context"
	"strconv"

	"github.com/Orzelius/cosi-testing/myresource"
	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/controller/generic/qtransform"
	"go.uber.org/zap"
)

var QTransformController = qtransform.NewQController(
	qtransform.Settings[*myresource.IntResource, *myresource.StringResource]{
		Name: "IntToStrController",
		MapMetadataFunc: func(ir *myresource.IntResource) *myresource.StringResource {
			return myresource.NewStringResource(ir.Metadata().ID(), "")
		},
		UnmapMetadataFunc: func(sr *myresource.StringResource) *myresource.IntResource {
			return myresource.NewIntResource(sr.Metadata().ID(), 0)
		},
		TransformFunc: func(ctx context.Context, r controller.Reader, l *zap.Logger, ir *myresource.IntResource, sr *myresource.StringResource) error {
			sr.TypedSpec().Val = strconv.Itoa(ir.TypedSpec().Val)
			return nil
		},
		FinalizerRemovalFunc: func(ctx context.Context, r controller.Reader, l *zap.Logger, ir *myresource.IntResource) error {
			return nil
		},
	},
	qtransform.WithConcurrency(2),
)
