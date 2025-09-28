package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initTracer(ctx context.Context) func() {
	exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint("otel-collector:4318"), otlptracehttp.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create otlp exporter: %v", err)
	}
	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("service-b"),
		),
	)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return func() { _ = tp.Shutdown(ctx) }
}

type CEPRequest struct {
	CEP string `json:"cep"`
}

var (
	viacepBase  = "https://viacep.com.br/ws"
	weatherBase = "http://api.weatherapi.com/v1/current.json"
)

var weatherKey string

func main() {
	ctx := context.Background()
	shutdown := initTracer(ctx)
	defer shutdown()

	_ = godotenv.Load()

	weatherKey = os.Getenv("WEATHERAPI_KEY")
	if weatherKey == "" {
		panic("WEATHERAPI_KEY não definido")
	}

	mux := http.NewServeMux()
	mux.Handle("/zipcode", otelhttp.NewHandler(http.HandlerFunc(handleZipcode), "HandleZipcode"))

	addr := ":8081"
	log.Printf("service-b listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func handleZipcode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/%s/json/", viacepBase, req.CEP))
	if err != nil {
		http.Error(w, "error viacep", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	var via struct {
		Localidade string `json:"localidade"`
		UF         string `json:"uf"`
		Erro       bool   `json:"erro"`
	}
	json.NewDecoder(resp.Body).Decode(&via)
	if via.Erro || via.Localidade == "" {
		// não seria could?
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	q := url.QueryEscape(fmt.Sprintf("%s,%s,BR", via.Localidade, via.UF))
	weatherURL := fmt.Sprintf("%s?key=%s&q=%s", weatherBase, weatherKey, q)
	resp, err = http.Get(weatherURL)
	if err != nil {
		http.Error(w, "error weatherapi", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	var wea struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	json.NewDecoder(resp.Body).Decode(&wea)

	tempC := wea.Current.TempC
	tempF := tempC*1.8 + 32
	tempK := tempC + 273
	out := map[string]string{
		"temp_C": strconv.FormatFloat(tempC, 'f', 1, 64),
		"temp_F": strconv.FormatFloat(tempF, 'f', 1, 64),
		"temp_K": strconv.FormatFloat(tempK, 'f', 1, 64),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

var ErrNotFound = fmt.Errorf("not found")
