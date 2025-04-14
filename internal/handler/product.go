package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ulyanovikovak/PVZ_avito/internal/middleware"
	"github.com/ulyanovikovak/PVZ_avito/internal/storage"
)

type AddProductRequest struct {
	Type  string `json:"type"`
	PVZID string `json:"pvzId"`
}

var validTypes = map[string]bool{
	"электроника": true,
	"одежда":      true,
	"обувь":       true,
}

func AddProductHandler(productStore *storage.ProductStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value(middleware.ContextKeyRole).(string)
		if role != "employee" {
			http.Error(w, `{"message":"only employee can add product"}`, http.StatusForbidden)
			return
		}

		var req AddProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !validTypes[req.Type] || req.PVZID == "" {
			http.Error(w, `{"message":"invalid product data"}`, http.StatusBadRequest)
			return
		}

		receptionId, err := productStore.GetLastOpenReception(req.PVZID)
		if err != nil {
			http.Error(w, `{"message":"db error"}`, http.StatusInternalServerError)
			return
		}
		if receptionId == "" {
			http.Error(w, `{"message":"no open reception found"}`, http.StatusBadRequest)
			return
		}

		product, err := productStore.AddProduct(req.Type, receptionId)
		if err != nil {
			http.Error(w, `{"message":"failed to add product"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(product)
	}
}

func DeleteLastProductHandler(productStore *storage.ProductStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value(middleware.ContextKeyRole).(string)
		if role != "employee" {
			http.Error(w, `{"message":"only employee can delete product"}`, http.StatusForbidden)
			return
		}

		pvzId := mux.Vars(r)["pvzId"]
		if pvzId == "" {
			http.Error(w, `{"message":"missing pvzId"}`, http.StatusBadRequest)
			return
		}

		err := productStore.DeleteLastProductFromReception(pvzId)
		if err != nil {
			http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
