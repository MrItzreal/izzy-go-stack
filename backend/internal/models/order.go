package models

import (
	"time"

	"github.com/your-username/your-repo/internal/database"
)

// Order represents an order in the system
type Order struct {
	ID              int         `json:"id"`
	UserID          int         `json:"user_id"`
	Status          string      `json:"status"`
	Total           int         `json:"total"` // Total in cents
	StripeSessionID string      `json:"stripe_session_id,omitempty"`
	Items           []OrderItem `json:"items,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     int       `json:"price"` // Price at time of purchase in cents
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetOrders returns all orders
func GetOrders(db *database.DB) ([]Order, error) {
	rows, err := db.Query(`
		SELECT id, user_id, status, total, stripe_session_id, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.Total, &o.StripeSessionID, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		
		// Get order items
		items, err := getOrderItems(db, o.ID)
		if err != nil {
			return nil, err
		}
		o.Items = items
		
		orders = append(orders, o)
	}

	return orders, nil
}

// GetOrderByID returns an order by ID
func GetOrderByID(db *database.DB, id int) (*Order, error) {
	var o Order
	err := db.QueryRow(`
		SELECT id, user_id, status, total, stripe_session_id, created_at, updated_at
		FROM orders
		WHERE id = $1
	`, id).Scan(&o.ID, &o.UserID, &o.Status, &o.Total, &o.StripeSessionID, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Get order items
	items, err := getOrderItems(db, o.ID)
	if err != nil {
		return nil, err
	}
	o.Items = items

	return &o, nil
}

// getOrderItems returns all items for an order
func getOrderItems(db *database.DB, orderID int) ([]OrderItem, error) {
	rows, err := db.Query(`
		SELECT id, order_id, product_id, quantity, price, created_at, updated_at
		FROM order_items
		WHERE order_id = $1
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var i OrderItem
		if err := rows.Scan(&i.ID, &i.OrderID, &i.ProductID, &i.Quantity, &i.Price, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, nil
}

// CreateOrder creates a new order
func CreateOrder(db *database.DB, o *Order) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	o.CreatedAt = now
	o.UpdatedAt = now

	// Insert order
	err = tx.QueryRow(`
		INSERT INTO orders (user_id, status, total, stripe_session_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, o.UserID, o.Status, o.Total, o.StripeSessionID, o.CreatedAt, o.UpdatedAt).Scan(&o.ID)
	if err != nil {
		return err
	}

	// Insert order items
	for i := range o.Items {
		item := &o.Items[i]
		item.OrderID = o.ID
		item.CreatedAt = now
		item.UpdatedAt = now

		err = tx.QueryRow(`
			INSERT INTO order_items (order_id, product_id, quantity, price, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, item.OrderID, item.ProductID, item.Quantity, item.Price, item.CreatedAt, item.UpdatedAt).Scan(&item.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// UpdateOrder updates an order
func UpdateOrder(db *database.DB, o *Order) error {
	o.UpdatedAt = time.Now()

	_, err := db.Exec(`
		UPDATE orders
		SET status = $1, stripe_session_id = $2, updated_at = $3
		WHERE id = $4
	`, o.Status, o.StripeSessionID, o.UpdatedAt, o.ID)

	return err
}

// DeleteOrder deletes an order
func DeleteOrder(db *database.DB, id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete order items
	_, err = tx.Exec(`DELETE FROM order_items WHERE order_id = $1`, id)
	if err != nil {
		return err
	}

	// Delete order
	_, err = tx.Exec(`DELETE FROM orders WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
