# ShortURL App

A full-stack URL shortener application built with Go (backend) and React (frontend). Instantly create short, memorable links for any URL and track usage.

## Features

- Generate short URLs for any valid link
- Automatic collision-free short code generation
- Click tracking (analytics)
- Optional expiration for short URLs
- Modern React frontend with instant feedback
- RESTful API backend with PostgreSQL

---

## Project Structure

```
shorturl-app/
  backend/    # Go backend API
  frontend/   # React frontend app
```

---

## Backend (Go)

- **API Endpoints:**
  - `POST /api/shorten` — Create a new short URL
  - `GET /{shortCode}` — Redirect to the original URL and increment click count

- **Tech Stack:** Go, Gorilla Mux, PostgreSQL, CORS, dotenv

- **Environment:**
  - Configure your database in `.env` or via `DATABASE_URL`
  - Example: `DATABASE_URL=postgres://username:password@localhost/shorturl_db?sslmode=disable`

- **Setup:**
  1. Install Go and PostgreSQL.
  2. Create a database:  
     `createdb shorturl_db`
  3. Copy `.env.example` to `.env` and set your DB credentials.
  4. Run the backend:
     ```sh
     cd backend
     go run .c
     ```

---

## Frontend (React)

- **Features:**  
  - Simple form to enter a long URL and get a short one
  - Displays errors and loading states
  - Uses Axios for API calls

- **Setup:**
  1. Install Node.js and npm.
  2. Install dependencies:
     ```sh
     cd frontend
     npm install
     ```
  3. Start the frontend:
     ```sh
     npm start
     ```
  4. The app runs at [http://localhost:3000](http://localhost:3000)

---

## Usage

1. Start both backend (`:8080`) and frontend (`:3000`).
2. Open the frontend in your browser.
3. Enter a long URL and click "Shorten".
4. Use the generated short URL to be redirected and increment click count.

---

## Customization

- **Short Code Length:**  
  Change the length in `backend/utils.go` (`GenerateRandomCode(6)`).
- **Allowed Origins:**  
  Update CORS settings in `backend/main.go`.
- **Expiration:**  
  Set expiration days in the frontend form (if enabled).

---

## License

MIT

---

Let me know if you want to add deployment, Docker, or more advanced usage instructions!
