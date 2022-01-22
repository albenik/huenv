package generator

import (
	"io"
)

func FormatSource(dst io.Writer, src io.Reader) error {
	return formatSource(dst, src)
}
