# FSM Processor

## Functions

### func [AddPeopleWithConsent](/consent.go#L15)

`func AddPeopleWithConsent(inputData InputData, peopleStore *PeopleStore) error`

AddPeopleWithConsent parses which people have given consent to check entitlement data
and adds them directly to the PeopleStore
Data sources: Consent 360 & Benefit Extract

### func [AddPeopleWithCtr](/income_check.go#L116)

`func AddPeopleWithCtr(inputData InputData, store *PeopleStore) error`

AddPeopleWithCtr adds people to the store who are receiging a
weekly cts entitlement greater than 0

### func [CleanString](/helpers.go#L16)

`func CleanString(str string) string`

CleanString replaces puncutation and spaces, and lowercases the string

### func [CompareCleanedStrings](/helpers.go#L26)

`func CompareCleanedStrings(a, b string) float64`

CompareCleanedStrings cleans inputs and passes to CompareStrings

### func [CompareStrings](/helpers.go#L21)

`func CompareStrings(a, b string) float64`

CompareStrings returns the jaro winkler distance from 0 (no similarity) to 1 (identical) between two strings

### func [GenerateAwardList](/award_list.go#L14)

`func GenerateAwardList(inputData InputData, store PeopleStore, name string)`

GenerateAwardList looks at the AwardDependents and generates an award list spreadsheet

### func [GenerateEducationReport](/education_report.go#L13)

`func GenerateEducationReport(inputData InputData, store PeopleStore, name string)`

GenerateEducationReport generates a spreadsheet of people who were not found in the school roll

### func [PeopleWithChildrenAtNlcSchool](/school_matcher.go#L22)

`func PeopleWithChildrenAtNlcSchool(inputData InputData, store PeopleStore) (matched []Dependent, unmatched []Dependent, err error)`

PeopleWithChildrenAtNlcSchool returns just the people from the store
that are likely matches for people in the school roll

### func [RespondWith](/response.go#L17)

`func RespondWith(fsmStore *PeopleStore, ctrStore *PeopleStore, err error)`

RespondWith stops execution and outputs response data as json

fsmStore - PeopleStore representing the final state of the FSM algorithm data
ctrStore - PeopleStore representing the final state of the FSM algorithm data
err - optional error that halted execution

### func [SplitByMinimumAge](/helpers.go#L58)

`func SplitByMinimumAge(inputData InputData, dependents []Dependent) (atThreshold []Dependent, belowThreshold []Dependent)`

SplitByMinimumAge splits the dependents into an array >= 16 and an array < 16 years old

## Sub Packages

* [llog](./llog)

* [spreadsheet](./spreadsheet)
