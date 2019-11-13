package spreadsheet

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
)

// CsvParser is a Parser implementation that handles CSVs
type CsvParser struct {
	file      *os.File
	csvReader *csv.Reader
}

// CsvRow represents a row in a CSV file
type CsvRow struct {
	line []string
}

// NewCsvParser creates a CsvParser with the given path, opening the file and preparing it for reading
func NewCsvParser(path string) *CsvParser {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %s", path)
	}
	fileReader := bufio.NewReader(file)
	csvReader := csv.NewReader(fileReader)

	// Skip column header row
	_, err = csvReader.Read()

	if err != nil {
		log.Fatal(err)
	}

	return &CsvParser{
		file:      file,
		csvReader: csvReader,
	}
}

// Next returns the next Row from the file, or errors if for example we reached the end
func (p *CsvParser) Next() (Row, error) {
	line, err := p.csvReader.Read()

	if err != nil {
		return CsvRow{}, err
	}

	row := CsvRow{line}
	return row, err
}

// Close closes the CSV file. No further operations will be possible.
func (p CsvParser) Close() {
	p.file.Close()
}

// Col returns the string at the specified index from the CsvRow
func (r CsvRow) Col(index int) string {
	if index < 0 || index > len(r.line)-1 {
		return ""
	}

	return r.line[index]
}
