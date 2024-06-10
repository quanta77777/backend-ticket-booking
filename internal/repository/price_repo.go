package repository

import (
	"database/sql"
	"movie-ticket-booking/internal/model"
)

type PriceRepository struct {
	DB *sql.DB
}

func NewPriceRepository(db *sql.DB) *PriceRepository {
	return &PriceRepository{db}
}

func (pr *PriceRepository) GetPricesByShowtimeID(showtimeID int) (map[string]int64, error) {
	prices := make(map[string]int64)
	query := `SELECT seat_type, price FROM price WHERE showtime_id = ?`
	rows, err := pr.DB.Query(query, showtimeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var seatType string
		var price int64
		if err := rows.Scan(&seatType, &price); err != nil {
			return nil, err
		}
		prices[seatType] = price
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prices, nil
}

func (pr *PriceRepository) CreatePriceForShowtime(req model.Price) error {
	query := `INSERT INTO price (showtime_id, seat_type, price) VALUES (?, ?, ?)`
	_, err := pr.DB.Exec(query, req.ShowtimeID, req.SeatType, req.Price)
	return err
}
