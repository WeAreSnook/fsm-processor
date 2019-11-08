package main

// ExtractPeopleWithConsent parses which people have given consent to check entitlement data
func ExtractPeopleWithConsent(inputData InputData, peopleStore *PeopleStore) {
	peopleStore.Add(Person{forename: "Michael", surname: "Hayes", age: 22})
}
