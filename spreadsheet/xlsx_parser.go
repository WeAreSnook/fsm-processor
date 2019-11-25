package spreadsheet

import (
	"github.com/tealeg/xlsx"
)

// XlsxParser is a spreadsheet.Parser implementation for Xlsx files
type XlsxParser struct {
	path       string
	file       *xlsx.File
	sheet      *xlsx.Sheet
	currentRow int
	numRows    int
	hasHeaders bool
	headers    []string
}

// XlsxRow is a spreadsheet.Row implementation for Xlsx files
type XlsxRow struct {
	p   *XlsxParser
	row *xlsx.Row
}

// NewXlsxParser returns an XlsxParser for the file at the given path
func NewXlsxParser(input ParserInput) (*XlsxParser, error) {
	xlFile, err := xlsx.OpenFile(input.Path)
	if err != nil {
		return nil, ErrUnableToParse{input.Path}
	}

	sheet := xlFile.Sheets[0]

	if sheet == nil {
		return nil, ErrUnableToParse{input.Path}
	}

	var headers []string
	if input.HasHeaders {
		headerRow := sheet.Row(0)
		for _, cell := range headerRow.Cells {
			headers = append(headers, cell.String())
		}
	}

	parser := &XlsxParser{
		path:       input.Path,
		file:       xlFile,
		sheet:      sheet,
		currentRow: 0,
		numRows:    sheet.MaxRow,
		hasHeaders: input.HasHeaders,
		headers:    headers,
	}

	err = AssertHeadersExist(parser, input.RequiredHeaders)

	return parser, err
}

// Next returns the next Row
func (p *XlsxParser) Next() (Row, error) {
	nextRow := p.currentRow + 1
	if nextRow > p.numRows {
		return XlsxRow{}, ErrEOF
	}

	p.currentRow = nextRow
	row := p.sheet.Row(nextRow)

	return XlsxRow{p: p, row: row}, nil
}

// Close is unimplemented and unnecessary for xlsx files
func (p XlsxParser) Close() {
	// Handled automatically by xlsx library
}

// SetHeaderNames sets header names, allowing retrieval of columns by name
func (p *XlsxParser) SetHeaderNames(names []string) {
	p.headers = names
	p.hasHeaders = true
}

// Headers returns the headers found or set on the current parsed file
func (p XlsxParser) Headers() []string {
	return p.headers
}

// Path returns the path used for the file being parsed
func (p XlsxParser) Path() string {
	return p.path
}

// Col returns the string in the specified column
func (r XlsxRow) Col(index int) string {
	if index < 0 || index > len(r.row.Cells)-1 {
		return ""
	}

	cell := r.row.Cells[index]
	if cell == nil {
		return ""
	}

	return cell.String()
}

// Headers returns the headers from the XLSX
func (r XlsxRow) Headers() []string {
	return r.p.headers
}
