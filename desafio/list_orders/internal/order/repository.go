package order

import (
	"database/sql"
)

type Repository struct {
	DB *sql.DB
}

func (r *Repository) List() ([]Order, error) {
	rows, err := r.DB.Query("SELECT id, customer, amount FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.Customer, &o.Amount); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *Repository) Create(customer string, amount float64) error {
	_, err := r.DB.Exec("INSERT INTO orders (customer, amount) VALUES ($1, $2)", customer, amount)
	return err
}
