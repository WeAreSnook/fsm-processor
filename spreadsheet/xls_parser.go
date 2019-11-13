package spreadsheet

import (
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
}

// XlsRow represents a row in an Xls sheet
type XlsRow struct {
	row *xls.Row
}

// NewXlsParser creates an XlsParser from a given file path
func NewXlsParser(path string) XlsParser {
	workbook, closer, err := xls.OpenWithCloser(path, "utf-8")

	if err != nil {
		// TODO FailWithError() that prints json to stdout
		log.Fatal(err)
	}

	sheet := workbook.GetSheet(0)
	if sheet == nil {
		log.Fatalf("Couldn't open sheet in xls: %s", path)
	}

	return XlsParser{
		workbook:   workbook,
		sheet:      sheet,
		closer:     closer,
		currentRow: 0,
	}
}

// Next returns the next XlsRow from the sheet
func (p *XlsParser) Next() (XlsRow, error) {
	nextRow := p.currentRow + 1
	if nextRow > int(p.sheet.MaxRow) {
		return XlsRow{}, ErrEOF
	}

	p.currentRow = nextRow
	row := p.sheet.Row(nextRow)

	return XlsRow{row}, nil
}

// Close closes the spreadsheet, making it unavailable for further operations
func (p *XlsParser) Close() {
	p.closer.Close()
}

// Col returns the string in the specified column
func (r *XlsRow) Col(index int) string {
	return r.Col(index)
}
