package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/addjam/fsm-processor/spreadsheet"
)

// AddPeopleWithConsent parses which people have given consent to check entitlement data
// and adds them directly to the PeopleStore
// Data sources: Consent 360 & Benefit Extract
func AddPeopleWithConsent(inputData InputData, peopleStore *PeopleStore) {
	consentByClaimNumber := extractConsentData(inputData)

	// Parse benefits extract
	spreadsheet.EachRow(inputData.benefitExtract, func(row spreadsheet.Row) {
		claimNumber, err := strconv.Atoi(row.Col(0))

		if err != nil {
			log.Printf("Error parsing claim number from benefits extract %s", row.Col(0))
			return
		}

		hasPermission := consentByClaimNumber[claimNumber]

		if hasPermission {
			peopleStore.Add(
				Person{
					Forename:          row.Col(4),
					Surname:           row.Col(3),
					ClaimNumber:       claimNumber,
					AgeYears:          0,
					BenefitExtractRow: row,
				},
			)
		}
	})
}

func extractConsentData(inputData InputData) map[int]bool {
	consentData := make(map[int]bool)

	spreadsheet.EachRow(inputData.consent360, func(row spreadsheet.Row) {
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
	})

	return consentData
}
