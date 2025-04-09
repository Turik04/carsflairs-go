package orders

import (
    "carsflairs-backend/pkg/db"
    "encoding/json"
    "net/http"
)

type Order struct {
    ID       int     `json:"id"`
    UserID   int     `json:"user_id"`
    FrameID  int     `json:"frame_id"`
    Quantity int     `json:"quantity"`
    Total    float64 `json:"total"`
}

// Создать заказ
func CreateOrder(w http.ResponseWriter, r *http.Request) {
    var o Order
    if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err := db.DB.Exec("INSERT INTO orders (user_id, frame_id, quantity, total) VALUES ($1, $2, $3, $4)",
        o.UserID, o.FrameID, o.Quantity, o.Total)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}
