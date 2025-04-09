package favorites

import (
	"carsflairs-backend/pkg/db"
	"carsflairs-backend/pkg/frames"
	"carsflairs-backend/pkg/auth" 
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
  "errors"

	"github.com/gorilla/mux"
	"github.com/golang-jwt/jwt"
)

func getUserIDFromToken(tokenString string) (int, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &auth.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("yourjwtsecret"), nil
	})

	if err != nil || !token.Valid {
			return 0, errors.New("invalid token")
	}

	return claims.ID, nil
}

func AddToFavorites(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	userID, err := getUserIDFromToken(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	frameID, err := strconv.Atoi(vars["frame_id"])
	if err != nil {
		http.Error(w, "Invalid frame ID", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("INSERT INTO favorites (user_id, frame_id) VALUES ($1, $2) ON CONFLICT (user_id, frame_id) DO NOTHING", userID, frameID)
	if err != nil {
		http.Error(w, "Error adding to favorites", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Frame added to favorites"))
}

func GetFavorites(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	userID, err := getUserIDFromToken(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.DB.Query("SELECT f.id, f.model, f.price, f.size, f.material, f.image FROM frames f JOIN favorites fav ON f.id = fav.frame_id WHERE fav.user_id = $1", userID)
	if err != nil {
		http.Error(w, "Error retrieving favorites", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var favorites []frames.Frame
	for rows.Next() {
		var frame frames.Frame
		if err := rows.Scan(&frame.ID, &frame.Model, &frame.Price, &frame.Size, &frame.Material, &frame.Image); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		favorites = append(favorites, frame)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(favorites)
}
