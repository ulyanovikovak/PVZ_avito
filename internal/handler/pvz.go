package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ulyanovikovak/PVZ_avito/internal/middleware"
	"github.com/ulyanovikovak/PVZ_avito/internal/storage"
)

type CreatePVZRequest struct {
	City string `json:"city"`
}

var allowedCities = map[string]bool{
	"Москва":          true,
	"Санкт-Петербург": true,
	"Казань":          true,
}

func CreatePVZHandler(pvzStore *storage.PVZStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value(middleware.ContextKeyRole).(string)
		if role != "moderator" {
			http.Error(w, `{"message":"only moderator can create PVZ"}`, http.StatusForbidden)
			return
		}

		var req CreatePVZRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !allowedCities[req.City] {
			http.Error(w, `{"message":"invalid city"}`, http.StatusBadRequest)
			return
		}

		pvz, err := pvzStore.CreatePVZ(req.City)
		if err != nil {
			http.Error(w, `{"message":"failed to create PVZ"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(pvz)
	}
}

func GetPVZsHandler(pvzStore *storage.PVZStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value(middleware.ContextKeyRole).(string)
		if role != "moderator" && role != "employee" {
			http.Error(w, `{"message":"access denied"}`, http.StatusForbidden)
			return
		}

		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		limit, _ := strconv.Atoi(q.Get("limit"))
		if page <= 0 {
			page = 1
		}
		if limit <= 0 || limit > 30 {
			limit = 10
		}

		var startPtr, endPtr *time.Time
		if sd := q.Get("startDate"); sd != "" {
			if t, err := time.Parse(time.RFC3339, sd); err == nil {
				startPtr = &t
			}
		}
		if ed := q.Get("endDate"); ed != "" {
			if t, err := time.Parse(time.RFC3339, ed); err == nil {
				endPtr = &t
			}
		}

		list, err := pvzStore.GetPVZsWithReceptionsAndProducts(startPtr, endPtr, page, limit)
		if err != nil {
			http.Error(w, `{"message":"db error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	}
}
