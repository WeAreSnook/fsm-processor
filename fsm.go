package main

import "fmt"

// GenerateFsmAwards runs the FSM algorithm, combining input spreadsheets
// and outputting both an award list and a report for education.
func GenerateFsmAwards(inputData InputData) PeopleStore {
	handleErr := func(err error, store PeopleStore) {
		if err != nil {
			RespondWith(&store, err)
			return
		}
	}

	store := PeopleStore{}

	err := AddPeopleWithConsent(inputData, &store)
	handleErr(err, store)
	fmt.Printf("%d people with consent\n", len(store.People))

	store.People, err = PeopleInHouseholdsWithChildren(inputData, store)
	handleErr(err, store)
	fmt.Printf("%d people after household check\n", len(store.People))

	store.People, err = PeopleWithQualifyingIncomes(inputData, store)
	handleErr(err, store)
	fmt.Printf("%d people after income qualifying\n", len(store.People))

	nlcDependents, nonNlcDependents, err := PeopleWithChildrenAtNlcSchool(inputData, store)
	handleErr(err, store)
	store.ReportForEducationDependents = nonNlcDependents
	store.AwardDependents = nlcDependents
	fmt.Printf("%d dependents in NLC schools, %d unmatched\n", len(nlcDependents), len(nonNlcDependents))

	store.AwardDependents = FillExistingGrants(inputData, store.AwardDependents)
	fmt.Printf("got %d AwardDependents filled\n", len(store.AwardDependents))

	store.AwardDependents = FilterOnlyNewEntitlements(store.AwardDependents)
	fmt.Printf("%d have new entitlements\n", len(store.AwardDependents))

	store.AwardDependents = FilterMinimumP1(store.AwardDependents)
	fmt.Printf("%d are in at least P1\n", len(store.AwardDependents))

	store.AwardDependents = FilterUsingExclusionList(inputData, store.AwardDependents)
	fmt.Printf("Filtered to %d dependents\n", len(store.AwardDependents))

	// atLeast16, below16 := splitByMinimumAge(inputData, store.AwardDependents)
	// fmt.Printf("%d people at least age 16, %d below\n", len(atLeast16), len(below16))

	// TODO for atLeast16 => waiting for a flag to be added to school roll indicating if they are still in education
	// inEducation := below16

	GenerateAwardList(inputData, store, "fsm")
	GenerateEducationReport(inputData, store, "fsm")

	return store
}
