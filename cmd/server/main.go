package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	baseHandlers "github.com/phosphene/go_bdd_reference/internal/http"
	"github.com/phosphene/go_bdd_reference/internal/user"
)

func main() {
	// 1. Database connection logic (optional for unit testing environment)
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	var userRepo user.Repository

	if dbHost != "" {
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPass, dbName)

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			fmt.Printf("Error opening database: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()
		userRepo = user.NewPostgresRepository(db)
	}

	// 2. Initialize Hanlers
	regHandler := user.NewRegistrationHandler(userRepo)

	// 3. Routing
	http.HandleFunc("/health", baseHandlers.HealthHandler)
	if userRepo != nil {
		http.Handle("/users", regHandler)
	}

	// 4. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
