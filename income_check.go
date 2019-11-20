package main

import (
	"fmt"
	"strconv"

	"github.com/addjam/fsm-processor/spreadsheet"
)

const ctcWtcAnnualIncomeFigure = 6420

// Required sheets:
// - Benefit Extract
// - dependants shbe
// - school roll
// - current awards
// - hb uc d (universal credit)
// - eligibility gap list

// incomeData represents the data for a single person
type incomeData struct {
	person                 Person
	taxCreditIncomeStepOne float32
	taxCreditIncomeStepTwo float32
	taxCreditFigure        float32
}

// PeopleWithQualifyingIncomes returns just the people in the provided store that qualify
// for FSM or CG. Updates the people to show this.
func PeopleWithQualifyingIncomes(inputData InputData, store PeopleStore) []Person {
	for _, person := range store.People {
		// Calculate step one/two data
		incomeData := calculateIncomeSteps(person)

		// Calculate tax credit figure
		if incomeData.taxCreditIncomeStepOne <= 300 {
			incomeData.taxCreditFigure = incomeData.taxCreditIncomeStepTwo * 52
		} else {
			incomeData.taxCreditFigure = (incomeData.taxCreditIncomeStepTwo - 300) * 52
		}
	}

	return []Person{}
}

func calculateIncomeSteps(person Person) incomeData {
	return incomeData{
		person:                 person,
		taxCreditIncomeStepOne: calculateStepOne(person),
		taxCreditIncomeStepTwo: calculateStepTwo(person),
	}
}

func calculateStepOne(person Person) float32 {
	colNames := []string{
		"Clmt Personal Pension",
		"Clmt State Retirement Pension (incl SERP's graduated pension etc)",
		"Ptnr Personal Pension",
		"Ptnr State Retirement Pension (incl SERP's graduated pension etc)",
		"Clmt Occupational Pension",
		"Ptnr Occupational Pension",
	}
	return sumFloatColumns(person.BenefitExtractRow, colNames)
}

func calculateStepTwo(person Person) float32 {
	colNames := []string{
		"Clmt AIF",
		"Clmt Employment (gross)",
		"Clmt Self-employment (gross)",
		"Clmt Student Grant/Loan",
		"Clmt Sub-tenants",
		"Clmt Boarders",
		"Clmt Government Training",
		"Clmt Statutory Sick Pay",
		"Clmt Widowed Parent's Allowance",
		"Clmt Apprenticeship",
		"Clmt Statutory Sick Pay",
		"Other weekly Income including In-Work Credit",
		"Ptnr AIF",
		"Ptnr Employment (gross)",
		"Ptnr Self-employment (gross)",
		"Ptnr Student Grant/Loan",
		"Ptnr Sub-tenants",
		"Ptnr Boarders",
		"Ptnr Training for Work/Community Action",
		"Ptnr New Deal 50+ Employment Credit",
		"Ptnr Government Training",
		"Ptnr Carer's Allowance",
		"Ptnr Statutory Sick Pay",
		"Ptnr Widowed Parent's Allowance",
		"Ptnr Apprenticeship",
		"Other weekly Income including In-Work Credit",
		"Clmt Savings Credit",
		"Ptnr Savings Credit",
		"Clmt Widows Benefit",
		"Ptnr Widows Benefit",
	}
	return sumFloatColumns(person.BenefitExtractRow, colNames)
}

func sumFloatColumns(row spreadsheet.Row, colNames []string) float32 {
	var result float32 = 0
	for _, colName := range colNames {
		cellStr := row.ColByName(colName)

		// Default to "0" for empty cells
		if cellStr == "" {
			cellStr = "0"
		}

		value, err := strconv.ParseFloat(cellStr, 32)
		if err != nil {
			fmt.Printf(`Error parsing float from cell value "%s" for col name "%s", falling back to 0`, cellStr, colName)
			fmt.Printf("\n")
			value = 0
		}
		result += float32(value)
	}
	return result
}
