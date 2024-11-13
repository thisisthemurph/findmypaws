package main

import (
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	dbConnStr := os.Getenv("DATABASE_CONNECTION_STRING")
	if dbConnStr == "" {
		panic("No database connection string found")
	}
	db, err := sqlx.Open("postgres", dbConnStr)
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	manager := NewManager(db, logger)

	// CORS middleware to allow specific headers, including 'anonymoususerid'
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},           // Allow the frontend origin
		AllowedHeaders: []string{"Content-Type", "anonymoususerid"}, // Allow specific headers
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},          // Allow methods as needed
	})

	// Route handler for /room
	http.HandleFunc("/room", func(w http.ResponseWriter, req *http.Request) {
		// Get the Room ID from the request (e.g., from a query parameter).
		participantID := req.URL.Query().Get("pid")
		if participantID == "" {
			http.Error(w, "Missing required parameter pid", http.StatusBadRequest)
			return
		}
		roomID := req.URL.Query().Get("r")
		if roomID == "" {
			http.Error(w, "Missing Room ID", http.StatusBadRequest)
			return
		}

		roomIdentifier, err := uuid.Parse(roomID)
		if err != nil {
			http.Error(w, "Invalid Room ID", http.StatusBadRequest)
			return
		}

		// Retrieve or create a new Room for the given Room ID.
		r, err := manager.GetOrCreateRoom(roomIdentifier, participantID)
		if err != nil {
			if errors.Is(err, ErrUnauthorized) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		r.ServeWS(w, req)
	})

	// Wrap the default HTTP handler with the CORS handler
	http.Handle("/", corsHandler.Handler(http.DefaultServeMux))

	log.Println("Listening on :42096")
	if err := http.ListenAndServe(":42096", nil); err != nil {
		log.Fatal(err)
	}
}
