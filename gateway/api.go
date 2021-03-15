package gateway

type Service struct {
	transactions []transaction
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

func New() *Service {
	return &Service{}
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
