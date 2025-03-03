package myresources

import (
	"github.com/cosi-project/runtime/pkg/resource"
)

// IntegerResource is implemented by resources holding ints.
type IntegerResource interface {
	Value() int
	SetValue(int)
}

// StringResource is implemented by resources holding strings.
type StringResource interface {
	Value() string
	SetValue(string)
}

// IntResourceType is the type of IntResource.
const IntResourceType = resource.Type("int")

// IntResource represents some integer value.
type IntResource = Resource[int, intSpec, *intSpec]

// NewIntResource creates new IntResource.
func NewIntResource(ns resource.Namespace, id resource.ID, value int) *IntResource {
	return NewResource[int, intSpec, *intSpec](resource.NewMetadata(ns, IntResourceType, id, resource.VersionUndefined), value)
}

type intSpec struct{ ValueGetSet[int] }

// StrResourceType is the type of StrResource.
const StrResourceType = resource.Type("str")

// StrResource represents some string value.
type StrResource = Resource[string, strSpec, *strSpec]

// NewStrResource creates new StrResource.
func NewStrResource(ns resource.Namespace, id resource.ID, value string) *StrResource {
	return NewResource[string, strSpec, *strSpec](resource.NewMetadata(ns, StrResourceType, id, resource.VersionUndefined), value)
}

type strSpec struct{ ValueGetSet[string] }

// SentenceResourceType is the type of SentenceResource.
const SentenceResourceType = resource.Type("sentence")

// SentenceResource represents some string value.
type SentenceResource = Resource[string, sentenceSpec, *sentenceSpec]

// NewSentenceResource creates new SentenceResource.
func NewSentenceResource(ns resource.Namespace, id resource.ID, value string) *SentenceResource {
	return NewResource[string, sentenceSpec, *sentenceSpec](resource.NewMetadata(ns, SentenceResourceType, id, resource.VersionUndefined), value)
}

type sentenceSpec struct{ ValueGetSet[string] }

// ValueGetSet is a basic building block for IntegerResource and StringResource implementations.
type ValueGetSet[T any] struct{ Val T }

func (s *ValueGetSet[T]) SetValue(t T) { s.Val = t }
func (s ValueGetSet[T]) Value() T      { return s.Val }
