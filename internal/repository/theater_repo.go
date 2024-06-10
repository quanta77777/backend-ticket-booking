package repository

import (
	"database/sql"

	"movie-ticket-booking/internal/model"
)

type TheaterRepository struct {
	DB *sql.DB
}

func NewTheaterRepository(db *sql.DB) *TheaterRepository {
	return &TheaterRepository{DB: db}
}

func (tr *TheaterRepository) AddTheater(branch_id int, name string) (int64, error) {
	query := "INSERT INTO theater (branch_id, name) VALUES ( ?, ?)"
	result, err := tr.DB.Exec(query, branch_id, name)

	if err != nil {
		return 0, err
	}
	theaterID, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return theaterID, nil
}

func (cr *TheaterRepository) GetAllTheaterByBranchID(branch_id int) ([]model.Theater, error) {
	rows, err := cr.DB.Query("SELECT theater_id, branch_id, name FROM theater WHERE branch_id = ?", branch_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var theaters []model.Theater
	for rows.Next() {
		var theater model.Theater
		err := rows.Scan(&theater.TheaterId, &theater.BranchId, &theater.Name)
		if err != nil {
			return nil, err
		}
		theaters = append(theaters, theater)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return theaters, nil
}

func (cr *TheaterRepository) GetTheaterByID(theaterID int) (*model.Theater, error) {
	var theater model.Theater
	err := cr.DB.QueryRow("SELECT * FROM theater WHERE theater_id = ?", theaterID).
		Scan(&theater.TheaterId, &theater.BranchId, &theater.Name)
	if err != nil {
		return nil, err
	}
	return &theater, nil
}
