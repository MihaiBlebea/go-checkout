package gateway

import (
	"errors"
)

var (
	InvalidTransactionErr error = errors.New("Invalid transaction id")
	UnavailableAmountErr  error = errors.New("Unavailable amount to capture")
	InvalidCurrencyErr    error = errors.New("Invalid currency")
	TransactionVoidedErr  error = errors.New("Transaction is voided")
)

func (s *Service) captureAmount(id string, amount int, currency string) (int, string, error) {
	if err := s.validateTransactionID(id); err != nil {
		return 0, "", err
	}

	trans := s.transactions[id]

	if err := validateAvailableAmount(trans, amount); err != nil {
		return 0, "", err
	}

	if trans.voided == true {
		return 0, "", TransactionVoidedErr
	}

	trans.captured += amount

	s.transactions[id] = trans

	return trans.amount - trans.captured, trans.currency, nil
}

func (s *Service) validateTransactionID(id string) error {
	_, ok := s.transactions[id]
	if ok == false {
		return InvalidTransactionErr
	}

	return nil
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
