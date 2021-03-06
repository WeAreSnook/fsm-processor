package spreadsheet

import "testing"

func TestCsvNext(t *testing.T) {
	t.Run("Parses space sparated files correctly", func(t *testing.T) {
		input := ParserInput{
			Path: "./testdata/space separated.txt",
		}
		parser, err := NewCsvParser(input)
		if err != nil {
			t.Fatalf("Error creating parser")
		}

		parser.SetSeparator(' ')

		row, err := parser.Next()
		if err != nil {
			t.Fatalf("Got an unexpected error %#v", err)
		}

		number := row.Col(0)
		description := row.Col(1)
		date := row.Col(2)

		if number != "1" {
			t.Fatalf(`Expected "1" but got "%s"`, number)
		}

		if description != "this is christmas" {
			t.Fatalf(`Expected "this is christmas" but got "%s"`, description)
		}

		if date != "25/12/2019" {
			t.Fatalf(`Expected "25/12/2019" but got "%s"`, date)
		}
	})

	t.Run("Allows retrieval of column data by name", func(t *testing.T) {
		parser, err := NewCsvParser(ParserInput{Path: "./testdata/csv with headers.txt", HasHeaders: true})

		if err != nil {
			t.Fatalf("Error creating parser")
		}

		row, err := parser.Next()
		if err != nil {
			t.Fatalf("Got an unexpected error %#v", err)
		}

		AssertColumnNamed(t, row, "ID", "1")
		AssertColumnNamed(t, row, "description", "this is christmas")
		AssertColumnNamed(t, row, "Date", "25/12/2019")
	})
}
