package storage

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ulyanovikovak/PVZ_avito/pkg/models"
)

type ProductStorage struct {
	db *sql.DB
}

func NewProductStorage(db *sql.DB) *ProductStorage {
	return &ProductStorage{db: db}
}

func (s *ProductStorage) GetLastOpenReception(pvzId string) (string, error) {
	var id string
	err := s.db.QueryRow(`
		SELECT id FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY date_time DESC
		LIMIT 1
	`, pvzId).Scan(&id)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return id, err
}

func (s *ProductStorage) AddProduct(productType string, receptionId string) (*models.Product, error) {
	id := uuid.New().String()
	now := time.Now()

	_, err := s.db.Exec(`
		INSERT INTO products (id, date_time, type, reception_id)
		VALUES ($1, $2, $3, $4)
	`, id, now, productType, receptionId)
	if err != nil {
		return nil, err
	}

	return &models.Product{
		ID:          id,
		DateTime:    now,
		Type:        productType,
		ReceptionID: receptionId,
	}, nil
}

func (s *ProductStorage) DeleteLastProductFromReception(pvzId string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var receptionId string
	err = tx.QueryRow(`
		SELECT id FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY date_time DESC LIMIT 1
	`, pvzId).Scan(&receptionId)
	if err == sql.ErrNoRows {
		return errors.New("no open reception")
	} else if err != nil {
		return err
	}

	var productId string
	err = tx.QueryRow(`
		SELECT id FROM products
		WHERE reception_id = $1
		ORDER BY date_time DESC
		LIMIT 1
	`, receptionId).Scan(&productId)
	if err == sql.ErrNoRows {
		return errors.New("no products in reception")
	} else if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM products WHERE id = $1`, productId)
	if err != nil {
		return err
	}

	return tx.Commit()
}
