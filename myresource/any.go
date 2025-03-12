package myresource

import (
	"github.com/Orzelius/cosi-testing/constants"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/meta"
	"github.com/cosi-project/runtime/pkg/resource/typed"
)

var _ resource.Resource = (*AnyResource)(nil)

func NewAnyResource(meta resource.Metadata, val any) *AnyResource {
	return typed.NewResource[AnySpec, AnyExtension](meta, AnySpec{Val: val})
}

const (
	AnyResourceType = resource.Type("Any")
)

type AnyResource = typed.Resource[AnySpec, AnyExtension]

type AnySpec = DeepCopyableSpec[any]

type AnyExtension struct{}

func (AnyExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{Type: AnyResourceType, DefaultNamespace: constants.NS}
}
