package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/addjam/fsm-processor/llog"
)

// RespondWith stops execution and outputs response data as json
//
// fsmStore - PeopleStore representing the final state of the FSM algorithm data
// ctrStore - PeopleStore representing the final state of the FSM algorithm data
// err - optional error that halted execution
func RespondWith(fsmStore *PeopleStore, ctrStore *PeopleStore, err error) {
	output := Output{
		Success: err == nil,
	}

	if err != nil {
		output.Error = err.Error()
	}

	output.FsmDebugData = generateDebugData(fsmStore)
	output.CtrDebugData = generateDebugData(ctrStore)
	output.Log = llog.Data()

	// Output as json
	json, err := json.Marshal(output)
	if err != nil {
		log.Fatal(`{ "success": false, "error": "Error marshalling json from store" }`)
	}

	fmt.Print(string(json))

	if !output.Success {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

// Output represents the result data
type Output struct {
	Success      bool   `json:"success"`
	FsmDebugData string `json:"fsm_debug,omitempty"`
	CtrDebugData string `json:"ctr_debug,omitempty"`
	Error        string `json:"error,omitempty"`
	Log          string `json:"log"`
}

func generateDebugData(store *PeopleStore) string {
	if store == nil {
		return "No data"
	}

	return fmt.Sprintf(`
		%d people in store,
	`,
		len(store.People))
}
