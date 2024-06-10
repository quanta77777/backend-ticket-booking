package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"movie-ticket-booking/internal/model"
)

type CinemaRepository struct {
	DB *sql.DB
}

func NewCinemaRepository(db *sql.DB) *CinemaRepository {
	return &CinemaRepository{DB: db}
}

func (cr *CinemaRepository) AddCinema(name, ImageURL, PublicID string) error {
	query := "INSERT INTO cinema (name, image_url, image_id) VALUES (?, ?, ?)"
	_, err := cr.DB.Exec(query, name, ImageURL, PublicID)

	if err != nil {
		return err
	}
	return nil
}

func (cr *CinemaRepository) GetAllCinema() ([]model.Cinema, error) {
	rows, err := cr.DB.Query("SELECT * FROM cinema")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cinemas []model.Cinema

	for rows.Next() {
		var cinema model.Cinema
		var createdAtBytes []byte
		err := rows.Scan(&cinema.CinemaID, &cinema.Name, &cinema.ImageURL, &cinema.ImageID, &createdAtBytes)
		if err != nil {
			return nil, err
		}
		createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
		if err != nil {
			return nil, err
		}
		cinema.CreatedAt = createdAt
		cinemas = append(cinemas, cinema)
	}
	return cinemas, nil
}

func (cr *CinemaRepository) GetCinemaByID(cinemaID int) (*model.Cinema, error) {
	var cinema model.Cinema
	var createdAtBytes []byte
	err := cr.DB.QueryRow("SELECT * FROM cinema WHERE cinema_id = ?", cinemaID).
		Scan(&cinema.CinemaID, &cinema.Name, &cinema.ImageURL, &cinema.ImageID, &createdAtBytes)
	if err != nil {
		return nil, err
	}
	createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
	if err != nil {
		return nil, err
	}
	cinema.CreatedAt = createdAt

	return &cinema, nil
}

func (cr *CinemaRepository) UpdateCinema(cinemaID int, name, imageURL, imageID string) error {
	query := "UPDATE cinema SET"
	var args []interface{}
	if name != "" {
		query += " name=?,"
		args = append(args, name)
	}
	if imageURL != "" {
		query += " image_url=?,"
		args = append(args, imageURL)
	}
	if imageID != "" {
		query += " image_id=?,"
		args = append(args, imageID)
	}

	// Loại bỏ dấu phẩy cuối cùng nếu có
	query = strings.TrimSuffix(query, ",")

	// Thêm điều kiện WHERE cho movieID
	query += " WHERE movie_id = ?"
	args = append(args, cinemaID)

	// Thực hiện câu lệnh SQL cập nhật
	_, err := cr.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (cr *CinemaRepository) DeleteCinema(cinemaID int) error {
	if cinemaID <= 0 {
		return errors.New("invalid cinemaID")
	}

	// Xóa chi nhánh trước
	branchQuery := "DELETE FROM branch WHERE cinema_id = ?"
	_, err := cr.DB.Exec(branchQuery, cinemaID)
	if err != nil {
		return fmt.Errorf("failed to delete branches: %v", err)
	}

	// Tiếp tục xóa chuỗi rạp
	cinemaQuery := "DELETE FROM cinema WHERE cinema_id = ?"
	_, err = cr.DB.Exec(cinemaQuery, cinemaID)
	if err != nil {
		return fmt.Errorf("failed to delete cinema: %v", err)
	}

	return nil
}
