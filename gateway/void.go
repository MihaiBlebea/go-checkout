package gateway

func (s *Service) voidTransaction(id string) (int, string, error) {
	if err := validateTransactionID(s.transactions, id); err != nil {
		return 0, "", err
	}

	trans := s.transactions[id]

	trans.state = VoidState

	s.transactions[id] = trans

	return calcRemainRefundAmount(&trans), trans.currency, nil
}
