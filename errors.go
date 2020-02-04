package main

import "fmt"

// ErrInvalidInputPath represents an error opening/parsing the given file
type ErrInvalidInputPath struct {
	filePath string
}

func (e ErrInvalidInputPath) Error() string {
	return fmt.Sprintf(`Invalid input path "%s"`, e.filePath)
}
