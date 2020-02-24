package main

import (
	"strconv"
	"sync"

	"github.com/addjam/fsm-processor/llog"
	"github.com/addjam/fsm-processor/spreadsheet"
)

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
	combinedQualifier      bool
	qualifierType          string
}

// PeopleWithQualifyingIncomes returns just the people in the provided store that qualify
// for FSM or CG. Updates the people to show this.
func PeopleWithQualifyingIncomes(inputData InputData, store PeopleStore) ([]Person, error) {
	var people []Person

	universalCreditParser, err := spreadsheet.NewParser(inputData.universalCredit)
	if err != nil {
		return people, err
	}
	universalCreditParser.SetHeaderNames([]string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "aa", "ab", "ac", "ad", "ae",
	})

	var rowsByClaimNum map[int]spreadsheet.Row = make(map[int]spreadsheet.Row)
	err = spreadsheet.EachParserRow(universalCreditParser, func(r spreadsheet.Row) {
		claimNumStr := spreadsheet.ColByName(r, "b")
		claimNum, err := strconv.Atoi(claimNumStr)

		if err != nil {
			llog.Printf(`Error converting "%s"`, claimNumStr)
			llog.Printf("\n%#v\n", r)
		} else {

			rowsByClaimNum[claimNum] = r
		}
	})

	if err != nil {
		return people, err
	}

	var wg sync.WaitGroup
	qualifyingPeopleChan := make(chan Person)

	for _, person := range store.People {
		wg.Add(1)
		ucRow := rowsByClaimNum[person.ClaimNumber]
		go qualifyPerson(inputData, person, ucRow, qualifyingPeopleChan, &wg)
	}

	go func() {
		wg.Wait()
		close(qualifyingPeopleChan)
	}()

	for person := range qualifyingPeopleChan {
		people = append(people, person)
	}

	return people, nil
}

// AddPeopleWithCtr adds people to the store who are receiging a
// weekly cts entitlement greater than 0
func AddPeopleWithCtr(inputData InputData, store *PeopleStore) error {
	return spreadsheet.EachRow(inputData.benefitExtract, func(r spreadsheet.Row) {
		weeklyCtsEntitlement := spreadsheet.FloatColByName(r, "Weekly CTS  entitlement")

		if weeklyCtsEntitlement <= 0.0 {
			return
		}

		person, err := NewPersonFromBenefitExtract(r)

		if err != nil {
			llog.Println("Error creating person from benefit extract")
			return
		}

		store.Add(person)
	})
}

// Concurrently qualifies person based on icnome data
func qualifyPerson(inputData InputData, p Person, universalCreditRow spreadsheet.Row, ch chan Person, w *sync.WaitGroup) {
	defer w.Done()

	// Calculate step one/two data
	incomeData := calculateIncomeSteps(p)

	// Calculate tax credit figure
	if incomeData.taxCreditIncomeStepOne <= 300 {
		incomeData.taxCreditFigure = incomeData.taxCreditIncomeStepTwo * 52
	} else {
		incomeData.taxCreditFigure = (incomeData.taxCreditIncomeStepTwo - 300) * 52
	}

	if incomeData.taxCreditFigure < 0 {
		incomeData.taxCreditFigure = 0
	}

	// Check for FSM & CG combined qualification
	incomeData.combinedQualifier, incomeData.qualifierType = determineCombinedQualifier(inputData, p, incomeData, universalCreditRow)
	if incomeData.combinedQualifier {
		p.QualiferType = incomeData.qualifierType
		for i, d := range p.Dependents {
			d.NewCG = true
			d.NewFSM = true
			d.Person = p
			p.Dependents[i] = d
		}

		ch <- p
		return
	}

	// Check for CG-only qualification via weekly cts entitlement being greater than 0.0
	weeklyCtsEntitlement := spreadsheet.FloatColByName(p.BenefitExtractRow, "Weekly CTS  entitlement")
	if weeklyCtsEntitlement > 0.0 {
		for i, d := range p.Dependents {
			d.NewCG = true
			d.Person = p
			p.Dependents[i] = d
		}

		ch <- p
	}
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

func determineCombinedQualifier(inputData InputData, p Person, incomeData incomeData, universalCreditRow spreadsheet.Row) (bool, string) {
	row := p.BenefitExtractRow

	wtc := spreadsheet.FloatColByName(row, "Clmt Working Tax Credits") + spreadsheet.FloatColByName(row, "Ptnr Working Tax Credits")
	ctc := spreadsheet.FloatColByName(row, "Child tax credit - Claimant") + spreadsheet.FloatColByName(row, "Child tax credit - Partner")

	qualifierA := wtc == 0 && ctc > 0 && incomeData.taxCreditFigure <= inputData.ctcFigure

	qualifierB := wtc > 0 && ctc > 0 && incomeData.taxCreditFigure <= inputData.ctcWtcFigure

	passportedStdClaimIndicator := spreadsheet.ColByName(row, "Passported / Standard claim indicator")
	passportQualifier := passportedStdClaimIndicator == "ESA(IR)" ||
		passportedStdClaimIndicator == "Income Support" ||
		passportedStdClaimIndicator == "JSA(IB)"

	ucQualifier := false
	if universalCreditRow != nil {
		benefitAmountStr := spreadsheet.ColByName(universalCreditRow, "aa")
		benefitAmount, err := strconv.Atoi(benefitAmountStr)
		if err == nil {
			ucQualifier = float32(benefitAmount) < inputData.benefitAmount
		}
	}

	qualifies := qualifierA || qualifierB || passportQualifier || ucQualifier

	qualifyType := ""
	if qualifierA {
		qualifyType = "CTC ONLY"
	} else if qualifierB {
		qualifyType = "CTC & WTC"
	} else if passportQualifier {
		qualifyType = "PASSPORTED"
	} else if ucQualifier {
		qualifyType = "UC QUALIFIER"
	}

	return qualifies, qualifyType
}

func sumFloatColumns(row spreadsheet.Row, colNames []string) float32 {
	var result float32 = 0
	for _, colName := range colNames {
		cellStr := spreadsheet.ColByName(row, colName)

		// Default to "0" for empty cells
		if cellStr == "" {
			cellStr = "0"
		}

		value, err := strconv.ParseFloat(cellStr, 32)
		if err != nil {
			llog.Printf(`Error parsing float from cell value "%s" for col name "%s", falling back to 0`, cellStr, colName)
			llog.Printf("\n")
			value = 0
		}
		result += float32(value)
	}
	return result
}
