package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://username:password@localhost/shorturl_db?sslmode=disable"
		log.Println("Using default database URL. Set DATABASE_URL environment variable for production.")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Create database instance and table
	database := NewDatabase(db)
	if err := database.CreateTable(); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Create handlers
	appHandlers := NewHandlers(database)

	// Create router
	r := mux.NewRouter()

	// Add CORS middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)(r)

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/shorten", appHandlers.ShortenURL).Methods("POST")

	// Redirect route (catch-all for short codes)
	r.PathPrefix("/").HandlerFunc(appHandlers.RedirectURL)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Printf("API endpoint: http://localhost:%s/api/shorten\n", port)
	fmt.Printf("Redirect endpoint: http://localhost:%s/{shortCode}\n", port)

	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
