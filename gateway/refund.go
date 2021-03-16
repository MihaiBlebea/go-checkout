package gateway

func (s *Service) refundAmount(id string, amount int, currency string) (int, string, error) {
	if err := validateTransactionID(s.transactions, id); err != nil {
		return 0, "", err
	}

	trans := s.transactions[id]

	if trans.state != VoidState {
		return 0, "", TransactionVoidedErr
	}

	trans.state = RefundState
	trans.refunded += amount

	if err := validateRefundAmount(&trans); err != nil {
		return 0, "", err
	}

	s.transactions[id] = trans

	return calcRemainRefundAmount(&trans), trans.currency, nil
}

func calcRemainRefundAmount(trans *transaction) int {
	return trans.captured - trans.refunded
}
