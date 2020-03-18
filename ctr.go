package main

import "github.com/addjam/fsm-processor/llog"

// GenerateCtrBasedAwards combined spreadsheet data and the output from the FSM algorithm
// to determine who can get clothing grant based on CTR. Outputs as an awards spreadsheet and
// a report for education
func GenerateCtrBasedAwards(inputData InputData, fsmStore PeopleStore) PeopleStore {
	handleErr := func(err error, store PeopleStore) {
		if err != nil {
			RespondWith(nil, &store, err)
			return
		}
	}

	store := PeopleStore{}
	var err error

	err = AddPeopleWithCtr(inputData, &store)
	handleErr(err, store)
	llog.Printf("%d people with CTR\n", len(store.People))

	store.People, err = PeopleInHouseholdsWithChildren(inputData, store)
	handleErr(err, store)
	llog.Printf("%d people after household check\n", len(store.People))

	// Mark everyone as CG eligible
	for i, p := range store.People {
		for j, d := range p.Dependents {
			d.NewCG = inputData.awardCG
			p.Dependents[j] = d
		}
		store.People[i] = p
	}

	nlcDependents, nonNlcDependents, err := PeopleWithChildrenAtNlcSchool(inputData, store)
	handleErr(err, store)
	store.ReportForEducationDependents = nonNlcDependents
	store.AwardDependents = nlcDependents
	llog.Printf("%d dependents in NLC schools, %d unmatched\n", len(nlcDependents), len(nonNlcDependents))

	store.AwardDependents = FillExistingGrants(inputData, store.AwardDependents)
	llog.Printf("got %d AwardDependents filled\n", len(store.AwardDependents))

	store.AwardDependents = filterNotReceivingCG(store.AwardDependents)
	llog.Printf("%d not receiving CG\n", len(store.AwardDependents))

	store.AwardDependents = FilterUsingExclusionList(inputData, store.AwardDependents)
	llog.Printf("%d after filtering exclusion list\n", len(store.AwardDependents))

	store.AwardDependents = filterDependents(store.AwardDependents, fsmStore.AwardDependents)
	llog.Printf("%d after filtering from awards list\n", len(store.AwardDependents))

	store.AwardDependents = FilterMinimumP1(store.AwardDependents)
	llog.Printf("%d in minimum P1\n", len(store.AwardDependents))

	GenerateAwardList(inputData, store, "ctr")
	GenerateEducationReport(inputData, store, "ctr")

	return store
}

func filterNotReceivingCG(dependents []Dependent) []Dependent {
	filtered := []Dependent{}

	for _, d := range dependents {
		if !d.ExistingCG {
			filtered = append(filtered, d)
		}
	}

	return filtered
}

func filterDependents(dependents []Dependent, filterList []Dependent) []Dependent {
	filtered := []Dependent{}

	for _, d := range dependents {
		isFiltered := false

		for _, filtered := range filterList {
			if d.Seemis == filtered.Seemis {
				isFiltered = true
				break
			}
		}

		if !isFiltered {
			filtered = append(filtered, d)
		}
	}

	return filtered
}
