package main

import "testing"

func TestFindExisting(t *testing.T) {
	chris := Person{forename: "Chris", surname: "Sloey", age: 29}
	michael := Person{forename: "Michael", surname: "Hayes", age: 31}

	store := PersonStore{
		people: []Person{
			chris,
			michael,
		},
	}

	t.Run("finds the correct person", func(t *testing.T) {
		want := chris
		got, err := store.FindExisting(Person{forename: chris.forename, surname: chris.surname})

		if err != nil {
			t.Errorf("error searching for %#v", want)
		}

		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("error when person doesn't exist", func(t *testing.T) {
		unknownPerson := Person{forename: "Doesn't", surname: "Exist"}
		_, err := store.FindExisting(unknownPerson)

		if err == nil {
			t.Errorf("No error for missing person")
		}
	})
}
