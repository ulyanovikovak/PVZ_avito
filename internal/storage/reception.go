package storage

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ulyanovikovak/PVZ_avito/pkg/models"
)

type ReceptionStorage struct {
	db *sql.DB
}

func NewReceptionStorage(db *sql.DB) *ReceptionStorage {
	return &ReceptionStorage{db: db}
}

func (s *ReceptionStorage) HasOpenReception(pvzId string) (bool, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM receptions WHERE pvz_id = $1 AND status = 'in_progress'`, pvzId).Scan(&count)
	return count > 0, err
}

func (s *ReceptionStorage) CreateReception(pvzId string) (*models.Reception, error) {
	id := uuid.New().String()
	now := time.Now()
	_, err := s.db.Exec(`INSERT INTO receptions (id, pvz_id, date_time, status) VALUES ($1, $2, $3, 'in_progress')`, id, pvzId, now)
	if err != nil {
		return nil, err
	}
	return &models.Reception{ID: id, DateTime: now, PVZID: pvzId, Status: "in_progress"}, nil
}

func (s *ReceptionStorage) CloseLastReception(pvzId string) (*models.Reception, error) {
	var id string
	var date time.Time

	err := s.db.QueryRow(`
		SELECT id, date_time FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY date_time DESC
		LIMIT 1
	`, pvzId).Scan(&id, &date)

	if err == sql.ErrNoRows {
		return nil, errors.New("no open reception")
	} else if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(`UPDATE receptions SET status = 'close' WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	return &models.Reception{
		ID:       id,
		DateTime: date,
		PVZID:    pvzId,
		Status:   "close",
	}, nil
}
