package spreadsheet

import (
	"path/filepath"
)

// Parser is an interface for types that can parse a spreadsheet by Row
type Parser interface {
	Next() (Row, error)
	Close()
	SetHeaderNames([]string)
	Headers() []string
	Path() string
}

// Row refers to a row in a spreadsheet, which has many columns
type Row interface {
	Headers() []string
	Col(int) string
}

// Format represents the format of the spreadsheete, e.g. xls, csv, etc
type Format int

const (
	// Auto to auto-detect the format. Based on extension and supports xls, xlsx and falls back to csv
	Auto = iota

	// Csv is CSV format
	Csv Format = iota

	// Ssv is space spearated
	Ssv = iota

	// Xls is Xls excel files, from excel up to 2004
	Xls = iota

	// Xlsx is a modern xlsx excel file
	Xlsx = iota
)

//ParserInput represents a spreadsheet and associated options/validations
type ParserInput struct {
	Path            string
	HasHeaders      bool
	RequiredHeaders []string
	Format          Format
}

// NewParser creates a parser appropriate for the spreadsheet at the given path.
// Supports:
//   - CSV
//   - xls
//   - xlsx
func NewParser(input ParserInput) (Parser, error) {
	inputFormat := input.Format
	if inputFormat == Auto {
		extension := filepath.Ext(input.Path)
		inputFormat = formatFromExtension(extension)
	}

	switch inputFormat {
	case Xls:
		return NewXlsParser(input)
	case Xlsx:
		return NewXlsxParser(input)
	case Csv:
		return NewCsvParser(input)
	case Ssv:
		return NewCsvParser(input)
	}

	return nil, nil
}

func formatFromExtension(ext string) Format {
	switch ext {
	case ".xls":
		return Xls
	case ".xlsx":
		return Xlsx
	case ".txt":
		return Csv
	case ".csv":
		return Csv
	}

	return Csv
}
