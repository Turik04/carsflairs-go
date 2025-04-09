package main

import (
	"carsflairs-backend/pkg/auth"
	"carsflairs-backend/pkg/db"
	"carsflairs-backend/pkg/favorites"
	"carsflairs-backend/pkg/frames"
	"carsflairs-backend/pkg/orders"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db.ConnectDB()
	defer db.CloseDB()

	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/register", auth.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")

	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(auth.AuthMiddleware)

	protected.HandleFunc("/frames", frames.GetFrames).Methods("GET")
	protected.HandleFunc("/frames", frames.CreateFrame).Methods("POST")
	protected.HandleFunc("/orders", orders.CreateOrder).Methods("POST")
	protected.HandleFunc("/favorites/{frame_id:[0-9]+}", favorites.AddToFavorites).Methods("POST")
	protected.HandleFunc("/favorites", favorites.GetFavorites).Methods("GET")
	protected.HandleFunc("/frames", auth.RoleMiddleware("Admin", frames.UpdateFrame)).Methods("PUT")
	protected.HandleFunc("/frames", frames.GetFrames).Methods("GET")
	protected.HandleFunc("/frames", frames.CreateFrame).Methods("POST")
	protected.HandleFunc("/frames/{id:[0-9]+}", frames.UpdateFrame).Methods("PUT")
	protected.HandleFunc("/frames/{id:[0-9]+}", frames.DeleteFrame).Methods("DELETE")

	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				return
			}
			h.ServeHTTP(w, r)
		})
	})

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Carsflairs API!")
}
