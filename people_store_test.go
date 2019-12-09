package main

import "testing"

func TestAdd(t *testing.T) {
	chris := Person{Forename: "Chris", Surname: "Sloey"}
	michael := Person{Forename: "Michael", Surname: "Hayes"}
	newPerson := Person{Forename: "Bob", Surname: "WOW"}

	store := PeopleStore{
		People: []Person{
			chris,
			michael,
		},
	}

	t.Run("adds person successfully", func(t *testing.T) {
		store.Add(newPerson)

		found := false
		for _, p := range store.People {
			if p.Forename == newPerson.Forename && p.Surname == newPerson.Surname {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("user wasn't added")
		}
	})
}

func TestUpdate(t *testing.T) {
	chris := Person{Forename: "Chris", Surname: "Sloey"}
	michael := Person{Forename: "Michael", Surname: "Hayes"}

	store := PeopleStore{
		People: []Person{
			chris,
			michael,
		},
	}

	t.Run("updates correct person", func(t *testing.T) {
		updatedChrisDetails := Person{Forename: "Christopher", Surname: "Sloey"}
		err := store.Update(updatedChrisDetails)

		if err != nil {
			t.Errorf("error updating user %#v", err)
		}

		found := false
		for _, p := range store.People {
			if p.Forename == updatedChrisDetails.Forename && p.Surname == updatedChrisDetails.Surname {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("user wasn't updated")
		}
	})
}
