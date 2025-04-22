package models

import (
	"time"

	"github.com/your-username/your-repo/internal/database"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	ClerkID   string    `json:"clerk_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetUsers returns all users
func GetUsers(db *database.DB) ([]User, error) {
	rows, err := db.Query(`
		SELECT id, clerk_id, email, name, created_at, updated_at
		FROM users
		ORDER BY email
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.ClerkID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// GetUserByID returns a user by ID
func GetUserByID(db *database.DB, id int) (*User, error) {
	var u User
	err := db.QueryRow(`
		SELECT id, clerk_id, email, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(&u.ID, &u.ClerkID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// CreateUser creates a new user
func CreateUser(db *database.DB, u *User) error {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	return db.QueryRow(`
		INSERT INTO users (clerk_id, email, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, u.ClerkID, u.Email, u.Name, u.CreatedAt, u.UpdatedAt).Scan(&u.ID)
}

// UpdateUser updates a user
func UpdateUser(db *database.DB, u *User) error {
	u.UpdatedAt = time.Now()

	_, err := db.Exec(`
		UPDATE users
		SET email = $1, name = $2, updated_at = $3
		WHERE id = $4
	`, u.Email, u.Name, u.UpdatedAt, u.ID)

	return err
}

// DeleteUser deletes a user
func DeleteUser(db *database.DB, id int) error {
	_, err := db.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}
