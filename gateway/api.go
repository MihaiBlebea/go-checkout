package gateway

type Service struct {
	transactions map[string]transaction
}

type AuthorizeOptions struct {
	NameOnCard  string `json:"name_on_card"`
	CardNumber  string `json:"card_number"`
	ExpireYear  int    `json:"expire_year"`
	ExpireMonth int    `json:"expire_month"`
	CVV         string `json:"cvv"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
}

type Transaction struct {
	ID       string
	Amount   int
	Currency string
	Voided   bool
	Refunded int
	Captured int
}

func New() *Service {
	return &Service{make(map[string]transaction)}
}

func (s *Service) AuthorizePayment(options AuthorizeOptions) (string, error) {
	return s.authorizePayment(options)
}

func (s *Service) CaptureAmount(id string, amount int, currency string) (int, string, error) {
	return s.captureAmount(id, amount, currency)
}

func (s *Service) VoidTransaction(id string) (int, string, error) {
	return s.voidTransaction(id)
}

func (s *Service) RefundAmount(id string, amount int, currency string) (int, string, error) {
	return s.refundAmount(id, amount, currency)
}

func (s *Service) ListTransactions() (transactions []Transaction) {
	for _, trans := range s.transactions {
		transactions = append(transactions, Transaction{
			ID:       trans.id,
			Amount:   trans.amount,
			Currency: trans.currency,
			Voided:   trans.voided,
			Refunded: trans.refunded,
			Captured: trans.captured,
		})
	}

	return transactions
}
