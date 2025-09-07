package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
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
		fmt.Println("Returned existing short URL")
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

// Authentication handlers

// LoginRequest represents the login request structure
type LoginRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

// SignupRequest represents the signup request structure
type SignupRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response structure
type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	UserID  string `json:"user_id,omitempty"`
}

// Login handles POST /api/login
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login called")
	var req LoginRequest

	// Parse JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Println("Parsed login request for user:", req.UserID)

	// Validate input
	if req.UserID == "" || req.Password == "" {
		response := AuthResponse{
			Success: false,
			Message: "User ID and password are required",
		}
		fmt.Println("Validation failed: missing user ID or password")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Println("User Id and Password is fine:")

	// Get user from database
	user, err := h.db.GetUserByUserID(req.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			response := AuthResponse{
				Success: false,
				Message: "Invalid user ID or password",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}
		fmt.Println("Database error during login:", err)
		response := AuthResponse{
			Success: false,
			Message: "Internal server error",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Println("User found with ID:", user.UserID)

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		response := AuthResponse{
			Success: false,
			Message: "Invalid user ID or password",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Println("Password verified for user:", req.UserID)

	// Generate simple token (in production, use JWT or similar)
	token := fmt.Sprintf("token-%s-%d", req.UserID, time.Now().Unix())

	response := AuthResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		UserID:  req.UserID,
	}

	fmt.Println("Login successful for user:", req.UserID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Signup handles POST /api/signup
func (h *Handlers) Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signup called")
	var req SignupRequest

	// Parse JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Println("Parsed signup request for user:", req.UserID)

	// Validate input
	if req.UserID == "" || req.Password == "" {
		response := AuthResponse{
			Success: false,
			Message: "User ID and password are required",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Println("User Id and Password is fine:")

	// Basic validation for user ID (alphanumeric, minimum length)
	if len(req.UserID) < 3 {
		response := AuthResponse{
			Success: false,
			Message: "User ID must be at least 3 characters long",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Println("validated user ID:", req.UserID)

	// Basic validation for password (minimum length)
	if len(req.Password) < 6 {
		response := AuthResponse{
			Success: false,
			Message: "Password must be at least 6 characters long",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Println("validated password length for user:", req.UserID)

	// Check if user already exists
	exists, err := h.db.UserExists(req.UserID)
	if err != nil {
		fmt.Println("Database error checking user existence:", err)
		response := AuthResponse{
			Success: false,
			Message: "Internal server error",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Println("User existence check for", req.UserID, ":", exists)

	if exists {
		response := AuthResponse{
			Success: false,
			Message: "User ID already exists",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Println("User ID is available:", req.UserID)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		response := AuthResponse{
			Success: false,
			Message: "Internal server error",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Println("Password hashed successfully for user:", req.UserID)

	// Create user in database
	user, err := h.db.CreateUser(req.UserID, string(hashedPassword))
	if err != nil {
		fmt.Println("Error creating user:", err)
		response := AuthResponse{
			Success: false,
			Message: "Failed to create user",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Println("User created with ID:", user.ID)

	// Generate simple token (in production, use JWT or similar)
	token := fmt.Sprintf("token-%s-%d", req.UserID, time.Now().Unix())

	response := AuthResponse{
		Success: true,
		Message: "Signup successful",
		Token:   token,
		UserID:  user.UserID,
	}

	fmt.Println("Signup successful for user:", req.UserID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
