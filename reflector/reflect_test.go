package reflector_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/albenik/huenv/reflector"
	"github.com/albenik/huenv/unmarshal"
)

type TestReflectConfig struct {
	Field1     string                       `env:"FIELD1"`
	Field2     bool                         `env:"FIELD2,optional"`
	Field3     int                          `env:"FIELD3,reqif(Field1=foo)"`
	Subconfig1 *TestReflectConfigSubconfig1 `env:"*,reqif(Field2=true)"`
	Subconfig2 *TestReflectConfigSubconfig2 `env:"SUB2_*,reqif(Field3=7)"`
	Subconfig3 *TestReflectConfigSubconfig1 `env:"SUB3" unmarshaler:"github.com/albenik/huenv/reflector TestUnmarshaler1"`
}

type TestReflectConfigSubconfig1 struct {
	Field1 string `env:"SUB1FIELD1"`
}

type TestReflectConfigSubconfig2 struct {
	Field1 string `env:"SUB2FIELD1"`
	Field2 string `env:"SUB2FIELD2,optional"`
	Field3 string `env:"SUB2FIELD3,enum(foo,bar,baz)"`
}

func TestReflector_Reflect(t *testing.T) {
	require.NoError(t, reflector.Register((*reflector.TestUnmarshaler1)(nil)))

	const (
		pkgUnmarshal      = "github.com/albenik/huenv/unmarshal"
		stringUnmarshaler = "String"
		boolUnmarshaler   = "Bool"
		intUnmarshaler    = "Int"

		pkgTest         = "github.com/albenik/huenv/reflector"
		pkg1Unmarshaler = "TestUnmarshaler1"
	)

	expected := &reflector.Result{
		ConfigPkg:  "reflector_test",
		ConfigType: "TestReflectConfig",
		Packages: map[string]struct{}{
			pkgUnmarshal: {},
			pkgTest:      {},
		},
		Envs: map[string]*reflector.Target{
			"FIELD1": {
				Field: &reflector.TargetField{
					Name: "Field1",
					Unmarshaler: unmarshal.UnmarshalerName{
						Package: pkgUnmarshal,
						Type:    stringUnmarshaler,
					},
				},
				Condition: reflector.ConditionRequired(true),
			},
			"FIELD2": {
				Field: &reflector.TargetField{
					Name: "Field2",
					Unmarshaler: unmarshal.UnmarshalerName{
						Package: pkgUnmarshal,
						Type:    boolUnmarshaler,
					},
				},
				Condition: reflector.ConditionRequired(false),
			},
			"FIELD3": {
				Field: &reflector.TargetField{
					Name: "Field3",
					Unmarshaler: unmarshal.UnmarshalerName{
						Package: pkgUnmarshal,
						Type:    intUnmarshaler,
					},
				},
				Condition: &reflector.ConditionRequireIf{
					Target: &reflector.TargetField{
						Name: "Field1",
						Unmarshaler: unmarshal.UnmarshalerName{
							Package: pkgUnmarshal,
							Type:    stringUnmarshaler,
						},
					},
					ValueStr: "foo",
				},
			},
			"SUB1FIELD1": {
				Field: &reflector.TargetField{
					Name: "Subconfig1.Field1",
					Unmarshaler: unmarshal.UnmarshalerName{
						Package: pkgUnmarshal,
						Type:    stringUnmarshaler,
					},
				},
				Condition: &reflector.ConditionRequireIf{
					Target: &reflector.TargetField{
						Name: "Field2",
						Unmarshaler: unmarshal.UnmarshalerName{
							Package: pkgUnmarshal,
							Type:    boolUnmarshaler,
						},
					},
					ValueStr: "true",
				},
			},
			"SUB2_SUB2FIELD1": {
				Field: &reflector.TargetField{
					Name: "Subconfig2.Field1",
					Unmarshaler: unmarshal.UnmarshalerName{
						Package: pkgUnmarshal,
						Type:    stringUnmarshaler,
					},
				},
				Condition: &reflector.ConditionRequireIf{
					Target: &reflector.TargetField{
						Name: "Field3",
						Unmarshaler: unmarshal.UnmarshalerName{
							Package: pkgUnmarshal,
							Type:    intUnmarshaler,
						},
					},
					ValueStr: "7",
				},
			},
			"SUB2_SUB2FIELD2": {
				Field: &reflector.TargetField{
					Name: "Subconfig2.Field2",
					Unmarshaler: unmarshal.UnmarshalerName{
						Package: pkgUnmarshal,
						Type:    stringUnmarshaler,
					},
				},
				Condition: reflector.ConditionRequired(false),
			},
			"SUB2_SUB2FIELD3": {
				Field: &reflector.TargetField{
					Name: "Subconfig2.Field3",
					Unmarshaler: unmarshal.UnmarshalerName{
						Package: pkgUnmarshal,
						Type:    stringUnmarshaler,
					},
				},
				Condition: &reflector.ConditionRequireIfCombined{
					First: &reflector.ConditionRequireIf{
						Target: &reflector.TargetField{
							Name: "Field3",
							Unmarshaler: unmarshal.UnmarshalerName{
								Package: pkgUnmarshal,
								Type:    intUnmarshaler,
							},
						},
						ValueStr: "7",
					},
					Second: reflector.ConditionEnum([]string{"foo", "bar", "baz"}),
				},
			},
			"SUB3": {
				Field: &reflector.TargetField{
					Name: "Subconfig3",
					Unmarshaler: unmarshal.UnmarshalerName{
						Package: pkgTest,
						Type:    pkg1Unmarshaler,
					},
				},
				Condition: reflector.ConditionRequired(true),
			},
		},
	}

	conf := new(TestReflectConfig)
	result, err := reflector.New().Reflect(conf)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReflector_CombineConditions(t *testing.T) {
	type L3 struct {
		Str string `env:"STR"`
	}

	type L2 struct {
		L3 *L3 `env:"L3_*"`
	}

	type L1 struct {
		Foo string `env:"FOO,enum(bar,baz)"`
		L2  *L2    `env:"L2_*,reqif(Foo=bar)"`
	}

	expect := &reflector.Result{
		ConfigPkg:  "reflector_test",
		ConfigType: "L1",
		Packages: map[string]struct{}{
			"github.com/albenik/huenv/unmarshal": {},
		},
		Envs: map[string]*reflector.Target{
			"FOO": {
				Field: &reflector.TargetField{
					Name: "Foo",
					Unmarshaler: unmarshal.UnmarshalerName{
						Package: "github.com/albenik/huenv/unmarshal",
						Type:    "String",
					},
				},
				Condition: reflector.ConditionEnum{"bar", "baz"},
			},
			"L2_L3_STR": {
				Field: &reflector.TargetField{
					Name: "L2.L3.Str",
					Unmarshaler: unmarshal.UnmarshalerName{
						Package: "github.com/albenik/huenv/unmarshal",
						Type:    "String",
					},
				},
				Condition: &reflector.ConditionRequireIf{
					Target: &reflector.TargetField{
						Name: "Foo",
						Unmarshaler: unmarshal.UnmarshalerName{
							Package: "github.com/albenik/huenv/unmarshal",
							Type:    "String",
						},
					},
					ValueStr: "bar",
				},
			},
		},
	}

	res, err := reflector.New().Reflect(new(L1))
	require.NoError(t, err)
	require.Equal(t, expect, res)
}
