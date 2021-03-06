package gateway

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"
)

var (
	InvalidTransactionErr  = errors.New("Invalid transaction id")
	InvalidAmountErr       = errors.New("Invalid payment amount id")
	InvalidNameOnCardErr   = errors.New("Invalid name on card")
	InvalidCvvErr          = errors.New("Invalid cvv. Must be 3 characters long")
	UnavailableAmountErr   = errors.New("Unavailable amount")
	InvalidCurrencyErr     = errors.New("Invalid currency")
	TransactionVoidedErr   = errors.New("Transaction is voided")
	TransactionRefundedErr = errors.New("Transaction is refunded")
	ExpiredCardErr         = errors.New("Card is expired")
	InvalidCardNumberErr   = errors.New("Card number is invalid")
	AuthFailedErr          = errors.New("Authorisation failed")
	CaptureFailedErr       = errors.New("Capture failed")
	RefundFailedErr        = errors.New("Refund failed")
)

func validateNameOnCard(name string) error {
	if strings.Contains(name, " ") == false {
		return InvalidNameOnCardErr
	}

	letters := strings.Split(strings.ReplaceAll(name, " ", ""), "")
	if len(letters) < 4 {
		return InvalidNameOnCardErr
	}

	return nil
}

func validateCVV(cvv string) error {
	re := regexp.MustCompile(`^[0-9]{3}$`)
	if re.MatchString(cvv) == false {
		return InvalidCvvErr
	}

	return nil
}

func validateAmount(amount int) error {
	if amount == 0 {
		return InvalidAmountErr
	}

	return nil
}

func validateTransactionID(transactions map[string]transaction, id string) error {
	_, ok := transactions[id]
	if ok == false {
		return InvalidTransactionErr
	}

	return nil
}

func validateCaptureAmount(trans *transaction) error {
	if trans.amount < trans.captured {
		return UnavailableAmountErr
	}

	return nil
}

func validateRefundAmount(trans *transaction) error {
	if trans.captured < trans.refunded {
		return UnavailableAmountErr
	}

	return nil
}

func validateCurrency(trans *transaction, currency string) error {
	if trans.currency != currency {
		return InvalidCurrencyErr
	}

	return nil
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
