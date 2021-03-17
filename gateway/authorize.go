package gateway

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

type transaction struct {
	id          string
	state       TransactioState
	nameOnCard  string
	cardNumber  string
	expireYear  int
	expireMonth int
	cvv         string
	amount      int
	captured    int
	refunded    int
	currency    string
}

func (s *Service) authorizePayment(options AuthorizeOptions) (string, error) {
	if err := validateDate(options.ExpireMonth, options.ExpireYear); err != nil {
		return "", err
	}

	if err := validateCardNumber(options.CardNumber); err != nil {
		return "", err
	}

	if err := validateNameOnCard(options.NameOnCard); err != nil {
		return "", err
	}

	if err := validateCVV(options.CVV); err != nil {
		return "", err
	}

	if err := validateAmount(options.Amount); err != nil {
		return "", err
	}

	if isSandboxCard(options.CardNumber) {
		return "", AuthFailedErr
	}

	trans := transaction{
		id:          genToken(),
		nameOnCard:  options.NameOnCard,
		cardNumber:  options.CardNumber,
		expireMonth: options.ExpireMonth,
		expireYear:  options.ExpireYear,
		cvv:         options.CVV,
		amount:      options.Amount,
		currency:    options.Currency,
	}
	s.storeCard(trans)

	return trans.id, nil
}

// TODO: move this to it's own module or file (consolidate)
func isSandboxCard(cardNumber string) bool {
	cardNum := strings.ReplaceAll(cardNumber, " ", "")
	sandoxCardNum := strings.ReplaceAll(AuthFailCard, " ", "")

	return cardNum == sandoxCardNum
}

func genToken() string {
	return uuid.NewV4().String()
}

func (s *Service) storeCard(transaction transaction) {
	s.transactions[transaction.id] = transaction
}
