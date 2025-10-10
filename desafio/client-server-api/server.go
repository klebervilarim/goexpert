package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

const (
	apiURL       = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	apiTimeout   = 200 * time.Millisecond
	dbTimeout    = 10 * time.Millisecond
	sqliteDBFile = "cotacoes.db"
)

type apiResponse map[string]struct {
	Bid        string `json:"bid"`
	CreateDate string `json:"create_date"`
}

type reply struct {
	Bid string `json:"bid"`
}

func main() {
	db, err := sql.Open("sqlite", sqliteDBFile)
	if err != nil {
		log.Fatalf("erro abrindo sqlite: %v", err)
	}
	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT NOT NULL,
		create_date TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatalf("erro criando tabela: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		handleCotacao(w, r, db)
	})

	// tenta porta 8080, se não, usa 8081
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("porta %d ocupada, usando 8081", port)
		port = 8081
		listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatalf("não foi possível abrir porta 8081: %v", err)
		}
	}

	log.Printf("Servidor iniciado em :%d", port)
	if err := http.Serve(listener, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handleCotacao(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	reqCtx := r.Context()

	apiCtx, cancelAPI := context.WithTimeout(reqCtx, apiTimeout)
	defer cancelAPI()

	req, err := http.NewRequestWithContext(apiCtx, http.MethodGet, apiURL, nil)
	if err != nil {
		log.Printf("erro criando request API: %v", err)
		http.Error(w, "erro interno", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("erro ao chamar API: %v", err)
		http.Error(w, "erro ao obter cotação", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		log.Printf("erro decodificando API: %v", err)
		http.Error(w, "erro parse resposta", http.StatusInternalServerError)
		return
	}

	data, ok := apiResp["USDBRL"]
	if !ok {
		for _, v := range apiResp {
			data = v
			break
		}
	}

	bid := data.Bid
	createDate := data.CreateDate

	dbCtx, cancelDB := context.WithTimeout(reqCtx, dbTimeout)
	defer cancelDB()
	_, err = db.ExecContext(dbCtx, "INSERT INTO cotacoes (bid, create_date) VALUES (?, ?)", bid, createDate)
	if err != nil {
		log.Printf("erro ao persistir no DB: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply{Bid: bid})
}
