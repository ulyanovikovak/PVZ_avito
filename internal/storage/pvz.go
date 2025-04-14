package storage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/ulyanovikovak/PVZ_avito/pkg/models"
)

type PVZStorage struct {
	db *sql.DB
}

func NewPVZStorage(db *sql.DB) *PVZStorage {
	return &PVZStorage{db: db}
}

func (s *PVZStorage) CreatePVZ(city string) (*models.PVZ, error) {
	id := uuid.New().String()
	now := time.Now()
	_, err := s.db.Exec(`INSERT INTO pvz (id, registration_date, city) VALUES ($1, $2, $3)`, id, now, city)
	if err != nil {
		return nil, err
	}
	return &models.PVZ{ID: id, RegistrationDate: now, City: city}, nil
}

func (s *PVZStorage) GetPVZsWithReceptionsAndProducts(start, end *time.Time, page, limit int) ([]map[string]interface{}, error) {
	offset := (page - 1) * limit

	query := `
		SELECT p.id, p.registration_date, p.city,
		       r.id, r.date_time, r.status,
		       pr.id, pr.date_time, pr.type
		FROM pvz p
		LEFT JOIN receptions r ON p.id = r.pvz_id
		LEFT JOIN products pr ON r.id = pr.reception_id
		WHERE ($1::timestamp IS NULL OR r.date_time >= $1)
		  AND ($2::timestamp IS NULL OR r.date_time <= $2)
		ORDER BY p.registration_date DESC, r.date_time DESC, pr.date_time DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := s.db.Query(query, start, end, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type pvzKey struct{ ID string }
	type recKey struct{ ID string }

	result := map[string]map[string]interface{}{}

	for rows.Next() {
		var pvzID, pvzCity string
		var pvzReg time.Time
		var recID, recStatus sql.NullString
		var recDate sql.NullTime
		var prodID, prodType sql.NullString
		var prodDate sql.NullTime

		err := rows.Scan(&pvzID, &pvzReg, &pvzCity, &recID, &recDate, &recStatus, &prodID, &prodDate, &prodType)
		if err != nil {
			return nil, err
		}

		if _, exists := result[pvzID]; !exists {
			result[pvzID] = map[string]interface{}{
				"pvz": map[string]interface{}{
					"id":               pvzID,
					"city":             pvzCity,
					"registrationDate": pvzReg,
				},
				"receptions": []map[string]interface{}{},
			}
		}

		if recID.Valid {
			receptions := result[pvzID]["receptions"].([]map[string]interface{})
			var target map[string]interface{}

			for _, r := range receptions {
				if r["reception"].(map[string]interface{})["id"] == recID.String {
					target = r
					break
				}
			}

			if target == nil {
				target = map[string]interface{}{
					"reception": map[string]interface{}{
						"id":       recID.String,
						"dateTime": recDate.Time,
						"status":   recStatus.String,
						"pvzId":    pvzID,
					},
					"products": []map[string]interface{}{},
				}
				receptions = append(receptions, target)
				result[pvzID]["receptions"] = receptions
			}

			if prodID.Valid {
				products := target["products"].([]map[string]interface{})
				products = append(products, map[string]interface{}{
					"id":          prodID.String,
					"type":        prodType.String,
					"dateTime":    prodDate.Time,
					"receptionId": recID.String,
				})
				target["products"] = products
			}
		}
	}

	final := make([]map[string]interface{}, 0, len(result))
	for _, v := range result {
		final = append(final, v)
	}

	return final, nil
}
