package models

import (
	"time"

	"github.com/your-username/your-repo/internal/database"
)

// Product represents a product in the system
type Product struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          int       `json:"price"` // Price in cents
	StripeProductID string    `json:"stripe_product_id,omitempty"`
	StripePriceID  string    `json:"stripe_price_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// GetProducts returns all products
func GetProducts(db *database.DB) ([]Product, error) {
	rows, err := db.Query(`
		SELECT id, name, description, price, stripe_product_id, stripe_price_id, created_at, updated_at
		FROM products
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.StripeProductID, &p.StripePriceID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// GetProductByID returns a product by ID
func GetProductByID(db *database.DB, id int) (*Product, error) {
	var p Product
	err := db.QueryRow(`
		SELECT id, name, description, price, stripe_product_id, stripe_price_id, created_at, updated_at
		FROM products
		WHERE id = $1
	`, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.StripeProductID, &p.StripePriceID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// CreateProduct creates a new product
func CreateProduct(db *database.DB, p *Product) error {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	return db.QueryRow(`
		INSERT INTO products (name, description, price, stripe_product_id, stripe_price_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, p.Name, p.Description, p.Price, p.StripeProductID, p.StripePriceID, p.CreatedAt, p.UpdatedAt).Scan(&p.ID)
}

// UpdateProduct updates a product
func UpdateProduct(db *database.DB, p *Product) error {
	p.UpdatedAt = time.Now()

	_, err := db.Exec(`
		UPDATE products
		SET name = $1, description = $2, price = $3, stripe_product_id = $4, stripe_price_id = $5, updated_at = $6
		WHERE id = $7
	`, p.Name, p.Description, p.Price, p.StripeProductID, p.StripePriceID, p.UpdatedAt, p.ID)

	return err
}

// DeleteProduct deletes a product
func DeleteProduct(db *database.DB, id int) error {
	_, err := db.Exec(`DELETE FROM products WHERE id = $1`, id)
	return err
}
