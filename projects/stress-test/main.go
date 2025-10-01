package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
}

type Report struct {
	TotalTime     time.Duration
	TotalRequests int
	Status200     int
	StatusCodes   map[int]int
}

func main() {
	// Definir flags CLI
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 0, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")
	flag.Parse()

	// Validar parâmetros
	if *url == "" || *requests <= 0 || *concurrency <= 0 {
		fmt.Println("Uso: go run main.go --url=<URL> --requests=<NUM> --concurrency=<NUM>")
		fmt.Println("Exemplo: go run main.go --url=http://google.com --requests=1000 --concurrency=10")
		return
	}

	fmt.Printf("Iniciando teste de stress...\n")
	fmt.Printf("URL: %s\n", *url)
	fmt.Printf("Requests: %d\n", *requests)
	fmt.Printf("Concorrência: %d\n", *concurrency)
	fmt.Printf("----------------------------------------\n")

	// Executar teste
	results := runStressTest(*url, *requests, *concurrency)

	// Gerar relatório
	report := generateReport(results)
	printReport(report)
}

func runStressTest(url string, totalRequests, concurrency int) []Result {
	results := make([]Result, totalRequests)
	var wg sync.WaitGroup
	resultsChan := make(chan Result, totalRequests)

	// Canal para controlar concorrência com tamanho definido pela flag
	semaphore := make(chan struct{}, concurrency)

	startTime := time.Now()

	for range totalRequests {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Controle de concorrência
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Fazer request
			reqStart := time.Now()
			resp, err := http.Get(url)
			duration := time.Since(reqStart)

			result := Result{
				Duration: duration,
			}

			if err != nil {
				result.StatusCode = 0 // Erro de conexão
			} else {
				result.StatusCode = resp.StatusCode
				resp.Body.Close()
			}

			resultsChan <- result
		}()
	}

	// Aguardar todos os requests
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Coletar resultados
	i := 0
	for result := range resultsChan {
		results[i] = result
		i++
	}

	fmt.Printf("Teste concluído em: %v\n", time.Since(startTime))
	return results
}

func generateReport(results []Result) Report {
	report := Report{
		TotalRequests: len(results),
		StatusCodes:   make(map[int]int),
	}

	var totalDuration time.Duration
	for _, result := range results {
		totalDuration += result.Duration

		if result.StatusCode == 200 {
			report.Status200++
		}

		report.StatusCodes[result.StatusCode]++
	}

	report.TotalTime = totalDuration
	return report
}

func printReport(report Report) {
	fmt.Printf("\n=== RELATÓRIO ===\n")
	fmt.Printf("Tempo total gasto: %v\n", report.TotalTime)
	fmt.Printf("Quantidade total de requests: %d\n", report.TotalRequests)
	fmt.Printf("Requests com status 200: %d\n", report.Status200)
	fmt.Printf("\nDistribuição de códigos de status:\n")

	for status, count := range report.StatusCodes {
		if status == 0 {
			fmt.Printf("  Erro de conexão: %d\n", count)
		} else {
			fmt.Printf("  Status %d: %d\n", status, count)
		}
	}
}
