package huenv

import (
	"fmt"
	"os"
	"sort"

	"go.uber.org/multierr"

	"github.com/albenik/huenv/reflector"
	"github.com/albenik/huenv/unmarshal"
)

type Config interface {
	Envmap() map[string]*unmarshal.Target
}

func Init(conf interface{}, envprefix string) error {
	envconf, ok := conf.(Config)
	if !ok {
		return ErrNotGenerated
	}

	info, err := reflector.New().Reflect(conf)
	if err != nil {
		return fmt.Errorf("reflect: %w", err)
	}

	envs := envconf.Envmap()
	envKeys := make(sort.StringSlice, 0, len(envs))
	for k := range envs {
		envKeys = append(envKeys, k)
	}

	envKeys.Sort()
	envKeys = reflector.SortWithDeps(envKeys, info.Envs)

	for _, k := range envKeys {
		if _, ok = envs[k]; !ok {
			return ErrOutdated
		}
	}

	var outerr error
	for _, k := range envKeys {
		if err = envs[k].Unmarshal(os.Getenv(envprefix + k)); err != nil {
			outerr = multierr.Append(outerr, &KeyError{k, err})
		}
	}

	return outerr
}
