package spreadsheet

import "testing"

// AssertColumnNamed assets the column with the given name matches the expected output
func AssertColumnNamed(t *testing.T, row Row, name, want string) {
	t.Helper()

	got := row.ColByName(name)

	if got != want {
		t.Fatalf(`Expected "%s" but got "%s"`, want, got)
	}
}
