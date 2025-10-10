package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Estrutura para o retorno do ViaCEP
type ViaCEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	Gia         string `json:"gia"`
	DDD         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Erro        bool   `json:"erro,omitempty"`
}

// Handler para buscar um CEP e gerar cidade.txt
func GetCEP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cep := vars["cep"]
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Erro ao consultar ViaCEP", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler resposta do ViaCEP", http.StatusInternalServerError)
		return
	}

	var cepResp ViaCEPResponse
	if err := json.Unmarshal(body, &cepResp); err != nil {
		http.Error(w, "Erro ao processar resposta", http.StatusInternalServerError)
		return
	}

	if cepResp.Erro {
		http.Error(w, "CEP não encontrado", http.StatusNotFound)
		return
	}

	// Gera o arquivo cidade.txt com os dados desejados
	cidadeInfo := fmt.Sprintf("CEP: %s, Logradouro: %s, Localidade: %s, UF: %s\n", cepResp.CEP, cepResp.Logradouro, cepResp.Localidade, cepResp.UF)
	if err := ioutil.WriteFile("cidade.txt", []byte(cidadeInfo), 0644); err != nil {
		log.Println("Erro ao escrever arquivo cidade.txt:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cepResp)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/cep/{cep}", GetCEP).Methods("GET")
	fmt.Println("Servidor iniciado em :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
