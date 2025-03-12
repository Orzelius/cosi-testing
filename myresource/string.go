package myresource

import (
	"github.com/Orzelius/cosi-testing/constants"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/meta"
	"github.com/cosi-project/runtime/pkg/resource/typed"
)

var _ resource.Resource = (*StringResource)(nil)

func NewStringResource(id resource.ID, val string) *StringResource {
	return typed.NewResource[StringSpec, StringExtension](
		resource.NewMetadata(constants.NS, StringResourceType, id, resource.VersionUndefined),
		StringSpec{Val: val},
	)
}

const (
	StringResourceType = resource.Type("string")
)

type StringResource = typed.Resource[StringSpec, StringExtension]

type StringSpec = DeepCopyableSpec[string]

type StringExtension struct{}

func (StringExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{Type: StringResourceType, DefaultNamespace: constants.NS}
}
