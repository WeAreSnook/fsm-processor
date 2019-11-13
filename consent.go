package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/addjam/fsm-processor/spreadsheet"
	"github.com/extrame/xls"
)

// ExtractPeopleWithConsent parses which people have given consent to check entitlement data
func ExtractPeopleWithConsent(inputData InputData, peopleStore *PeopleStore) {
	consentByClaimNumber := extractConsentData(inputData)

	// Parse benefits extract
	benefitExtractParser := spreadsheet.NewCsvParser(inputData.benefitExtractPath)

	// Parse files
	for {
		line, err := benefitExtractParser.Next()

		if err == spreadsheet.ErrEOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		claimNumber, err := strconv.Atoi(line.Col(0))

		if err != nil {
			log.Printf("Error parsing claim number from benefits extract %s", line.Col(0))
			continue
		}

		hasPermission := consentByClaimNumber[claimNumber]

		if hasPermission {
			peopleStore.Add(
				Person{
					forename:    line.Col(4),
					surname:     line.Col(3),
					claimNumber: line.Col(0),
					ageYears:    0,
				},
			)
		}
	}
}

func extractConsentData(inputData InputData) map[int]bool {
	// Consent spreadsheet is used to determine which people on the benefit extract have given
	// permission
	consentFile, err := xls.Open(inputData.consent360Path, "utf-8")
	if err != nil {
		// TODO FailWithError() that prints json to stdout
		log.Fatal(err)
	}

	consentSheet := consentFile.GetSheet(0)
	if consentSheet == nil {
		log.Fatal("Couldn't open sheet in consnet file")
	}

	// Consent data mapping claim number to entitled bool
	consentData := make(map[int]bool)
	for i := 0; i <= int(consentSheet.MaxRow); i++ {
		row := consentSheet.Row(i)
		claimNumStr := strings.Replace(row.Col(2), "TEMP", "", 1)
		claimNum, err := strconv.Atoi(claimNumStr)

		if err != nil {
			// consent spreadsheet has claim numbers beginning with "TEMP" followed by 6 digits
			// benefit extract just seems to be numbers. Consent spreadsheet can also have e.g. 000123 but
			// benefit extract seems to present this as 123
			// fmt.Printf("Error parsing claim number %s", row.Col(2))
		}

		consentDesc := row.Col(0)
		hasPermission := consentDesc != "FSM&CG Consent Removed" // TODO this isn't in our example spreadsheet
		consentData[claimNum] = hasPermission
	}

	return consentData
}
