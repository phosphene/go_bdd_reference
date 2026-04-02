package main

import (
	"fmt"
	"net/http"
	"os"

	handlers "github.com/phosphene/go_bdd_reference/internal/http"
)

func main() {
	http.HandleFunc("/health", handlers.HealthHandler)

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
