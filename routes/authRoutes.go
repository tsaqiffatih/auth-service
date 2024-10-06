package routes

import (
	"github.com/gorilla/mux"
	"github.com/tsaqiffatih/auth-service/handlers"
	"github.com/tsaqiffatih/auth-service/middleware"
	"gorm.io/gorm"
)

func SetupAuthRoutes(r *mux.Router, db *gorm.DB) {
	r.HandleFunc("/register", handlers.RegisterHandler(db)).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler(db)).Methods("POST")

	authRoutes := r.PathPrefix("/auth").Subrouter()
	authRoutes.Use(middleware.AuthMiddleware)
	authRoutes.HandleFunc("/logout", handlers.LogoutHandler(db)).Methods("POST")
}
