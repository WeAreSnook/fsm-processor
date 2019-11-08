package main

import "fmt"

func main() {
	store := PeopleStore{}
	store.Add(Person{forename: "Chris"})
	store.Add(Person{forename: "Michael"})
	store.Add(Person{forename: "Charlotte"})

	fmt.Printf("%#v are the people in AJ", store.people)
}
