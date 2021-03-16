package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	gtway "github.com/MihaiBlebea/go-checkout/gateway"
	"github.com/MihaiBlebea/go-checkout/server/handler"
	hand "github.com/MihaiBlebea/go-checkout/server/handler"
	"github.com/MihaiBlebea/go-checkout/server/resp"
	"github.com/MihaiBlebea/go-checkout/server/validate"
)

func TestHealthCheckEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.HealthEndpoint().ServeHTTP)

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
	body := hand.AuthorizeRequest{
		NameOnCard:  "Mihai Blebea",
		CardNumber:  "4000 0000 0000 0259",
		ExpireYear:  2022,
		ExpireMonth: 4,
		CVV:         "755",
		Amount:      200,
		Currency:    "GBP",
	}

	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal()
	}

	req, err := http.NewRequest("POST", "/authorize", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(
		hand.AuthorizeEndpoint(gtway.New(), validate.Validate, resp.ErrorResponse).ServeHTTP,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response hand.AuthorizeResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success != true {
		t.Errorf("response success key: got %v want %v", response.Success, true)
	}

	if response.Message != "" {
		t.Errorf("response message key: got %v want %v", response.Message, "")
	}

	if isValidUUID(response.ID) != true {
		t.Errorf("response id is valid uuid format: got %v want %v", isValidUUID(response.ID), true)
	}

	if response.Amount != body.Amount {
		t.Errorf("response amount key: got %v want %v", response.Amount, body.Amount)
	}

	if response.Currency != body.Currency {
		t.Errorf("response amount key: got %v want %v", response.Currency, body.Currency)
	}
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
