package spreadsheet

import (
	"errors"
	"io"
	"log"
	"path/filepath"
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

// NewParser creates a parser appropriate for the spreadsheet at the given path.
// Supports:
//   - CSV
//   - xls
//   - xlsx
func NewParser(path string) Parser {
	extension := filepath.Ext(path)

	switch extension {
	case ".xls":
		return NewXlsParser(path)
	case ".xlsx":
		return NewXlsxParser(path)
	case ".txt":
		return NewCsvParser(path)
	case ".csv":
		return NewCsvParser(path)
	}

	log.Fatalf("No parser for extension %s of file at path %s\n", extension, path)

	return nil
}
