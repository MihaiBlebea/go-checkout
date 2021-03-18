package gateway_test

import (
	"log"
	"regexp"
	"testing"

	"github.com/MihaiBlebea/go-checkout/gateway"
)

func TestAuthorizePayment(t *testing.T) {
	invalidNameOnCardCase := gateway.AuthorizeOptions{
		NameOnCard:  "Mih",
		CardNumber:  "4000 0000 0000 0259",
		ExpireYear:  2099,
		ExpireMonth: 2,
		CVV:         "557",
		Amount:      200,
		Currency:    "GBP",
	}

	invalidCardNumberCase := gateway.AuthorizeOptions{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "44539067130279178",
		ExpireYear:  2099,
		ExpireMonth: 2,
		CVV:         "557",
		Amount:      200,
		Currency:    "GBP",
	}

	invalidExpireDateCase := gateway.AuthorizeOptions{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0200",
		ExpireYear:  1989,
		ExpireMonth: 7,
		CVV:         "557",
		Amount:      200,
		Currency:    "GBP",
	}

	invalidCvvCase := gateway.AuthorizeOptions{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0200",
		ExpireYear:  2099,
		ExpireMonth: 7,
		CVV:         "5572",
		Amount:      200,
		Currency:    "GBP",
	}

	invalidAmountCase := gateway.AuthorizeOptions{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0200",
		ExpireYear:  2099,
		ExpireMonth: 7,
		CVV:         "557",
		Amount:      0,
		Currency:    "GBP",
	}

	successCase := gateway.AuthorizeOptions{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0259",
		ExpireYear:  2099,
		ExpireMonth: 7,
		CVV:         "557",
		Amount:      200,
		Currency:    "GBP",
	}

	cases := []struct {
		title string
		input gateway.AuthorizeOptions
		hasID bool
		err   error
	}{
		{
			title: "Invalid number on card case",
			input: invalidNameOnCardCase,
			hasID: false,
			err:   gateway.InvalidNameOnCardErr,
		},
		{
			title: "Invalid card number case",
			input: invalidCardNumberCase,
			hasID: false,
			err:   gateway.InvalidCardNumberErr,
		},
		{
			title: "Card expired case",
			input: invalidExpireDateCase,
			hasID: false,
			err:   gateway.ExpiredCardErr,
		},
		{
			title: "Invalid card cvv case",
			input: invalidCvvCase,
			hasID: false,
			err:   gateway.InvalidCvvErr,
		},
		{
			title: "Invalid amount case",
			input: invalidAmountCase,
			hasID: false,
			err:   gateway.InvalidAmountErr,
		},
		{
			title: "Success case",
			input: successCase,
			hasID: true,
			err:   nil,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			g := gateway.New()

			id, err := g.AuthorizePayment(c.input)

			if c.hasID {
				if isValidUUID(id) != true {
					t.Errorf("transaction id: got %v want %v", isValidUUID(id), true)
				}
			} else {
				if id != "" {
					t.Errorf("transaction id: got %v want %v", id, "")
				}
			}

			if err != nil && err.Error() != c.err.Error() {
				t.Errorf("err: got %v want %v", err.Error(), c.err.Error())
			}
		})
	}
}

func TestCaptureAmount(t *testing.T) {
	g := gateway.New()

	id, err := g.AuthorizePayment(gateway.AuthorizeOptions{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0259",
		ExpireYear:  2099,
		ExpireMonth: 7,
		CVV:         "557",
		Amount:      200,
		Currency:    "GBP",
	})
	if err != nil {
		log.Fatal(err)
	}

	remain, _, err := g.CaptureAmount(id, 10, "GBP")

	expected := 190
	if remain != expected {
		t.Errorf("remaining amount: got %v want %v", remain, expected)
	}

	if err != nil {
		t.Errorf("err: got %v want %v", err.Error(), nil)
	}
}

func TestCaptureAmountMultipleAttempts(t *testing.T) {
	g := gateway.New()

	id, err := g.AuthorizePayment(gateway.AuthorizeOptions{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0259",
		ExpireYear:  2099,
		ExpireMonth: 7,
		CVV:         "557",
		Amount:      200,
		Currency:    "GBP",
	})
	if err != nil {
		log.Fatal(err)
	}

	var amount = 200
	var attempts int
	for amount > 0 {
		attempts += 1
		remain, _, err := g.CaptureAmount(id, 10, "GBP")
		amount = remain

		expected := 200 - attempts*10
		if remain != expected {
			t.Errorf("remaining amount: got %v want %v", remain, expected)
		}

		if err != nil {
			t.Errorf("err: got %v want %v", err.Error(), nil)
		}
	}

	if attempts != 20 {
		t.Errorf("count nr of attempts to capture full amount: got %v want %v", attempts, 20)
	}
}

func TestFailCaptureIfVoidState(t *testing.T) {
	g := gateway.New()

	id, err := g.AuthorizePayment(gateway.AuthorizeOptions{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0259",
		ExpireYear:  2099,
		ExpireMonth: 7,
		CVV:         "557",
		Amount:      200,
		Currency:    "GBP",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = g.VoidTransaction(id)
	if err != nil {
		log.Fatal(err)
	}

	remain, _, err := g.CaptureAmount(id, 10, "GBP")
	expected := 0
	if remain != expected {
		t.Errorf("remaining amount: got %v want %v", remain, expected)
	}

	if err != nil && err.Error() != gateway.TransactionVoidedErr.Error() {
		t.Errorf("err: got %v want %v", err.Error(), gateway.TransactionVoidedErr.Error())
	}
}

func TestFailCaptureIfRefundState(t *testing.T) {
	g := gateway.New()

	id, err := g.AuthorizePayment(gateway.AuthorizeOptions{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0259",
		ExpireYear:  2099,
		ExpireMonth: 7,
		CVV:         "557",
		Amount:      300,
		Currency:    "GBP",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Capture £2
	_, _, err = g.CaptureAmount(id, 200, "GBP")
	if err != nil {
		log.Fatal(err)
	}

	// Refund £1
	_, _, err = g.RefundAmount(id, 100, "GBP")
	if err != nil {
		log.Fatal(err)
	}

	// Try to refund another £1
	remain, _, err := g.CaptureAmount(id, 100, "GBP")
	expected := 0
	if remain != expected {
		t.Errorf("remaining amount: got %v want %v", remain, expected)
	}

	if err != nil && err.Error() != gateway.TransactionRefundedErr.Error() {
		t.Errorf("err: got %v want %v", err.Error(), gateway.TransactionRefundedErr.Error())
	}
}

func TestFailCaptureMoreThanAmmount(t *testing.T) {
	g := gateway.New()

	id, err := g.AuthorizePayment(gateway.AuthorizeOptions{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0259",
		ExpireYear:  2099,
		ExpireMonth: 7,
		CVV:         "557",
		Amount:      100,
		Currency:    "GBP",
	})
	if err != nil {
		log.Fatal(err)
	}

	remain, _, err := g.CaptureAmount(id, 200, "GBP")
	expected := 0
	if remain != expected {
		t.Errorf("remaining amount: got %v want %v", remain, expected)
	}

	if err != nil && err.Error() != gateway.UnavailableAmountErr.Error() {
		t.Errorf("err: got %v want %v", err.Error(), gateway.TransactionRefundedErr.Error())
	}
}

func TestRefundFullAmount(t *testing.T) {
	g := gateway.New()

	id, err := g.AuthorizePayment(gateway.AuthorizeOptions{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0259",
		ExpireYear:  2099,
		ExpireMonth: 7,
		CVV:         "557",
		Amount:      100,
		Currency:    "GBP",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = g.CaptureAmount(id, 100, "GBP")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 2; i++ {

		// First attempt for 50
		remain, _, err := g.RefundAmount(id, 50, "GBP")
		expcted := 100 - 50*(i+1)
		if remain != expcted {
			t.Errorf("remaining amount: got %v want %v", remain, expcted)
		}

		if err != nil {
			t.Errorf("err: got %v want %v", err.Error(), nil)
		}
	}

	// Third attempt for another 50 should fail
	remain, _, err := g.RefundAmount(id, 50, "GBP")
	if remain != 0 {
		t.Errorf("remaining amount: got %v want %v", remain, 0)
	}

	if err != nil && err.Error() != gateway.UnavailableAmountErr.Error() {
		t.Errorf("err: got %v want %v", err.Error(), gateway.UnavailableAmountErr.Error())
	}
}

func TestFailRefundIfNotCaptured(t *testing.T) {
	g := gateway.New()

	id, err := g.AuthorizePayment(gateway.AuthorizeOptions{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0259",
		ExpireYear:  2099,
		ExpireMonth: 7,
		CVV:         "557",
		Amount:      100,
		Currency:    "GBP",
	})
	if err != nil {
		log.Fatal(err)
	}

	remain, _, err := g.RefundAmount(id, 50, "GBP")
	if remain != 0 {
		t.Errorf("remaining amount: got %v want %v", remain, 0)
	}

	if err != nil && err.Error() != gateway.UnavailableAmountErr.Error() {
		t.Errorf("err: got %v want %v", err.Error(), gateway.UnavailableAmountErr.Error())
	}
}

func TestFailRefundIfVoidState(t *testing.T) {
	g := gateway.New()

	id, err := g.AuthorizePayment(gateway.AuthorizeOptions{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0259",
		ExpireYear:  2099,
		ExpireMonth: 7,
		CVV:         "557",
		Amount:      100,
		Currency:    "GBP",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = g.VoidTransaction(id)
	if err != nil {
		log.Fatal(err)
	}

	remain, _, err := g.RefundAmount(id, 50, "GBP")
	if remain != 0 {
		t.Errorf("remaining amount: got %v want %v", remain, 0)
	}

	if err != nil && err.Error() != gateway.TransactionVoidedErr.Error() {
		t.Errorf("err: got %v want %v", err.Error(), gateway.TransactionVoidedErr.Error())
	}
}

func TestVoidTransactionIfRefundState(t *testing.T) {
	g := gateway.New()

	id, err := g.AuthorizePayment(gateway.AuthorizeOptions{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0259",
		ExpireYear:  2099,
		ExpireMonth: 7,
		CVV:         "557",
		Amount:      100,
		Currency:    "GBP",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Capture amount
	_, _, err = g.CaptureAmount(id, 100, "GBP")
	if err != nil {
		log.Fatal(err)
	}

	// Refund amount
	_, _, err = g.RefundAmount(id, 50, "GBP")
	if err != nil {
		log.Fatal(err)
	}

	// Try to void the transaction
	remain, _, err := g.VoidTransaction(id)
	if err != nil {
		log.Fatal(err)
	}
	if remain != 50 {
		t.Errorf("remaining amount: got %v want %v", remain, 50)
	}

	if err != nil {
		t.Errorf("err: got %v want %v", err.Error(), nil)
	}
}

func TestListTransactions(t *testing.T) {
	g := gateway.New()

	var transIDs []string
	for i := 0; i < 10; i++ {
		id, err := g.AuthorizePayment(gateway.AuthorizeOptions{
			NameOnCard:  "Mihai Blebea",
			CardNumber:  "4000 0000 0000 0259",
			ExpireYear:  2099,
			ExpireMonth: 7,
			CVV:         "557",
			Amount:      100,
			Currency:    "GBP",
		})
		if err != nil {
			log.Fatal(err)
		}

		transIDs = append(transIDs, id)
	}

	// Capture amount
	transactions := g.ListTransactions()

	expected := 10
	if len(transactions) != expected {
		t.Errorf("transactions count: got %v want %v", len(transactions), expected)
	}

	for _, trans := range transactions {
		if _, ok := find(transIDs, trans.ID); ok != true {
			t.Errorf("expected id in slice for id %s: got %v want %v", trans.ID, ok, true)
		}
	}
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}

	return -1, false
}
