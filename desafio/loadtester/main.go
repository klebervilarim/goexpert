package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	// Parâmetros da CLI
	url := flag.String("url", "", "URL do serviço a ser testado")
	totalRequests := flag.Int("requests", 1, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" {
		fmt.Println("Erro: --url é obrigatório")
		return
	}

	fmt.Printf("Iniciando teste de carga em %s\n", *url)
	fmt.Printf("Total de requests: %d, Concurrency: %d\n", *totalRequests, *concurrency)

	start := time.Now()

	statusCounts := make(map[int]int)
	var mu sync.Mutex

	jobs := make(chan int, *totalRequests)
	wg := sync.WaitGroup{}

	// Workers
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				resp, err := http.Get(*url)
				code := 0
				if err != nil {
					code = -1
				} else {
					code = resp.StatusCode
					resp.Body.Close()
				}

				mu.Lock()
				statusCounts[code]++
				mu.Unlock()
			}
		}()
	}

	// Enviar jobs
	for i := 0; i < *totalRequests; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	duration := time.Since(start)

	// Relatório
	fmt.Println("\n===== Relatório =====")
	fmt.Printf("Tempo total: %s\n", duration)
	fmt.Printf("Total de requests: %d\n", *totalRequests)
	fmt.Printf("Requests com status 200: %d\n", statusCounts[200])
	fmt.Println("Distribuição de outros códigos:")
	for code, count := range statusCounts {
		if code != 200 {
			if code == -1 {
				fmt.Printf("Erros de requisição: %d\n", count)
			} else {
				fmt.Printf("Status %d: %d\n", code, count)
			}
		}
	}
}
