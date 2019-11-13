package spreadsheet

import (
	"errors"
	"io"
)

var (
	// ErrEOF is returned when the parser has reached the end of the file
	ErrEOF = io.EOF

	// ErrColumnOutOfBounds is returned when the given column index is too low or too high
	ErrColumnOutOfBounds = errors.New("Given column is out of bounds")
)

// Parser is an interface for types that can parse a spreadsheet by Row
type Parser interface {
	Next() (Row, error)
	Close()
}

// Row refers to a row in a spreadsheet, which has many columns
type Row interface {
	Col(int) string
}
