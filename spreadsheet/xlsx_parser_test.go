package spreadsheet

import "testing"

func TestXlsxNext(t *testing.T) {
	t.Run("Allows retrieval of column data by name", func(t *testing.T) {
		parser := NewXlsxParser("./testdata/Consent Report W360.xlsx", true)

		row, err := parser.Next()
		if err != nil {
			t.Fatalf("Got an unexpected error %#v", err)
		}

		AssertColumnNamed(t, row, "DocDesc", "FSM Application")
		AssertColumnNamed(t, row, "DocDate", "12/12/18 1:23:45")
		AssertColumnNamed(t, row, "CLAIMREFERENCE", "000017")
	})
}
