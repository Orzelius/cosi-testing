package mystate

import (
	"fmt"

	"github.com/cosi-project/runtime/pkg/resource"
)

type eNotFound struct {
	error
}

func (eNotFound) NotFoundError() {}

// ErrNotFound generates error compatible with state.ErrNotFound.
func ErrNotFound(r resource.Pointer) error {
	return eNotFound{
		fmt.Errorf("resource %s doesn't exist", r),
	}
}
