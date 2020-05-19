# spreadsheet

## Functions

### func [AssertColumnNamed](/testing.go#L8)

`func AssertColumnNamed(t *testing.T, row Row, name, want string)`

AssertColumnNamed assets the column with the given name matches the expected output

### func [AssertHeadersExist](/helpers.go#L65)

`func AssertHeadersExist(p Parser, expectedHeaders []string) error`

AssertHeadersExist ensures the provided headers exist and exits if they don't

### func [ColByName](/col_helpers.go#L10)

`func ColByName(r Row, name string) string`

ColByName returns the string in the cell at the specified column

### func [CountRows](/helpers.go#L43)

`func CountRows(input ParserInput) int`

CountRows returns the number of rows in a spreadsheet

### func [CreateIndex](/helpers.go#L52)

`func CreateIndex(i ParserInput, colName string, rowKeyCreator func(string) string) (map[string][]Row, error)`

CreateIndex returns a map of rowKey => []Row. rowKey is created by the keyCreator function, which takes a cell value and returns a rowKey

### func [EachParserRow](/helpers.go#L9)

`func EachParserRow(p Parser, f func(Row)) error`

EachParserRow calls func for each of the rows provided by a Parser
Automatically closes the parser

### func [EachRow](/helpers.go#L33)

`func EachRow(input ParserInput, f func(Row)) error`

EachRow takes the path of a spreadsheet and executes the func once for each row

### func [FloatColByName](/col_helpers.go#L20)

`func FloatColByName(r Row, name string) float32`

FloatColByName returns the float32 in the cell at the specified column
