package huenv_test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/albenik/huenv"
	"github.com/albenik/huenv/unmarshal"
)

var _ huenv.Config = (*TestConfig)(nil)

type TestConfig struct {
	Bool           bool          `env:"ENV_BOOL"`
	BoolOpt        bool          `env:"ENV_BOOL_OPT,optional"`
	String         string        `env:"ENV_STRING"`
	StringEnum     string        `env:"ENV_STRING_ENUM,enum(foo,bar,baz)"`
	StringOptional string        `env:"ENV_STRING_OPT,optional"`
	L2             *TestConfigL2 `env:"L2_*"`
}

type TestConfigL2 struct {
	IntField int64 `env:"INT"`
}

func (c *TestConfig) Envmap() map[string]*unmarshal.Target {
	to := unmarshal.NewTarget

	return map[string]*unmarshal.Target{
		"ENV_BOOL":        to(&c.Bool, new(unmarshal.Bool), unmarshal.Required()),
		"ENV_BOOL_OPT":    to(&c.BoolOpt, new(unmarshal.Bool), unmarshal.Optional()),
		"ENV_STRING":      to(&c.String, new(unmarshal.String), unmarshal.Required()),
		"ENV_STRING_ENUM": to(&c.StringEnum, new(unmarshal.String), unmarshal.Enum("foo", "bar")),
		"ENV_STRING_OPT":  to(&c.StringOptional, new(unmarshal.String), unmarshal.Optional()),
		"L2_INT":          to(&c.L2.IntField, new(unmarshal.Int64), unmarshal.Required()),
	}
}

func TestEnvconfig_Init(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Setenv("ENV_BOOL", "true")
		t.Setenv("ENV_STRING", "foo")
		t.Setenv("ENV_STRING_ENUM", "bar")
		t.Setenv("L2_INT", "777")

		conf := new(TestConfig)
		err := huenv.Init(conf, "")

		require.NoError(t, err)
		assert.Equal(t, &TestConfig{
			Bool:           true,
			BoolOpt:        false,
			String:         "foo",
			StringEnum:     "bar",
			StringOptional: "",
			L2: &TestConfigL2{
				IntField: 777,
			},
		}, conf)
	})

	// TODO Implement errors tests
	// t.Run("Errors", func(t *testing.T) {
	// 	type testConf struct {
	// 		FieldBool     bool
	// 		FieldBytes    []byte
	// 		FieldDuration time.Duration
	// 		FieldTime     time.Time
	// 		FieldString   string
	// 	}
	// })
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
