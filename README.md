# Short URL Generator

A simple URL shortener built with Go backend and React frontend.

## Features

- Shorten any URL with base62 encoding
- URL expiration support
- Click tracking
- Clean, responsive UI
- PostgreSQL database
- RESTful API

## Backend (Go)

### Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up PostgreSQL database:
   - Create a database named `shorturl_db`
   - Set the `DATABASE_URL` environment variable:
     ```bash
     export DATABASE_URL="postgres://username:password@localhost/shorturl_db?sslmode=disable"
     ```

4. Run the backend:
   ```bash
   go run .
   ```

The backend will run on `http://localhost:8080`

### API Endpoints

- `POST /api/shorten` - Create a short URL
  ```json
  {
    "url": "https://example.com",
    "expires_in_days": 30  // optional
  }
  ```

- `GET /{shortCode}` - Redirect to original URL

### Database Schema

The `short_urls` table includes:
- `id` - Primary key (SERIAL)
- `short_code` - The short URL code (VARCHAR(10))
- `original_url` - The original long URL (TEXT)
- `created_at` - Creation timestamp
- `expires_at` - Expiration timestamp (optional)
- `click_count` - Number of clicks

## Frontend (React)

### Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm start
   ```

The frontend will run on `http://localhost:3000`

## Project Structure

```
shorturl-app/
├── backend/
│   ├── main.go          # Main application
│   ├── models.go        # Database models and methods
│   ├── handlers.go      # HTTP handlers
│   ├── utils.go         # URL encoding/decoding
│   └── go.mod           # Go dependencies
├── frontend/
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── App.js       # Main app component
│   │   └── ...
│   └── package.json     # Node dependencies
└── README.md
```

## Environment Variables

- `DATABASE_URL` - PostgreSQL connection string
- `PORT` - Server port (default: 8080)
