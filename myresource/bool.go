package myresource

import (
	"github.com/Orzelius/cosi-testing/constants"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/meta"
	"github.com/cosi-project/runtime/pkg/resource/typed"
)

var _ resource.Resource = (*BoolResource)(nil)

func NewBoolResource(id resource.ID, val bool) *BoolResource {
	return typed.NewResource[BoolSpec, BoolExtension](
		resource.NewMetadata(constants.NS, BoolResourceType, id, resource.VersionUndefined),
		BoolSpec{Val: val},
	)
}

const (
	BoolResourceType = resource.Type("bool")
)

type BoolResource = typed.Resource[BoolSpec, BoolExtension]

type BoolSpec = DeepCopyableSpec[bool]

type BoolExtension struct{}

func (BoolExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{Type: BoolResourceType, DefaultNamespace: constants.NS}
}
