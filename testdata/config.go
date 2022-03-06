package testdata

type TestConfig1 struct {
	Field1       string   `env:"X_TEST_FIELD_1"`
	StringsSlice []string `env:"X_TEST_STR_SLICE"`
}
