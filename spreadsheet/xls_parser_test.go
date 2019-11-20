package spreadsheet

import "testing"

func TestXlsNext(t *testing.T) {
	t.Run("Allows retrieval of column data by name", func(t *testing.T) {
		parser, err := NewXlsParser(ParserInput{Path: "./testdata/Consent Report W360.xls"})
		if err != nil {
			t.Fatalf("Error creating parser")
		}

		parser.SetHeaderNames([]string{"DocDesc", "DocDate", "CLAIMREFERENCE"})

		row, err := parser.Next()
		if err != nil {
			t.Fatalf("Got an unexpected error %#v", err)
		}

		AssertColumnNamed(t, row, "DocDesc", "FSM Application")
		AssertColumnNamed(t, row, "CLAIMREFERENCE", "000017")
		AssertColumnNamed(t, row, "DocDate", "12/12/18 1:23:45")
	})
}
