package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Estrutura padrão para exibição de CEP
type CepResponse struct {
	Logradouro string
	Bairro     string
	Cidade     string
	UF         string
}

type Result struct {
	Source string
	Data   CepResponse
	Err    error
}

func main() {
	// Receber CEP como argumento de linha de comando
	cep := flag.String("cep", "", "CEP a ser consultado")
	flag.Parse()

	if *cep == "" {
		log.Fatal("Informe o CEP usando: -cep=05372100")
	}

	// Contexto com timeout de 1 segundo
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	results := make(chan Result, 2)

	// Chamar BrasilAPI
	go func() {
		url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", *cep)
		var resp struct {
			Street       string `json:"street"`
			Neighborhood string `json:"neighborhood"`
			City         string `json:"city"`
			State        string `json:"state"`
		}
		err := fetch(ctx, url, &resp)
		results <- Result{
			Source: "BrasilAPI",
			Data: CepResponse{
				Logradouro: resp.Street,
				Bairro:     resp.Neighborhood,
				Cidade:     resp.City,
				UF:         resp.State,
			},
			Err: err,
		}
	}()

	// Chamar ViaCEP
	go func() {
		url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", *cep)
		var resp struct {
			Logradouro string `json:"logradouro"`
			Bairro     string `json:"bairro"`
			Localidade string `json:"localidade"`
			UF         string `json:"uf"`
		}
		err := fetch(ctx, url, &resp)
		results <- Result{
			Source: "ViaCEP",
			Data: CepResponse{
				Logradouro: resp.Logradouro,
				Bairro:     resp.Bairro,
				Cidade:     resp.Localidade,
				UF:         resp.UF,
			},
			Err: err,
		}
	}()

	// Pegar a primeira resposta que chegar
	select {
	case res := <-results:
		if res.Err != nil {
			log.Fatalf("Erro na API %s: %v", res.Source, res.Err)
		}
		fmt.Printf("Resposta recebida da API: %s\n", res.Source)
		fmt.Printf("Logradouro: %s\nBairro: %s\nCidade: %s\nUF: %s\n",
			res.Data.Logradouro, res.Data.Bairro, res.Data.Cidade, res.Data.UF)
	case <-ctx.Done():
		log.Fatal("Timeout de 1 segundo atingido")
	}
}

// Função genérica para buscar JSON de URL
func fetch(ctx context.Context, url string, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}
