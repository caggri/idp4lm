package main

import (
	"log"
	"net/http"

	"github.com/caggri/idp4lm/internal/handlers"
	"github.com/caggri/idp4lm/internal/k8s"
)

func main() {
	log.Println("Initializing Kubernetes client...")
	if err := k8s.InitClient(); err != nil {
		log.Printf("Warning: Failed to initialize Kubernetes client: %v", err)
	} else {
		log.Println("Kubernetes client initialized successfully.")
	}

	log.Println("Setting up routes...")
	handlers.SetupRoutes()

	log.Println("Starting server on http://localhost:8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
