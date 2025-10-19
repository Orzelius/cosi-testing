package myresource

import (
	"github.com/Orzelius/cosi-testing/constants"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/meta"
	"github.com/cosi-project/runtime/pkg/resource/typed"
)

var _ resource.Resource = (*KubernetesInput)(nil)

func NewKubernetesInput(id resource.ID, val string) *KubernetesInput {
	return typed.NewResource[KubernetesInputSpec, KubernetesInputExtension](
		resource.NewMetadata(constants.NS, KubernetesInputType, id, resource.VersionUndefined),
		KubernetesInputSpec{Val: val},
	)
}

const (
	KubernetesInputType = resource.Type("kubernetes-input")
)

type KubernetesInput = typed.Resource[KubernetesInputSpec, KubernetesInputExtension]

type KubernetesInputSpec = DeepCopyableSpec[string]

type KubernetesInputExtension struct{}

func (KubernetesInputExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{Type: KubernetesInputType, DefaultNamespace: constants.NS}
}
