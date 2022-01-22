package generator

import (
	"bytes"
	"errors"
	"os/exec"
)

type DetailedError interface {
	Details() string
}

type CmdError struct {
	wrapped error
}

func (e *CmdError) Error() string {
	return e.wrapped.Error()
}

func (e *CmdError) Details() string {
	buf := bytes.NewBuffer(nil)
	ee := new(exec.ExitError)
	if errors.As(e.wrapped, &ee) {
		buf.Write(ee.Stderr)
	}
	return buf.String()
}

type CodegenError struct {
	wrapped error
	source  []byte
}

func (e *CodegenError) Error() string {
	return e.wrapped.Error()
}

func (e *CodegenError) Details() string {
	buf := bytes.NewBuffer(nil)
	ee := new(exec.ExitError)
	if errors.As(e.wrapped, &ee) {
		buf.WriteString("STDERR:\n")
		buf.Write(ee.Stderr)
		buf.WriteRune('\n')
	}
	buf.WriteString("SOURCE:\n")
	buf.Write(e.source)
	return buf.String()
}
