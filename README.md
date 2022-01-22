# EnvConfig for Humans & Robots

[![CI](https://github.com/albenik/huenv/actions/workflows/main.yml/badge.svg)](https://github.com/albenik/huenv/actions/workflows/main.yml)

The main goal — is the clear information about all configuration in one place

# Usage

```
config.go:

package config

//go:generate go run github.com/albenik/huenv/cmd/huenv -out confg_unmarshal.go example.com/project/name/config Config

type Config struct {
  ...
}
```
