package myresource

import (
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/meta"
	"github.com/cosi-project/runtime/pkg/resource/typed"
)

var _ resource.Resource = (*StringResource)(nil)

func NewStringResource(ns resource.Namespace, id resource.ID, val string) *StringResource {
	return typed.NewResource[StringSpec, StringExtension](
		resource.NewMetadata(ns, StringType, id, resource.VersionUndefined),
		StringSpec{Val: val},
	)
}

const (
	StringType = resource.Type("String")
)

type StringResource = typed.Resource[StringSpec, StringExtension]

type StringSpec = DeepCopyableSpec[string]

type StringExtension struct{}

func (StringExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{Type: StringType}
}
