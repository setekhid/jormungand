package proto

import (
	"io"
)

type Decoder func(io.Reader) io.Reader
type Encoder func(io.Writer) io.Writer
