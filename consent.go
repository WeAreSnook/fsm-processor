package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/addjam/fsm-processor/llog"
	"github.com/addjam/fsm-processor/spreadsheet"
)

// AddPeopleWithConsent parses which people have given consent to check entitlement data
// and adds them directly to the PeopleStore
// Data sources: Consent 360 & Benefit Extract
func AddPeopleWithConsent(inputData InputData, peopleStore *PeopleStore) error {
	consentDescByClaimNumber, err := extractConsentData(inputData)

	if err != nil {
		return err
	}

	// Parse benefits extract
	numPeople := 0
	err = spreadsheet.EachRow(inputData.benefitExtract, func(row spreadsheet.Row) {
		claimNumStr := spreadsheet.ColByName(row, "Claim Number")
		claimNumber, err := strconv.Atoi(claimNumStr)

		if err != nil {
			log.Printf("Error parsing claim number from benefits extract %s", claimNumStr)
			return
		}

		desc := consentDescByClaimNumber[claimNumber]
		hasPermission := desc != "FSM&CG Consent Removed" && desc != ""
		numPeople += 1

		if hasPermission {
			person, err := NewPersonFromBenefitExtract(row)

			if err != nil {
				llog.Println("Error creating person from benefit extract")
				return
			}

			person.ConsentDesc = desc
			peopleStore.Add(person)
		}
	})

	llog.Printf("%d rows checked for consent\n", numPeople)
	return err
}

func extractConsentData(inputData InputData) (map[int]string, error) {
	consentData := make(map[int]string)

	err := spreadsheet.EachRow(inputData.consent360, func(row spreadsheet.Row) {
		claimNumStr := strings.Replace(row.Col(2), "TEMP", "", 1)
		claimNum, err := strconv.Atoi(claimNumStr)

		if err != nil {
			// consent spreadsheet has claim numbers beginning with "TEMP" followed by 6 digits
			// benefit extract just seems to be numbers. Consent spreadsheet can also have e.g. 000123 but
			// benefit extract seems to present this as 123
			// llog.Printf("Error parsing claim number %s", row.Col(2))
		}

		consentDesc := row.Col(0)
		consentData[claimNum] = consentDesc
	})

	return consentData, err
}
