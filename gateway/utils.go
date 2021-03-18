package gateway

func (s *Service) getCardNumber(id string) (string, error) {
	s.RLock()
	defer s.RUnlock()

	if err := validateTransactionID(s.transactions, id); err != nil {
		return "", err
	}

	return s.transactions[id].cardNumber, nil
}
