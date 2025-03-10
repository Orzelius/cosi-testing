package myresource

import (
	"github.com/Orzelius/cosi-testing/constants"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/meta"
	"github.com/cosi-project/runtime/pkg/resource/typed"
)

var _ resource.Resource = (*IntResource)(nil)

func NewIntResource(id resource.ID, val int) *IntResource {
	return typed.NewResource[IntSpec, IntExtension](
		resource.NewMetadata(constants.NS, IntResourceType, id, resource.VersionUndefined),
		IntSpec{Val: val},
	)
}

const (
	IntResourceType = resource.Type("int")
)

type IntResource = typed.Resource[IntSpec, IntExtension]

type IntSpec = DeepCopyableSpec[int]

type IntExtension struct{}

func (IntExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{Type: IntResourceType}
}
