package gateway

func (s *Service) list() []Transaction {
	transactions := make([]Transaction, 0)

	s.RLock()
	defer s.RUnlock()

	for _, trans := range s.transactions {
		transactions = append(transactions, Transaction{
			ID:       trans.id,
			State:    trans.state,
			Amount:   trans.amount,
			Captured: trans.captured,
			Currency: trans.currency,
		})
	}

	return transactions
}
