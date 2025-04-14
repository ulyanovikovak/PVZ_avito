package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ulyanovikovak/PVZ_avito/internal/middleware"
	"github.com/ulyanovikovak/PVZ_avito/internal/storage"
)

type CreateReceptionRequest struct {
	PVZID string `json:"pvzId"`
}

func CreateReceptionHandler(store *storage.ReceptionStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value(middleware.ContextKeyRole).(string)
		if role != "employee" {
			http.Error(w, `{"message":"only employee can create reception"}`, http.StatusForbidden)
			return
		}

		var req CreateReceptionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PVZID == "" {
			http.Error(w, `{"message":"invalid request"}`, http.StatusBadRequest)
			return
		}

		exists, err := store.HasOpenReception(req.PVZID)
		if err != nil {
			http.Error(w, `{"message":"db error"}`, http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, `{"message":"open reception exists"}`, http.StatusBadRequest)
			return
		}

		reception, err := store.CreateReception(req.PVZID)
		if err != nil {
			http.Error(w, `{"message":"creation error"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(reception)
	}
}

func CloseLastReceptionHandler(store *storage.ReceptionStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value(middleware.ContextKeyRole).(string)
		if role != "employee" {
			http.Error(w, `{"message":"only employee can close reception"}`, http.StatusForbidden)
			return
		}

		pvzId := mux.Vars(r)["pvzId"]
		if pvzId == "" {
			http.Error(w, `{"message":"missing pvzId"}`, http.StatusBadRequest)
			return
		}

		reception, err := store.CloseLastReception(pvzId)
		if err != nil {
			http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reception)
	}
}
