package gateway

func (s *Service) refundAmount(id string, amount int, currency string) (int, string, error) {
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

	trans.refunded += amount

	s.transactions[id] = trans

	return trans.amount - trans.captured, trans.currency, nil
}
