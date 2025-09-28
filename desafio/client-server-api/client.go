// client.go
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	serverURL     = "http://localhost:8081/cotacao"
	clientTimeout = 300 * time.Millisecond
	outputFile    = "cotacao.txt"
)

type serverReply struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL, nil)
	if err != nil {
		log.Fatalf("erro criando request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("erro realizando request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("server retornou status %d: %s", resp.StatusCode, string(body))
	}

	var sr serverReply
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		log.Fatalf("erro decodificando resposta: %v", err)
	}

	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("erro criando arquivo: %v", err)
	}
	defer f.Close()

	_, _ = f.WriteString("Dólar: " + sr.Bid + "\n")
	log.Printf("Cotação salva em %s: %s", outputFile, sr.Bid)
}
