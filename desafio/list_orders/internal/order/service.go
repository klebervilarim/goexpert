package order

// Service é a camada de negócio da aplicação.
// Ela usa o Repository para acessar dados do banco.

type Service struct {
	Repo *Repository
}

// ListOrders retorna todas as orders do banco.
func (s *Service) ListOrders() ([]Order, error) {
	return s.Repo.List()
}

// CreateOrder cria uma nova order no banco.
func (s *Service) CreateOrder(customer string, amount float64) error {
	return s.Repo.Create(customer, amount)
}
