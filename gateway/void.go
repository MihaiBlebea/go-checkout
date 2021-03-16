package gateway

func (s *Service) voidTransaction(id string) (int, string, error) {
	if err := s.validateTransactionID(id); err != nil {
		return 0, "", err
	}

	trans := s.transactions[id]

	trans.voided = true

	s.transactions[id] = trans

	return trans.amount - trans.captured, trans.currency, nil
}
