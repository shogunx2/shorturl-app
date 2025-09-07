package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Handlers contains the HTTP handlers
type Handlers struct {
	db *Database
}

// NewHandlers creates a new handlers instance
func NewHandlers(db *Database) *Handlers {
	return &Handlers{db: db}
}

// ShortenURL handles POST /api/shorten
func (h *Handlers) ShortenURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ShortenURL called")
	var req ShortenRequest

	// Parse JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate URL
	if req.URL == "" {
		fmt.Println("URL is required")
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Normalize URL
	normalizedURL := NormalizeURL(req.URL)

	// Check if URL already exists
	existing, err := h.db.GetByOriginalURL(normalizedURL)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("Database error checking existing URL:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if existing != nil {
		// Return existing short URL
		fmt.Println("Found existing URL:", existing.ShortCode)
		response := ShortenResponse{
			ShortURL:    fmt.Sprintf("http://%s/%s", r.Host, existing.ShortCode),
			OriginalURL: existing.OriginalURL,
			CreatedAt:   existing.CreatedAt,
			ExpiresAt:   existing.ExpiresAt,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate short code before inserting
	var shortCode string
	for {
		shortCode = GenerateRandomCode(6)
		// Check for collision
		exists, err := h.db.GetByShortCode(shortCode)
		if err != nil && err != sql.ErrNoRows {
			fmt.Println("Database error checking short code:", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if exists == nil {
			break // unique code
		}
	}

	shortURL := &ShortURL{
		OriginalURL: normalizedURL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(),
		ClickCount:  0,
	}

	if req.ExpiresInDays != nil {
		expiresAt := time.Now().Add(time.Duration(*req.ExpiresInDays) * 24 * time.Hour)
		shortURL.ExpiresAt = &expiresAt
	}

	// Insert into database
	fmt.Println("Creating new short URL for:", normalizedURL)
	id, err := h.db.Create(shortURL)
	if err != nil {
		fmt.Println("Error creating short URL:", err)
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}
	fmt.Println("Created short URL with ID:", id)
	shortURL.ID = int(id)

	// Return response
	response := ShortenResponse{
		ShortURL:    fmt.Sprintf("http://%s/%s", r.Host, shortCode),
		OriginalURL: shortURL.OriginalURL,
		CreatedAt:   shortURL.CreatedAt,
		ExpiresAt:   shortURL.ExpiresAt,
	}

	fmt.Println("ShortURL created:", response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// RedirectURL handles GET /{shortCode}
func (h *Handlers) RedirectURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RedirectURL called")
	shortCode := r.URL.Path[1:] // Remove leading slash

	if shortCode == "" {
		fmt.Println("Error: Short code is required")
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	// Get URL from database by short code
	fmt.Println("Looking up URL for short code:", shortCode)
	shortURL, err := h.db.GetByShortCode(shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("URL not found for short code:", shortCode)
			h.renderErrorPage(w, "URL not found", http.StatusNotFound)
		} else {
			fmt.Println("Database error looking up URL:", err)
			h.renderErrorPage(w, "Database error", http.StatusInternalServerError)
		}
		return
	}
	fmt.Println("Found URL:", shortURL.OriginalURL)

	// Check if URL has expired
	if shortURL.ExpiresAt != nil && shortURL.ExpiresAt.Before(time.Now()) {
		fmt.Println("URL has expired:", shortURL.ExpiresAt)
		h.renderErrorPage(w, "URL has expired", http.StatusNotFound)
		return
	}

	// Increment click count
	if err := h.db.IncrementClickCount(shortURL.ID); err != nil {
		// Log error but don't fail the redirect
		fmt.Printf("Failed to increment click count: %v\n", err)
	} else {
		fmt.Println("Successfully incremented click count")
	}

	// Redirect to original URL
	fmt.Println("Redirecting to:", shortURL.OriginalURL)
	http.Redirect(w, r, shortURL.OriginalURL, http.StatusMovedPermanently)
}

// renderErrorPage renders a simple HTML error page
func (h *Handlers) renderErrorPage(w http.ResponseWriter, message string, statusCode int) {
	fmt.Printf("Rendering error page: %d - %s\n", statusCode, message)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>%d - Error</title>
    <style>
        body { 
            font-family: Arial, sans-serif; 
            text-align: center; 
            padding: 50px; 
            background-color: #f8f9fa;
        }
        .error { 
            color: #dc3545; 
            font-size: 4rem;
            margin-bottom: 1rem;
        }
        .message {
            color: #6c757d;
            font-size: 1.2rem;
        }
    </style>
</head>
<body>
    <h1 class="error">%d</h1>
    <p class="message">%s</p>
</body>
</html>`, statusCode, statusCode, message)

	fmt.Fprint(w, html)
}
