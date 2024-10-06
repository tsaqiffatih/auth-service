package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/tsaqiffatih/auth-service/config"
	"github.com/tsaqiffatih/auth-service/routes"
	"github.com/tsaqiffatih/auth-service/utils"
)

func main() {
	err := utils.LoadPrivateKey("private.key")
	if err != nil {
		log.Fatalf("Gagal memuat private key: %v", err)
	}

	err = utils.LoadPublicKey("public.key")
	if err != nil {
		log.Fatalf("Gagal memuat public key: %v", err)
	}

	// Initialize database
	db := config.InitDB()
	log.Println("Database Connected")

	// Setup router
	r := mux.NewRouter()
	routes.SetupAuthRoutes(r, db)

	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
