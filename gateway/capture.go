package gateway

import "errors"

var (
	InvalidTransactionErr error = errors.New("Invalid transaction id")
	UnavailableAmountErr  error = errors.New("Unavailable amount to capture")
	InvalidCurrencyErr    error = errors.New("Invalid currency")
)

func (s *Service) captureAmount(id string, amount int, currency string) (int, string, error) {
	trans, err := s.getTransactionByID(id)
	if err != nil {
		return 0, "", err
	}

	if err := validateAvailableAmount(*trans, amount); err != nil {
		return 0, "", err
	}

	trans.captured += amount

	return trans.amount - amount, trans.currency, nil
}

func (s *Service) getTransactionByID(id string) (*transaction, error) {
	for _, trans := range s.transactions {
		if trans.id == id {
			return &trans, nil
		}
	}

	return &transaction{}, nil
}

func validateAvailableAmount(trans transaction, amount int) error {
	if trans.amount < amount {
		return UnavailableAmountErr
	}

	return nil
}

func validateCurrency(trans transaction, currency string) error {
	if trans.currency != currency {
		return InvalidCurrencyErr
	}

	return nil
}
