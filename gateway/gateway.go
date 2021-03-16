package gateway

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"

	uuid "github.com/satori/go.uuid"
)

const (
	AuthFailCard    string = "4000 0000 0000 0119"
	CaptureFailCard string = "4000 0000 0000 0259"
	RefundFailCard  string = "4000 0000 0000 3238"
)

var (
	ExpiredCardErr       error = errors.New("Card is expired")
	InvalidCardNumberErr error = errors.New("Card number is invalid")
	AuthFailedErr        error = errors.New("Authorisation failed")
)

type transaction struct {
	id          string
	nameOnCard  string
	cardNumber  string
	expireYear  int
	expireMonth int
	cvv         string
	amount      int
	currency    string
	refunded    int
	captured    int
	voided      bool
}

func (s *Service) authorizePayment(options AuthorizeOptions) (string, error) {
	if err := validateDate(options.ExpireMonth, options.ExpireYear); err != nil {
		return "", err
	}

	if err := validateCardNumber(options.CardNumber); err != nil {
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

func validateDate(month int, year int) error {
	now := time.Now()
	expire := time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC)

	if now.After(expire) {
		return ExpiredCardErr
	}

	return nil
}

func validateCardNumber(cardNumber string) error {
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")
	re := regexp.MustCompile(`^(?:4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11})$`)
	if re.MatchString(cardNumber) == false {
		return InvalidCardNumberErr
	}

	if validateLuhnAlgoritm(cardNumber) == false {
		return InvalidCardNumberErr
	}

	return nil
}

func validateLuhnAlgoritm(cardNumber string) bool {
	sum := 0
	counter := 0

	for _, r := range reverse(cardNumber) {
		if unicode.IsDigit(r) {
			val := int(r - '0')

			if counter%2 == 1 {
				val = (val * 2)
				if val > 9 {
					val = val - 9
				}
			}
			sum += val
			counter++
			continue
		}

		if unicode.IsSpace(r) {
			continue
		}

		return false
	}

	if counter < 2 {
		return false
	}

	return (sum % 10) == 0
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

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
