package app

import (
	"testing"

	"github.com/pkg/errors"
)

func TestError(t *testing.T) {
	err1 := errors.Wrap(errors.New("Error1"), "funcA.CallA")
	err2 := errors.Wrap(err1, "FuncB.CallA")
	restErr := NewRestErrorWithMessage(ErrInternalServerError, "Error parsing smth", err2)
	t.Log(restErr)
	err := ParseErrors(errors.New("hai"))
	t.Log(err)
}
