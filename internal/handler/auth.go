package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ulyanovikovak/PVZ_avito/internal/config"
)

type DummyLoginRequest struct {
	Role string `json:"role"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || (req.Role != "employee" && req.Role != "moderator") {
		http.Error(w, `{"message":"invalid role"}`, http.StatusBadRequest)
		return
	}

	claims := jwt.MapClaims{
		"role": req.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(config.JwtSecret)
	if err != nil {
		http.Error(w, `{"message":"signing failed"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(TokenResponse{Token: signed})
}
