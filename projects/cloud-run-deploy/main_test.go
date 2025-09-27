package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func doReq(url string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	handler(w, req)
	return w
}

func TestInvalidCEP(t *testing.T) {
	w := doReq("/weather?cep=123")
	if w.Code != 422 {
		t.Fatalf("esperado 422, veio %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "invalid zipcode") {
		t.Fatalf("resposta errada: %s", w.Body.String())
	}
}

func TestCEPNotFound(t *testing.T) {
	viacepBase = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"erro": true}`))
	})).URL
	defer func() { viacepBase = "https://viacep.com.br/ws" }()

	w := doReq("/weather?cep=00000000")
	if w.Code != 404 {
		t.Fatalf("esperado 404, veio %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "can not find zipcode") {
		t.Fatalf("resposta errada: %s", w.Body.String())
	}
}

func TestSuccess(t *testing.T) {
	// mock do ViaCEP
	viacepSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"localidade":"São Paulo","uf":"SP"}`))
	}))
	defer viacepSrv.Close()
	viacepBase = viacepSrv.URL

	// mock do WeatherAPI
	weatherSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"current":{"temp_c":25.0}}`))
	}))
	defer weatherSrv.Close()
	weatherBase = weatherSrv.URL

	w := doReq("/weather?cep=01001000")
	if w.Code != 200 {
		t.Fatalf("esperado 200, veio %d. body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"temp_C":"25.0"`) {
		t.Fatalf("resposta errada: %s", w.Body.String())
	}
}
