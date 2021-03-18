package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	gtway "github.com/MihaiBlebea/go-checkout/gateway"
	sandbox "github.com/MihaiBlebea/go-checkout/sandbox_gateway"
	"github.com/MihaiBlebea/go-checkout/server/handler"
	hand "github.com/MihaiBlebea/go-checkout/server/handler"
	"github.com/MihaiBlebea/go-checkout/server/validate"
)

var logger *MuteLogger

type MuteLogger struct{}

func (l *MuteLogger) Info(args ...interface{})    {}
func (l *MuteLogger) Trace(args ...interface{})   {}
func (l *MuteLogger) Debug(args ...interface{})   {}
func (l *MuteLogger) Print(args ...interface{})   {}
func (l *MuteLogger) Warn(args ...interface{})    {}
func (l *MuteLogger) Warning(args ...interface{}) {}
func (l *MuteLogger) Error(args ...interface{})   {}
func (l *MuteLogger) Fatal(args ...interface{})   {}
func (l *MuteLogger) Panic(args ...interface{})   {}

func TestHealthCheckEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.HealthEndpoint(logger).ServeHTTP)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response struct {
		OK bool `json:"ok"`
	}

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.OK != true {
		t.Errorf("handler returned unexpected body: got %v want %v", response.OK, true)
	}
}

func TestAuthorizeEndpoint(t *testing.T) {
	requestSuccess := hand.AuthorizeRequest{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0259",
		ExpireYear:  2022,
		ExpireMonth: 4,
		CVV:         "755",
		Amount:      200,
		Currency:    "GBP",
	}

	responseSuccess := hand.AuthorizeResponse{
		Success:  true,
		Message:  "",
		Amount:   200,
		Currency: "GBP",
	}

	requestFailingCard := hand.AuthorizeRequest{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0119",
		ExpireYear:  2022,
		ExpireMonth: 4,
		CVV:         "755",
		Amount:      200,
		Currency:    "GBP",
	}

	responseFailingCard := hand.AuthorizeResponse{
		Success:  false,
		Message:  "Authorisation failed",
		Amount:   0,
		Currency: "",
	}

	requestMissingKey := hand.AuthorizeRequest{
		NameOnCard: "Mihai Blebea",
		CardNumber: "4000 0000 0000 0119",
		CVV:        "755",
		Amount:     200,
		Currency:   "GBP",
	}

	responseMissingKey := hand.AuthorizeResponse{
		Success:  false,
		Message:  "Field ExpireYear is required",
		Amount:   0,
		Currency: "",
	}

	requestMissingStringKey := hand.AuthorizeRequest{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0119",
		ExpireYear:  2022,
		ExpireMonth: 4,
		Amount:      200,
		Currency:    "GBP",
	}

	responseMissingStringKey := hand.AuthorizeResponse{
		Success:  false,
		Message:  "Field CVV is required",
		Amount:   0,
		Currency: "",
	}

	cases := []struct {
		input hand.AuthorizeRequest
		want  hand.AuthorizeResponse
		code  int
	}{
		{
			input: requestSuccess,
			want:  responseSuccess,
			code:  200,
		},
		{
			input: requestFailingCard,
			want:  responseFailingCard,
			code:  400,
		},
		{
			input: requestMissingKey,
			want:  responseMissingKey,
			code:  400,
		},
		{
			input: requestMissingStringKey,
			want:  responseMissingStringKey,
			code:  400,
		},
	}

	for _, c := range cases {
		b, err := json.Marshal(c.input)
		if err != nil {
			t.Fatal()
		}

		req, err := http.NewRequest("POST", "/authorize", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		// failed card aware gateway
		g := sandbox.New(gtway.New())

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(
			hand.AuthorizeEndpoint(g, logger, validate.Validate).ServeHTTP,
		)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != c.code {
			t.Errorf("http status code: got %d want %d", rr.Code, http.StatusOK)
		}

		var response hand.AuthorizeResponse
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		if response.Success != c.want.Success {
			t.Errorf("response success key: got %v want %v", response.Success, c.want.Success)
		}

		if response.Message != c.want.Message {
			t.Errorf("response message key: got %v want %v", response.Message, c.want.Message)
		}

		if response.ID != "" && isValidUUID(response.ID) != true {
			t.Errorf("response id is valid uuid format: got %v want %v", isValidUUID(response.ID), true)
		}

		if response.Amount != c.want.Amount {
			t.Errorf("response amount key: got %v want %v", response.Amount, c.want.Amount)
		}

		if response.Currency != c.want.Currency {
			t.Errorf("response amount key: got %v want %v", response.Currency, c.want.Currency)
		}
	}
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
