package auth

import (
	"carsflairs-backend/pkg/db"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("yourjwtsecret") // тот же ключ!

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Claims struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// Регистрация (всегда с ролью "User")
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	role := "User"
	_, err := db.DB.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)", user.Username, hashedPassword, role)
	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "User registered successfully")
}

// Логин
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var stored User
	err := db.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username=$1", user.Username).Scan(
		&stored.ID, &stored.Username, &stored.Password, &stored.Role,
	)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(user.Password)) != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Если это Turik — автоматически делаем его админом
	if stored.Username == "Turik" && stored.Role != "Admin" {
		_, err := db.DB.Exec("UPDATE users SET role=$1 WHERE id=$2", "Admin", stored.ID)
		if err == nil {
			stored.Role = "Admin"
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		ID:   stored.ID,
		Role: stored.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	})
	tokenString, _ := token.SignedString(secretKey)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
