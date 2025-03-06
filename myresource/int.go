package myresource

import (
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/meta"
	"github.com/cosi-project/runtime/pkg/resource/typed"
)

var _ resource.Resource = (*IntResource)(nil)

func NewIntResource(ns resource.Namespace, id resource.ID, val int) *IntResource {
	return typed.NewResource[IntSpec, IntExtension](
		resource.NewMetadata(ns, IntType, id, resource.VersionUndefined),
		IntSpec{Val: val},
	)
}

const (
	IntType = resource.Type("int")
)

type IntResource = typed.Resource[IntSpec, IntExtension]

type IntSpec = DeepCopyableSpec[int]

type IntExtension struct{}

func (IntExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{Type: IntType}
}
