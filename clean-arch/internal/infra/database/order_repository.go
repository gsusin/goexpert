package database

import (
	"database/sql"
	"strconv"

	"github.com/devfullcycle/20-CleanArch/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.Db.QueryRow("Select count(*) from orders").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *OrderRepository) GetOrders() (*[]entity.Order, error) {
	rows, err := r.Db.Query("SELECT id, price, tax, final_price FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := []entity.Order{}
	for rows.Next() {
		var id, price, tax, final_price string
		if err := rows.Scan(&id, &price, &tax, &final_price); err != nil {
			return nil, err
		}
		price_f, _ := strconv.ParseFloat(price, 32)
		tax_f, _ := strconv.ParseFloat(tax, 32)
		final_price_f, _ := strconv.ParseFloat(final_price, 32)
		orders = append(orders, entity.Order{
			ID:         id,
			Price:      price_f,
			Tax:        tax_f,
			FinalPrice: final_price_f,
		})
	}
	return &orders, nil
}
