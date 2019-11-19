package people

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// RespondWith stops execution and outputs some json
func RespondWith(store Store, err error) {
	// TODO update output data when we have some output to return
	output := Output{
		Success:        err == nil,
		OutputFilePath: "none yet",
	}

	if err != nil {
		output.Error = err.Error()
	} else {
		output.DebugData = fmt.Sprintf("%d people extracted", len(store.People))
	}

	// Format json
	json, err := json.Marshal(output)
	if err != nil {
		log.Fatal("Error marshalling json from store")
	}

	fmt.Println(string(json))

	if !output.Success {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

// Output represents the result data
type Output struct {
	Success        bool   `json:"success"`
	OutputFilePath string `json:"output_file_path"`
	DebugData      string `json:"debug,omitempty"`
	Error          string `json:"error,omitempty"`
}
