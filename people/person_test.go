package people

import "testing"

func TestIsSameAs(t *testing.T) {
	t.Run("returns true for users with the same name", func(t *testing.T) {
		userA := Person{Forename: "bob", Surname: "smith"}
		userB := Person{Forename: "bob", Surname: "smith"}

		if !userA.IsSameAs(userB) {
			t.Errorf("Expected %#v to be same as %#v but got false", userA, userB)
		}
	})

	t.Run("returns true for users with the same name", func(t *testing.T) {
		userA := Person{Forename: "bob", Surname: "smith"}
		userB := Person{Forename: "jane", Surname: "smith"}

		if userA.IsSameAs(userB) {
			t.Errorf("Expected %#v NOT to be same as %#v but got false", userA, userB)
		}
	})
}
