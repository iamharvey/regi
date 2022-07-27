package io

import (
	"io"
)

// Streams defines io options for cli.
type Streams struct {
	// In, os.Stdin
	In io.Reader

	// Out, os.Stdout
	Out io.Writer

	// ErrOut, os.Stderr
	ErrOut io.Writer
}
