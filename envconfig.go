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

	envmap := envconf.Envmap()
	for name := range info.Envs {
		if _, ok = envmap[name]; !ok {
			return ErrOutdated
		}
	}

	keys := make(sort.StringSlice, 0, len(envmap))
	for k := range envmap {
		keys = append(keys, k)
	}
	keys.Sort()

	var outerr error

	for _, key := range keys {
		if err = envmap[key].Unmarshal(os.Getenv(envprefix + key)); err != nil {
			outerr = multierr.Append(outerr, &KeyError{key, err})
		}
	}

	return outerr
}
