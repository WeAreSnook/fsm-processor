package spreadsheet

import (
	"fmt"
	"io"
	"log"

	"github.com/extrame/xls"
)

// XlsParser is a Parser implementation that handles xls spreadsheets
type XlsParser struct {
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
func NewXlsParser(path string) *XlsParser {
	workbook, closer, err := xls.OpenWithCloser(path, "utf-8")

	if err != nil {
		// TODO FailWithError() that prints json to stdout
		log.Fatal(err)
	}

	sheet := workbook.GetSheet(0)
	if sheet == nil {
		log.Fatalf("Couldn't open sheet in xls: %s", path)
	}

	return &XlsParser{
		workbook:   workbook,
		sheet:      sheet,
		closer:     closer,
		currentRow: 0,
		hasHeaders: false,
	}
}

// Next returns the next Row from the sheet
func (p *XlsParser) Next() (Row, error) {
	nextRow := p.currentRow + 1
	if nextRow > int(p.sheet.MaxRow) {
		return XlsRow{}, ErrEOF
	}

	p.currentRow = nextRow
	row := p.sheet.Row(nextRow)
	return XlsRow{p, row}, nil
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

// Col returns the string in the specified column
func (r XlsRow) Col(index int) string {
	return r.row.Col(index)
}

// ColByName returns the string in the cell at the specified column
//
// NOTE: SetHeaderNames should be called to enable this, as Xls headers aren't
// automatically parsed
func (r XlsRow) ColByName(name string) string {
	if !r.p.hasHeaders {
		return ""
	}

	index := indexOf(r.p.headers, name)
	fmt.Printf("Index is %d\n", index)
	if index < 0 {
		return ""
	}

	return r.Col(index)
}
