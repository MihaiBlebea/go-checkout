package gateway

func (s *Service) voidTransaction(id string) (int, string, error) {
	trans, err := s.getTransactionByID(id)
	if err != nil {
		return 0, "", err
	}

	trans.voided = true

	return trans.amount - trans.captured, trans.currency, nil
}
