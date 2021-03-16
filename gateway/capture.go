package gateway

func (s *Service) captureAmount(id string, amount int, currency string) (int, string, error) {
	if err := validateTransactionID(s.transactions, id); err != nil {
		return 0, "", err
	}

	trans := s.transactions[id]

	if trans.state != CaptureState {
		return 0, "", TransactionVoidedErr
	}

	trans.captured += amount

	if err := validateCaptureAmount(&trans); err != nil {
		return 0, "", err
	}

	s.transactions[id] = trans

	return calcRemainCaptureAmount(&trans), trans.currency, nil
}

func calcRemainCaptureAmount(trans *transaction) int {
	return trans.amount - trans.captured
}
