# Short URL Generator

A comprehensive URL shortener application built with Go backend and React frontend, featuring user authentication and advanced URL management.

## Features

- **User Authentication**: Secure signup and login system
- **URL Shortening**: Transform long URLs into short, shareable links with base62 encoding
- **URL Visiting**: Visit short URLs safely to preview destination before redirecting
- **User Dashboard**: Personalized interface for authenticated users
- **URL Expiration**: Set custom expiration dates for shortened URLs
- **Click Tracking**: Monitor usage statistics for your URLs
- **Responsive UI**: Clean, modern interface that works on all devices
- **PostgreSQL Database**: Reliable data storage with proper indexing
- **RESTful API**: Well-structured API endpoints for all functionality

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

#### Authentication
- `POST /api/signup` - Create a new user account
  ```json
  {
    "user_id": "your_username",
    "password": "your_password"
  }
  ```

- `POST /api/login` - Authenticate user
  ```json
  {
    "user_id": "your_username", 
    "password": "your_password"
  }
  ```

#### URL Management
- `POST /api/shorten` - Create a short URL
  ```json
  {
    "url": "https://example.com",
    "expires_in_days": 30  // optional
  }
  ```

- `GET /{shortCode}` - Redirect to original URL (increments click count)

### Database Schema

The application uses two main tables:

#### `short_urls` table:
- `id` - Primary key (SERIAL)
- `short_code` - The short URL code (VARCHAR(10))
- `original_url` - The original long URL (TEXT)
- `created_at` - Creation timestamp
- `expires_at` - Expiration timestamp (optional)
- `click_count` - Number of clicks

#### `users` table:
- `id` - Primary key (SERIAL)
- `user_id` - Unique username (VARCHAR(255))
- `password` - Hashed password (VARCHAR(255))
- `created_at` - Account creation timestamp
- `updated_at` - Last update timestamp

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

### Application Features

The React frontend provides several key features:

1. **User Authentication**
   - Secure signup and login forms
   - Session management
   - User-specific dashboard

2. **URL Shortening**
   - Create short URLs from long ones
   - Set optional expiration dates
   - Copy shortened URLs to clipboard

3. **URL Visiting**
   - Safe preview of short URLs before redirecting
   - Extract short codes from various URL formats
   - Open original URLs in new tabs

4. **Dashboard**
   - Personalized welcome interface
   - Quick access to main features
   - Easy navigation between functions

## Project Structure

```
shorturl-app/
├── backend/
│   ├── main.go          # Main application with routing
│   ├── models.go        # Database models and methods
│   ├── handlers.go      # HTTP handlers for all endpoints
│   ├── utils.go         # URL encoding/decoding utilities
│   ├── go.mod           # Go dependencies
│   └── .env             # Environment variables
├── frontend/
│   ├── src/
│   │   ├── components/  # React components
│   │   │   ├── Auth.js           # Login/Signup component
│   │   │   ├── Dashboard.js      # User dashboard
│   │   │   ├── CreateShortUrl.js # URL shortening
│   │   │   ├── VisitShortUrl.js  # URL visiting
│   │   │   └── UrlShortener.js   # Main URL shortener
│   │   ├── App.js       # Main app component
│   │   └── ...
│   └── package.json     # Node dependencies
└── README.md
```

## Environment Variables

- `DATABASE_URL` - PostgreSQL connection string (required)
- `PORT` - Server port (default: 8080)

## Getting Started

1. **Set up the database**: Create a PostgreSQL database and configure the connection string
2. **Start the backend**: Navigate to the backend directory and run `go run .`
3. **Start the frontend**: Navigate to the frontend directory and run `npm start`
4. **Create an account**: Sign up with a username and password
5. **Start shortening**: Use the dashboard to create and manage your short URLs

## Security Features

- Password hashing for secure user authentication
- Input validation and sanitization
- CORS configuration for secure cross-origin requests
- SQL injection prevention with parameterized queries

## Usage Flow

1. **Sign up/Login**: Create an account or log in with existing credentials
2. **Dashboard**: Access your personalized dashboard
3. **Create Short URLs**: Enter long URLs to generate short versions
4. **Visit URLs**: Safely preview short URLs before being redirected
5. **Track Usage**: Monitor click counts for your shortened URLs
