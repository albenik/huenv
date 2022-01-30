package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"go.uber.org/multierr"
)

func BuildAndRun(dst, cfgPkg, cfgType string, buildFlags []string) (outerr error) {
	src, err := GenerateProgram(cfgPkg, cfgType)
	if err != nil {
		return fmt.Errorf("codegen: internal error: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	tmpDir, err := ioutil.TempDir(wd, "huenv_")
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			outerr = fmt.Errorf("failed to remove temp directory %q: %w", tmpDir, err)
		}
	}()

	var binFilename = "huenv"
	if runtime.GOOS == "windows" {
		binFilename += ".exe"
	} else {
		binFilename += ".bin"
	}

	if err = buildProgram(tmpDir, binFilename, src, buildFlags); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	cmdArgs := make([]string, 0, 2)
	if dst != "" {
		cmdArgs = append(cmdArgs, "-out", dst)
		defer func() {
			if outerr != nil {
				outerr = multierr.Append(outerr, removeDestinationFile(dst))
			}
		}()
	}
	cmd := exec.Command(filepath.Join(tmpDir, binFilename), cmdArgs...)
	cmd.Stdout = os.Stdout
	if err = runCmd(cmd); err != nil {
		return fmt.Errorf("exec: %w", &CmdError{err})
	}
	return nil
}

func GenerateProgram(srcPkg, srcTyp string) ([]byte, error) {
	srcTyp = srcTyp[(strings.LastIndex(srcTyp, ".") + 1):]

	buf := bytes.NewBuffer(nil)
	err := programSourceTemplate.Execute(buf, &struct {
		ConfigPackage string
		ConfigType    string
	}{
		ConfigPackage: srcPkg,
		ConfigType:    srcTyp,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func buildProgram(dir string, bin string, src []byte, buildFlags []string) error {
	const srcFilename = "main.go"
	if err := os.WriteFile(filepath.Join(dir, srcFilename), src, 0600); err != nil {
		return err
	}

	// Build the program.
	cmdArgs := []string{"build"}
	cmdArgs = append(cmdArgs, buildFlags...)
	cmdArgs = append(cmdArgs, "-o", bin, srcFilename)

	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = dir

	if err := runCmd(cmd); err != nil {
		return &CmdError{err}
	}
	return nil
}

func removeDestinationFile(name string) error {
	info, err := os.Stat(name)
	if err != nil {
		return fmt.Errorf("stat %q: %w", name, err)
	}
	if info.IsDir() {
		return fmt.Errorf("invalid destination: %q", name)
	}
	if err = os.Remove(name); err != nil {
		return fmt.Errorf("remove %q: %w", name, err)
	}
	return nil
}

var programSourceTemplate = template.Must(template.New("program").
	Parse(`package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	config {{ printf "%q" .ConfigPackage }}

	"github.com/albenik/huenv/generator"
	"github.com/albenik/huenv/reflector"
)

func main() {
	out := flag.String("out", "", "Output file. Defaults to stdout.")
	flag.Parse()

	if err := exec(*out); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		var de generator.DetailedError
		if errors.As(err, &de) {
			fmt.Fprintln(os.Stderr, de.Details())
		}
		os.Exit(1)
	}
}

func exec(out string) (outerr error) {
	dst := io.Writer(os.Stdout)

	if out != "" {
		file, err := os.Create(out)
		if err != nil {
			return err
		}
		defer func() {
			if err := file.Close(); err != nil {
				outerr = err
			}
		}()
		dst = file
	}

	info, err := reflector.New().Reflect(new(config.{{ .ConfigType }}))
	if err != nil {
		return err
	}

	if err = new(generator.ConfigGenerator).Generate(dst, info); err != nil {
		return err
	}

	return nil
}
`))
