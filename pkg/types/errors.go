package types

import (
	"fmt"
)

const rootCodeSpace = "redisson"

var (
	ErrWaitTimeout = register(rootCodeSpace, 10000, "wait timeout")
	ErrMismatch    = register(rootCodeSpace, 20001, "identity mismatch")
)

var usedCode = map[string]struct{}{}

func register(codeSpace string, code uint32, desc string) Error {
	err := sdkError{
		codeSpace: codeSpace,
		code:      code,
		desc:      desc,
	}
	usedCode[fmt.Sprintf("%s:%d", codeSpace, code)] = struct{}{}
	return err
}

type Error interface {
	Error() string
	Code() uint32
	CodeSpace() string
}

type sdkError struct {
	codeSpace string
	code      uint32
	desc      string
}

func (s sdkError) Error() string {
	return s.desc
}
func (s sdkError) CodeSpace() string {
	return s.codeSpace
}
func (s sdkError) Code() uint32 {
	return s.code
}
