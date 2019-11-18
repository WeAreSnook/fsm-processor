package spreadsheet

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
)

// CsvParser is a Parser implementation that handles CSVs
type CsvParser struct {
	file       *os.File
	csvReader  *csv.Reader
	headers    []string
	hasHeaders bool
}

// CsvRow represents a row in a CSV file
type CsvRow struct {
	p    *CsvParser
	line []string
}

// NewCsvParser creates a CsvParser with the given path, opening the file and preparing it for reading
func NewCsvParser(path string, hasHeaders bool) *CsvParser {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %s", path)
	}
	fileReader := bufio.NewReader(file)
	csvReader := csv.NewReader(fileReader)

	var headers []string
	if hasHeaders {
		// Skip column header row
		line, err := csvReader.Read()
		if err != nil {
			log.Fatal(err)
		}

		headers = line
	}

	return &CsvParser{
		file:       file,
		csvReader:  csvReader,
		headers:    headers,
		hasHeaders: hasHeaders,
	}
}

// Next returns the next Row from the file, or errors if for example we reached the end
func (p *CsvParser) Next() (Row, error) {
	line, err := p.csvReader.Read()

	if err != nil {
		return CsvRow{}, err
	}

	row := CsvRow{p, line}
	return row, err
}

// Close closes the CSV file. No further operations will be possible.
func (p CsvParser) Close() {
	p.file.Close()
}

// SetHeaderNames sets header names, allowing retrieval of columns by name
func (p *CsvParser) SetHeaderNames(names []string) {
	p.headers = names
	p.hasHeaders = true
}

// SetSeparator changes the delimiter parsed in the provided file. Default is a comma.
func (p CsvParser) SetSeparator(r rune) {
	p.csvReader.Comma = r
}

// Col returns the string at the specified index from the CsvRow
func (r CsvRow) Col(index int) string {
	if index < 0 || index > len(r.line)-1 {
		return ""
	}

	return r.line[index]
}

// ColByName returns the string in the cell at the specified column
func (r CsvRow) ColByName(name string) string {
	index := indexOf(r.p.headers, name)
	if index < 0 {
		return ""
	}

	return r.Col(index)
}
