package myresource

import (
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/typed"
)

type Marshalable struct {
	Val any                `yaml:"spec"`
	Md  *resource.Metadata `yaml:"metadata"`
}
type UnMarshalable struct {
	Val ValVal            `yaml:"spec"`
	Md  resource.Metadata `yaml:"metadata"`
}
type ValVal struct {
	Val any
}

func (r UnMarshalable) Metadata() *resource.Metadata { return &r.Md }
func (r UnMarshalable) Spec() any                    { return r.Val }
func (r UnMarshalable) DeepCopy() resource.Resource  { return UnMarshalable{Val: r.Val, Md: r.Md} }

type DeepCopyableSpec[T any] struct {
	Val T
}

func (s DeepCopyableSpec[T]) DeepCopy() DeepCopyableSpec[T] {
	return DeepCopyableSpec[T]{Val: s.Val}
}

func Cast(meta resource.Metadata, val any) resource.Resource {
	switch meta.Type() {
	case IntResourceType:
		v, _ := val.(int)
		return typed.NewResource[IntSpec, IntExtension](meta, IntSpec{Val: v})
	case StringResourceType:
		v, _ := val.(string)
		return typed.NewResource[StringSpec, StringExtension](meta, StringSpec{Val: v})
	case AnyResourceType:
		return typed.NewResource[AnySpec, AnyExtension](meta, AnySpec{Val: val})
	}
	panic("failed to cast")
}
