package spreadsheet

import (
	"encoding/csv"
	"log"
)

// EachParserRow calls func for each of the rows provided by a Parser
// Automatically closes the parser
func EachParserRow(p Parser, f func(Row)) {
	defer p.Close()

	for {
		row, err := p.Next()

		if err == ErrEOF {
			break
		} else if err != nil {
			if err, ok := err.(*csv.ParseError); ok && err.Err == csv.ErrFieldCount {
				// Some files have extra data which presents as a field count error, e.g. hb-uc.d has some trailing data.
				break
			}

			log.Fatal(err)
		}

		f(row)
	}
}

// EachRow takes the path of a spreadsheet and executes the func once for each row
func EachRow(input ParserInput, f func(Row)) error {
	parser, err := NewParser(input)
	if err != nil {
		return err
	}

	EachParserRow(parser, f)

	return nil
}

// CountRows returns the number of rows in a spreadsheet
func CountRows(input ParserInput) int {
	total := 0
	EachRow(input, func(r Row) {
		total++
	})
	return total
}

// CreateIndex returns a map of rowKey => []Row. rowKey is created by the keyCreator function, which takes a cell value and returns a rowKey
func CreateIndex(i ParserInput, colName string, rowKeyCreator func(string) string) (map[string][]Row, error) {
	index := make(map[string][]Row)

	err := EachRow(i, func(r Row) {
		baseKey := ColByName(r, colName)
		key := rowKeyCreator(baseKey)
		index[key] = append(index[key], r)
	})

	return index, err
}

// AssertHeadersExist ensures the provided headers exist and exits if they don't
func AssertHeadersExist(p Parser, expectedHeaders []string) error {
	for _, hdr := range expectedHeaders {
		if indexOf(p.Headers(), hdr) < 0 {
			return ErrMissingHeader{filePath: p.Path(), header: hdr}
		}
	}

	return nil
}

func indexOf(in []string, target string) int {
	for index, header := range in {
		if header == target {
			return index
		}
	}
	return -1
}
