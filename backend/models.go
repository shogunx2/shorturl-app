package main

import (
	"database/sql"
	"fmt"
	"time"
)

// ShortURL represents a shortened URL in the database
type ShortURL struct {
	ID          int        `json:"id" db:"id"`
	ShortCode   string     `json:"short_code" db:"short_code"`
	OriginalURL string     `json:"original_url" db:"original_url"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at" db:"expires_at"`
	ClickCount  int        `json:"click_count" db:"click_count"`
}

// ShortenRequest represents the request body for shortening a URL
type ShortenRequest struct {
	URL           string `json:"url"`
	ExpiresInDays *int   `json:"expires_in_days,omitempty"`
}

// ShortenResponse represents the response for a shortened URL
type ShortenResponse struct {
	ShortURL    string     `json:"short_url"`
	OriginalURL string     `json:"original_url"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// Database represents the database connection
type Database struct {
	conn *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(conn *sql.DB) *Database {
	return &Database{conn: conn}
}

// Create inserts a new short URL and returns its ID
func (db *Database) Create(shortURL *ShortURL) (int64, error) {
	fmt.Println("Create called with shortURL:", shortURL)
	query := `
		INSERT INTO short_urls (short_code, original_url, created_at, expires_at, click_count)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	var id int64
	err := db.conn.QueryRow(query, shortURL.ShortCode, shortURL.OriginalURL, shortURL.CreatedAt, shortURL.ExpiresAt, shortURL.ClickCount).Scan(&id)
	return id, err
}

// GetByID retrieves a short URL by its ID
func (db *Database) GetByID(id int) (*ShortURL, error) {
	query := `SELECT id, short_code, original_url, created_at, expires_at, click_count FROM short_urls WHERE id = $1`

	shortURL := &ShortURL{}
	err := db.conn.QueryRow(query, id).Scan(
		&shortURL.ID,
		&shortURL.ShortCode,
		&shortURL.OriginalURL,
		&shortURL.CreatedAt,
		&shortURL.ExpiresAt,
		&shortURL.ClickCount,
	)

	if err != nil {
		return nil, err
	}

	return shortURL, nil
}

// GetByOriginalURL retrieves a short URL by its original URL
func (db *Database) GetByOriginalURL(originalURL string) (*ShortURL, error) {
	query := `SELECT id, short_code, original_url, created_at, expires_at, click_count FROM short_urls WHERE original_url = $1`

	shortURL := &ShortURL{}
	err := db.conn.QueryRow(query, originalURL).Scan(
		&shortURL.ID,
		&shortURL.ShortCode,
		&shortURL.OriginalURL,
		&shortURL.CreatedAt,
		&shortURL.ExpiresAt,
		&shortURL.ClickCount,
	)

	if err != nil {
		return nil, err
	}

	return shortURL, nil
}

// UpdateShortCode updates the short code for a given ID
func (db *Database) UpdateShortCode(id int, shortCode string) error {
	query := `UPDATE short_urls SET short_code = $1 WHERE id = $2`
	_, err := db.conn.Exec(query, shortCode, id)
	return err
}

// IncrementClickCount increments the click count for a given ID
func (db *Database) IncrementClickCount(id int) error {
	query := `UPDATE short_urls SET click_count = click_count + 1 WHERE id = $1`
	_, err := db.conn.Exec(query, id)
	return err
}

// GetByShortCode retrieves a short URL by its short code
func (db *Database) GetByShortCode(shortCode string) (*ShortURL, error) {
	query := `SELECT id, short_code, original_url, created_at, expires_at, click_count FROM short_urls WHERE short_code = $1`

	shortURL := &ShortURL{}
	err := db.conn.QueryRow(query, shortCode).Scan(
		&shortURL.ID,
		&shortURL.ShortCode,
		&shortURL.OriginalURL,
		&shortURL.CreatedAt,
		&shortURL.ExpiresAt,
		&shortURL.ClickCount,
	)

	if err != nil {
		return nil, err
	}

	return shortURL, nil
}

// CreateTable creates the short_urls table if it doesn't exist
func (db *Database) CreateTable() error {
	fmt.Println("CreateTable called")
	query := `
		CREATE TABLE IF NOT EXISTS short_urls (
			id SERIAL PRIMARY KEY,
			short_code VARCHAR(10) UNIQUE NOT NULL,
			original_url TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			expires_at TIMESTAMP,
			click_count INTEGER NOT NULL DEFAULT 0
		);
		
		CREATE INDEX IF NOT EXISTS idx_short_code ON short_urls(short_code);
		CREATE INDEX IF NOT EXISTS idx_original_url ON short_urls(original_url);
	`

	_, err := db.conn.Exec(query)
	return err
}

// User represents a user in the database
type User struct {
	ID        int       `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserTable creates the users table if it doesn't exist
func (db *Database) CreateUserTable() error {
	fmt.Println("CreateUserTable called")
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			user_id VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
		
		CREATE INDEX IF NOT EXISTS idx_user_id ON users(user_id);
	`

	_, err := db.conn.Exec(query)
	return err
}

// CreateUser inserts a new user and returns the user struct
func (db *Database) CreateUser(userID, hashedPassword string) (*User, error) {
	fmt.Println("CreateUser called with user_id:", userID)
	now := time.Now()
	query := `
		INSERT INTO users (user_id, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, password, created_at, updated_at`

	user := &User{}
	err := db.conn.QueryRow(query, userID, hashedPassword, now, now).Scan(
		&user.ID,
		&user.UserID,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUserID retrieves a user by their user_id
func (db *Database) GetUserByUserID(userID string) (*User, error) {
	query := `SELECT id, user_id, password, created_at, updated_at FROM users WHERE user_id = $1`

	user := &User{}
	err := db.conn.QueryRow(query, userID).Scan(
		&user.ID,
		&user.UserID,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UserExists checks if a user with the given user_id exists
func (db *Database) UserExists(userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`

	var exists bool
	err := db.conn.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
