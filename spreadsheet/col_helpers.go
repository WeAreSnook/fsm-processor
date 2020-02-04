package spreadsheet

import (
	"strconv"

	"github.com/addjam/fsm-processor/llog"
)

// ColByName returns the string in the cell at the specified column
func ColByName(r Row, name string) string {
	index := indexOf(r.Headers(), name)
	if index < 0 {
		return ""
	}

	return r.Col(index)
}

// FloatColByName returns the float32 in the cell at the specified column
func FloatColByName(r Row, name string) float32 {
	str := ColByName(r, name)

	if str == "" {
		return 0
	}

	value, err := strconv.ParseFloat(str, 32)
	if err != nil {
		llog.Printf(`Error parsing float for column "%s", falling back to 0`, name)
		return 0
	}

	return float32(value)
}
