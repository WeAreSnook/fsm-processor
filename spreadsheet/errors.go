package spreadsheet

import (
	"errors"
	"fmt"
	"io"
)

var (
	// ErrEOF is returned when the parser has reached the end of the file
	ErrEOF = io.EOF

	// ErrColumnOutOfBounds is returned when the given column index is too low or too high
	ErrColumnOutOfBounds = errors.New("Given column is out of bounds")
)

// ErrMissingHeader is an error message representing which header in which file was missing
type ErrMissingHeader struct {
	filePath string
	header   string
}

func (e ErrMissingHeader) Error() string {
	return fmt.Sprintf(`Expected file "%s" to have header "%s" but it was missing`, e.filePath, e.header)
}
