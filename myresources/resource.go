package myresources

import "github.com/cosi-project/runtime/pkg/resource"

// Resource represents some T value.
type Resource[T any, S Spec[T], SS SpecPtr[T, S]] struct {
	Val S
	md  resource.Metadata
}

// NewResource creates new Resource.
func NewResource[T any, S Spec[T], SS SpecPtr[T, S]](md resource.Metadata, value T) *Resource[T, S, SS] {
	var s S
	ss := SS(&s)
	ss.SetValue(value)

	r := &Resource[T, S, SS]{
		md:  md,
		Val: s,
	}

	return r
}

// Metadata implements resource.Resource.
func (r *Resource[T, S, SS]) Metadata() *resource.Metadata {
	return &r.md
}

// Spec implements resource.Resource.
func (r *Resource[T, S, SS]) Spec() any {
	return r.Val
}

// Value returns a value inside the spec.
func (r *Resource[T, S, SS]) Value() T { //nolint:ireturn
	return r.Val.Value()
}

// SetValue set spec with provided value.
func (r *Resource[T, S, SS]) SetValue(v T) {
	val := SS(&r.Val)
	val.SetValue(v)
}

// DeepCopy implements resource.Resource.
func (r *Resource[T, S, SS]) DeepCopy() resource.Resource { //nolint:ireturn
	return &Resource[T, S, SS]{
		md:  r.md,
		Val: r.Val,
	}
}

// SpecPtr requires Spec to be a pointer and have a set of methods.
type SpecPtr[T, S any] interface {
	*S
	Spec[T]
	SetValue(T)
}

// Spec requires spec to have a set of Get methods.
type Spec[T any] interface {
	Value() T
}
