package gateway

import (
	"testing"
)

func TestValidateNameOnCard(t *testing.T) {
	cases := []struct {
		title string
		input string
		want  error
	}{
		{
			title: "Invalid 3 letter name",
			input: "Mih",
			want:  InvalidNameOnCardErr,
		},
		{
			title: "Invalid just first name",
			input: "Mihai",
			want:  InvalidNameOnCardErr,
		},
		{
			title: "Valid full name",
			input: "Mihai Blebea",
			want:  nil,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			err := validateNameOnCard(c.input)

			if err != c.want {
				t.Errorf("err: got %v want %v", err, c.want)
			}
		})
	}
}

func TestValidateCvv(t *testing.T) {
	cases := []struct {
		title string
		input string
		want  error
	}{
		{
			title: "Invalid 1 number cvv",
			input: "2",
			want:  InvalidCvvErr,
		},
		{
			title: "Invalid 2 number cvv",
			input: "27",
			want:  InvalidCvvErr,
		},
		{
			title: "Invalid 4 number cvv",
			input: "1234",
			want:  InvalidCvvErr,
		},
		{
			title: "Invalid letters in cvv",
			input: "12F",
			want:  InvalidCvvErr,
		},
		{
			title: "Valid cvv ",
			input: "557",
			want:  nil,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			err := validateCVV(c.input)

			if err != c.want {
				t.Errorf("err: got %v want %v", err, c.want)
			}
		})
	}
}

func TestValidateAmount(t *testing.T) {
	cases := []struct {
		title string
		input int
		want  error
	}{
		{
			title: "Invalid 0 amount",
			input: 0,
			want:  InvalidAmountErr,
		},
		{
			title: "Valid 100 amount",
			input: 100,
			want:  nil,
		},
		{
			title: "Valid 1 amount",
			input: 1,
			want:  nil,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			err := validateAmount(c.input)

			if err != c.want {
				t.Errorf("err: got %v want %v", err, c.want)
			}
		})
	}
}

func TestValidateTransactionID(t *testing.T) {
	transactions := make(map[string]transaction)
	transactions["1234"] = transaction{}
	transactions["abcd"] = transaction{}
	transactions["1234abcd"] = transaction{}

	cases := []struct {
		title string
		input string
		want  error
	}{
		{
			title: "Invalid abcde id",
			input: "abcde",
			want:  InvalidTransactionErr,
		},
		{
			title: "Valid 1234 id",
			input: "1234",
			want:  nil,
		},
		{
			title: "Valid abcd id",
			input: "abcd",
			want:  nil,
		},
		{
			title: "Valid 1234abcd id",
			input: "1234abcd",
			want:  nil,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			err := validateTransactionID(transactions, c.input)

			if err != c.want {
				t.Errorf("err: got %v want %v", err, c.want)
			}
		})
	}
}

func TestValidateCaptureAmount(t *testing.T) {
	cases := []struct {
		title string
		input transaction
		want  error
	}{
		{
			title: "Valid captured 20 of 200",
			input: transaction{amount: 200, captured: 20},
			want:  nil,
		},
		{
			title: "Invalid captured 200 of 20",
			input: transaction{amount: 20, captured: 200},
			want:  UnavailableAmountErr,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			err := validateCaptureAmount(&c.input)

			if err != c.want {
				t.Errorf("err: got %v want %v", err, c.want)
			}
		})
	}
}

func TestValidateRefundAmount(t *testing.T) {
	cases := []struct {
		title string
		input transaction
		want  error
	}{
		{
			title: "Valid captured 200 refunded 20",
			input: transaction{captured: 200, refunded: 20},
			want:  nil,
		},
		{
			title: "Invalid captured 20 refunded 200",
			input: transaction{captured: 20, refunded: 200},
			want:  UnavailableAmountErr,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			err := validateRefundAmount(&c.input)

			if err != c.want {
				t.Errorf("err: got %v want %v", err, c.want)
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	cases := []struct {
		title string
		month int
		year  int
		want  error
	}{
		{
			title: "Invalid expired card 2/1999",
			month: 2,
			year:  1999,
			want:  ExpiredCardErr,
		},
		{
			title: "Invalid expired card 7/1989",
			month: 7,
			year:  1989,
			want:  ExpiredCardErr,
		},
		{
			title: "Valid expired card 1/2099",
			month: 1,
			year:  2099,
			want:  nil,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			err := validateDate(c.month, c.year)

			if err != c.want {
				t.Errorf("err: got %v want %v", err, c.want)
			}
		})
	}
}

func TestValidateCardNumber(t *testing.T) {
	cases := []struct {
		title string
		card  string
		want  error
	}{
		{
			title: "Valid card number 4001 5900 0000 0001",
			card:  "4001 5900 0000 0001",
			want:  nil,
		},
		{
			title: "Valid card number 4111 1111 1111 1111",
			card:  "4111 1111 1111 1111",
			want:  nil,
		},
		{
			title: "Valid card number 4111111111111111 no spaces",
			card:  "4111111111111111",
			want:  nil,
		},
		{
			title: "Invalid card number 4111 1111 1111 1112",
			card:  "4111 1111 1111 1112",
			want:  InvalidCardNumberErr,
		},
		{
			title: "Invalid card number 4001 5900 0000 0005",
			card:  "4001 5900 0000 0005",
			want:  InvalidCardNumberErr,
		},
		{
			title: "Invalid card number with letters 4001 5900 0000 00abcd",
			card:  "4001 5900 0000 00abcd",
			want:  InvalidCardNumberErr,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			err := validateCardNumber(c.card)

			if err != c.want {
				t.Errorf("err: got %v want %v", err, c.want)
			}
		})
	}
}

func TestReverseString(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{
			input: "abcd",
			want:  "dcba",
		},
		{
			input: "1234",
			want:  "4321",
		},
		{
			input: "12 34",
			want:  "43 21",
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			output := reverse(c.input)

			if output != c.want {
				t.Errorf("err: got %v want %v", output, c.want)
			}
		})
	}
}
