package spreadsheet

import (
	"io"

	"github.com/extrame/xls"
)

// XlsParser is a Parser implementation that handles xls spreadsheets
type XlsParser struct {
	path       string
	workbook   *xls.WorkBook
	sheet      *xls.WorkSheet
	closer     io.Closer
	currentRow int
	headers    []string
	hasHeaders bool
}

// XlsRow represents a row in an Xls sheet
type XlsRow struct {
	p   *XlsParser
	row *xls.Row
}

// NewXlsParser creates an XlsParser from a given file path
func NewXlsParser(input ParserInput) (*XlsParser, error) {
	workbook, closer, err := xls.OpenWithCloser(input.Path, "utf-8")
	if err != nil {
		return nil, err
	}

	sheet := workbook.GetSheet(0)
	if sheet == nil {
		return nil, ErrUnableToParse{input.Path}
	}

	var headers []string
	if input.HasHeaders {
		headerRow := sheet.Row(0)
		for i := 0; i <= headerRow.LastCol(); i++ {
			cell := headerRow.Col(i)
			headers = append(headers, cell)
		}
	}

	parser := &XlsParser{
		path:       input.Path,
		workbook:   workbook,
		sheet:      sheet,
		closer:     closer,
		currentRow: 0,
		hasHeaders: input.HasHeaders,
		headers:    headers,
	}

	err = AssertHeadersExist(parser, input.RequiredHeaders)

	return parser, err
}

// Next returns the next Row from the sheet
func (p *XlsParser) Next() (Row, error) {
	nextRow := p.currentRow + 1
	if nextRow > int(p.sheet.MaxRow) {
		return XlsRow{}, ErrEOF
	}

	p.currentRow = nextRow
	row := p.sheet.Row(nextRow)
	return XlsRow{p: p, row: row}, nil
}

// Close closes the spreadsheet, making it unavailable for further operations
func (p XlsParser) Close() {
	p.closer.Close()
}

// SetHeaderNames sets header names, allowing retrieval of columns by name
func (p *XlsParser) SetHeaderNames(names []string) {
	p.headers = names
	p.hasHeaders = true
}

// Headers returns the headers found or set on the current parsed file
func (p XlsParser) Headers() []string {
	return p.headers
}

// Path returns the path used for the file being parsed
func (p XlsParser) Path() string {
	return p.path
}

// Col returns the string in the specified column
func (r XlsRow) Col(index int) string {
	return r.row.Col(index)
}

// Headers returns the headers from the XLS
func (r XlsRow) Headers() []string {
	return r.p.headers
}
