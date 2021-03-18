package sandbox_gateway

import (
	gtway "github.com/MihaiBlebea/go-checkout/gateway"
)

const (
	AuthFailCard    string = "4000 0000 0000 0119"
	CaptureFailCard string = "4000 0000 0000 0259"
	RefundFailCard  string = "4000 0000 0000 3238"
)

type Service struct {
	gateway *gtway.Service
}

func New(gateway *gtway.Service) *Service {
	return &Service{
		gateway,
	}
}

func (s *Service) AuthorizePayment(options gtway.AuthorizeOptions) (string, error) {
	if options.CardNumber == AuthFailCard {
		return "", gtway.AuthFailedErr
	}

	return s.gateway.AuthorizePayment(options)
}

func (s *Service) CaptureAmount(id string, amount int, currency string) (int, string, error) {
	cardNumber, err := s.gateway.GetCardNumber(id)
	if err != nil {
		return 0, "", err
	}

	if cardNumber == CaptureFailCard {
		return 0, "", gtway.CaptureFailedErr
	}

	return s.gateway.CaptureAmount(id, amount, currency)
}

func (s *Service) VoidTransaction(id string) (int, string, error) {
	return s.gateway.VoidTransaction(id)
}

func (s *Service) RefundAmount(id string, amount int, currency string) (int, string, error) {
	cardNumber, err := s.gateway.GetCardNumber(id)
	if err != nil {
		return 0, "", err
	}

	if cardNumber == RefundFailCard {
		return 0, "", gtway.RefundFailedErr
	}

	return s.gateway.RefundAmount(id, amount, currency)
}

func (s *Service) ListTransactions() []gtway.Transaction {
	return s.gateway.ListTransactions()
}
