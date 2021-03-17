package validate_test

import (
	"errors"
	"testing"

	"github.com/MihaiBlebea/go-checkout/server/validate"
)

func TestValidateRequiredField(t *testing.T) {
	type ShoppingList struct {
		Bread  int `validate:"required"`
		Beer   int `validate:"required"`
		Coffee int `json:"coffee"`
	}
	cases := []struct {
		title string
		input ShoppingList
		want  error
	}{
		{
			title: "Valid required keys present",
			input: ShoppingList{1, 2, 0},
			want:  nil,
		},
		{
			title: "Invalid required keys missing Bread",
			input: ShoppingList{0, 0, 1},
			want:  errors.New("Field Bread is required"),
		},
		{
			title: "Invalid required keys missing Beer",
			input: ShoppingList{1, 0, 0},
			want:  errors.New("Field Beer is required"),
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			err := validate.Validate(&c.input)

			if err != nil && err.Error() != c.want.Error() {
				t.Errorf("err: got %v want %v", err.Error(), c.want.Error())
			}
		})
	}
}
