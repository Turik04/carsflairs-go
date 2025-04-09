package frames

import (
    "carsflairs-backend/pkg/db"
    "encoding/json"
    "net/http"
		_ "github.com/lib/pq"
		"fmt"

		"github.com/gorilla/mux"
)

type Frame struct {
    ID       int     `json:"id"`
    Model    string  `json:"model"`
    Price    float64 `json:"price"`
    Size     string  `json:"size"`
    Material string  `json:"material"`
    Image    string  `json:"image"`
}

func GetFrames(w http.ResponseWriter, r *http.Request) {
    rows, err := db.DB.Query("SELECT id, model, price, size, material, image FROM frames")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var frames []Frame
    for rows.Next() {
        var f Frame
        if err := rows.Scan(&f.ID, &f.Model, &f.Price, &f.Size, &f.Material, &f.Image); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        frames = append(frames, f)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(frames)
}

func CreateFrame(w http.ResponseWriter, r *http.Request) {
	var f Frame
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec("INSERT INTO frames (model, price, size, material, image) VALUES ($1, $2, $3, $4, $5)",
		f.Model, f.Price, f.Size, f.Material, f.Image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Frame created successfully",
	})
}



func UpdateFrame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("Updating frame with ID:", id) // Логирование

	var f Frame
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
	}

	_, err := db.DB.Exec(
			"UPDATE frames SET model=$1, price=$2, size=$3, material=$4, image=$5 WHERE id=$6",
			f.Model, f.Price, f.Size, f.Material, f.Image, id,
	)
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Frame updated successfully"))
}



func DeleteFrame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.DB.Exec("DELETE FROM frames WHERE id=$1", id)
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Frame deleted successfully"))
}
