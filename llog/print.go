package llog

import "fmt"

var (
	logData string

	// PrintToStdout determines if log output is written in real time to stdout
	PrintToStdout = false
)

// Print data
func Print(data string) {
	if PrintToStdout {
		fmt.Print(data)
	}

	logData += data
}

// Println prints data followed by a newline
func Println(data string) {
	Printf("%s\n", data)
}

// Printf mirrors llog.Printf( params. Logs data
func Printf(format string, a ...interface{}) {
	output := fmt.Sprintf(format, a...)
	Print(output)
}

// Data returns all logged data so far
func Data() string {
	return logData
}
