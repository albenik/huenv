package generator

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

func formatSource(dst io.Writer, src io.Reader) error {
	cmd := exec.Command("gofmt")
	cmd.Stdin = src
	cmd.Stdout = dst
	return runCmd(cmd)
}

func runCmd(cmd *exec.Cmd) error {
	stderr := bytes.NewBuffer(nil)
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		ee := new(exec.ExitError)
		if errors.As(err, &ee) {
			ee.Stderr = stderr.Bytes()
		}
		return fmt.Errorf("%s: %w", cmd, ee)
	}
	return nil
}
