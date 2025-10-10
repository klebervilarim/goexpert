package order

import (
	"encoding/json"
	"net/http"
)

func RESTHandler(repo *Repository) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			orders, err := repo.List()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			json.NewEncoder(w).Encode(orders)
			return
		}

		if r.Method == http.MethodPost {
			var o Order
			if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			if err := repo.Create(o.Customer, o.Amount); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.WriteHeader(http.StatusCreated)
			return
		}
	})

	return mux
}
