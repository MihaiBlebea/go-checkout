package sandbox_gateway_test

import (
	"regexp"
	"testing"

	"github.com/MihaiBlebea/go-checkout/gateway"
	sandbox "github.com/MihaiBlebea/go-checkout/sandbox_gateway"
)

func TestSanboxCardPayment(t *testing.T) {
	options := gateway.AuthorizeOptions{
		NameOnCard:  "Mihai Blebea",
		ExpireYear:  2099,
		ExpireMonth: 2,
		CVV:         "557",
		Amount:      200,
		Currency:    "GBP",
	}

	authCardCase := options
	authCardCase.CardNumber = sandbox.AuthFailCard

	captureCardCase := options
	captureCardCase.CardNumber = sandbox.CaptureFailCard

	refundCardCase := options
	refundCardCase.CardNumber = sandbox.RefundFailCard

	validCardCase := options
	validCardCase.CardNumber = "4111 1111 1111 1111"

	cases := []struct {
		title      string
		input      gateway.AuthorizeOptions
		authErr    error
		captureErr error
		refundErr  error
	}{
		{
			title:      "Auth failed card",
			input:      authCardCase,
			authErr:    gateway.AuthFailedErr,
			captureErr: nil,
			refundErr:  nil,
		},
		{
			title:      "Capture failed card",
			input:      captureCardCase,
			authErr:    nil,
			captureErr: gateway.CaptureFailedErr,
			refundErr:  nil,
		},
		{
			title:      "Refund failed card",
			input:      refundCardCase,
			authErr:    nil,
			captureErr: nil,
			refundErr:  gateway.RefundFailedErr,
		},
		{
			title:      "valid card case",
			input:      validCardCase,
			authErr:    nil,
			captureErr: nil,
			refundErr:  nil,
		},
	}

	g := gateway.New()

	sandboxGateway := sandbox.New(g)

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {

			id, authErr := sandboxGateway.AuthorizePayment(c.input)
			if c.authErr != nil {
				if authErr != c.authErr {
					t.Errorf("err: got %v want %v", authErr, c.authErr)
				}
			}

			_, _, captureErr := sandboxGateway.CaptureAmount(id, 100, "GBP")
			if c.captureErr != nil {
				if captureErr != c.captureErr {
					t.Errorf("err: got %v want %v", captureErr, c.captureErr)
				}
			}

			_, _, refundErr := sandboxGateway.RefundAmount(id, 50, "GBP")
			if c.refundErr != nil {
				if refundErr != c.refundErr {
					t.Errorf("err: got %v want %v", refundErr, c.refundErr)
				}
			}
		})
	}
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
