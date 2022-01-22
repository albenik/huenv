package huenv_test

import (
	"bytes"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/albenik/huenv"
	"github.com/albenik/huenv/unmarshal"
)

var _ huenv.Config = (*testConfig)(nil)

type testConfig struct {
	Bool           bool   `env:"ENV_BOOL"`
	BoolOpt        bool   `env:"ENV_BOOL_OPT,optional"`
	String         string `env:"ENV_STRING"`
	StringEnum     string `env:"ENV_STRING_ENUM,enum(foo,bar,baz)"`
	StringOptional string `env:"ENV_STRING_OPT,optional"`
}

func (c *testConfig) Envmap() map[string]*unmarshal.Target {
	to := unmarshal.NewTarget

	return map[string]*unmarshal.Target{
		"ENV_BOOL":        to(&c.Bool, new(unmarshal.Bool), unmarshal.Required()),
		"ENV_BOOL_OPT":    to(&c.BoolOpt, new(unmarshal.Bool), unmarshal.Optional()),
		"ENV_STRING":      to(&c.String, new(unmarshal.String), unmarshal.Required()),
		"ENV_STRING_ENUM": to(&c.StringEnum, new(unmarshal.String), unmarshal.Enum("foo", "bar")),
		"ENV_STRING_OPT":  to(&c.StringOptional, new(unmarshal.String), unmarshal.Optional()),
	}
}

func TestEnvconfig_Init(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Setenv("ENV_BOOL", "true")
		t.Setenv("ENV_STRING", "foo")
		t.Setenv("ENV_STRING_ENUM", "bar")

		conf := new(testConfig)
		err := huenv.Init(conf, "")

		require.NoError(t, err)
		assert.Equal(t, &testConfig{
			Bool:           true,
			BoolOpt:        false,
			String:         "foo",
			StringEnum:     "bar",
			StringOptional: "",
		}, conf)
	})

	t.Run("Errors", func(t *testing.T) {
		type testConf struct {
			FieldBool     bool
			FieldBytes    []byte
			FieldDuration time.Duration
			FieldTime     time.Time
			FieldString   string
		}
	})
}

func TestGenerator(t *testing.T) {
	stderr := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)

	cmd := exec.Command(
		"go", "run", "github.com/albenik/huenv/cmd/huenv",
		"-out", "testdata/config_unmarshal.go",
		"github.com/albenik/huenv/testdata", "TestConfig1",
	)
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	err := cmd.Run()
	if !assert.NoError(t, err) {
		t.Logf("=== STDERR ===\n%s", stderr.String())
	}
	t.Logf("=== STDOUT ===\n%s", stdout.String())
}
