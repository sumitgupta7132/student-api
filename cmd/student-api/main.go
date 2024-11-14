package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sumitgupta7132/student-api/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// database setup
	// setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Student api"))
	})
	// setup server

	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}
	fmt.Printf("Server started at %s\n", cfg.Address)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start server")
	}

}
