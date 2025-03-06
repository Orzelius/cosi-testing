package myresource

import (
	"github.com/cosi-project/runtime/pkg/resource"
)

type Marshalable struct {
	Val valVal `yaml:"spec"`
	Md  resource.Metadata
}
type valVal struct {
	Val any
}

func (r Marshalable) Metadata() *resource.Metadata { return &r.Md }
func (r Marshalable) Spec() any                    { return r.Val }
func (r Marshalable) DeepCopy() resource.Resource  { return Marshalable{Val: r.Val, Md: r.Md} }

type DeepCopyableSpec[T any] struct {
	Val T
}

func (s DeepCopyableSpec[T]) DeepCopy() DeepCopyableSpec[T] {
	return DeepCopyableSpec[T]{Val: s.Val}
}
