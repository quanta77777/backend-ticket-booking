package repository

import (
	"database/sql"
	"movie-ticket-booking/internal/model"
	"strings"
	"time"
)

type BranchRepository struct {
	DB *sql.DB
}

func NewBranchRepository(db *sql.DB) *BranchRepository {
	return &BranchRepository{DB: db}
}

func (br *BranchRepository) AddBranch(cinema_id int, name, address, imageURl, imageID string) (int64, error) {
	query := "INSERT INTO branch (cinema_id, name, address, image_url, image_id) VALUES (?, ?, ?, ?, ?)"
	result, err := br.DB.Exec(query, cinema_id, name, address, imageURl, imageID)

	if err != nil {
		return 0, err
	}
	branchID, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return branchID, nil
}

func (br *BranchRepository) GetAllBranch() ([]model.Branch, error) {
	rows, err := br.DB.Query("SELECT * FROM branch")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var branches []model.Branch

	for rows.Next() {
		var branch model.Branch
		var createdAtBytes []byte
		err := rows.Scan(&branch.BranchId, &branch.CinemaId, &branch.Name, &branch.Address, &branch.ImageURL, &branch.ImageID, &createdAtBytes)
		if err != nil {
			return nil, err
		}
		createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
		if err != nil {
			return nil, err
		}
		branch.CreatedAt = createdAt

		branches = append(branches, branch)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return branches, nil
}

func (br *BranchRepository) GetBranchByID(branchID int) (*model.Branch, error) {
	var branch model.Branch
	var createdAtBytes []byte
	err := br.DB.QueryRow("SELECT * FROM branch WHERE branch_id = ?", branchID).Scan(&branch.BranchId, &branch.CinemaId, &branch.Name, &branch.Address, &branch.ImageURL, &branch.ImageID, &createdAtBytes)
	if err != nil {
		return nil, err
	}
	createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
	if err != nil {
		return nil, err
	}
	branch.CreatedAt = createdAt

	return &branch, nil
}

func (br *BranchRepository) GetBranchesByCinemaID(cinemaID int) ([]model.Branch, error) {
	rows, err := br.DB.Query("SELECT * FROM branch WHERE cinema_id = ?", cinemaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var branches []model.Branch
	for rows.Next() {
		var branch model.Branch
		var createdAtBytes []byte
		err := rows.Scan(&branch.BranchId, &branch.CinemaId, &branch.Name, &branch.Address, &branch.ImageURL, &branch.ImageID, &createdAtBytes)
		if err != nil {
			return nil, err
		}
		createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
		if err != nil {
			return nil, err
		}
		branch.CreatedAt = createdAt
		branches = append(branches, branch)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return branches, nil
}

func (br *BranchRepository) UpdateBranch(branchID int, name, address, imageURL, imageID string) error {
	query := "UPDATE cinema_chain SET"
	var args []interface{}
	if name != "" {
		query += " name=?,"
		args = append(args, name)
	}
	if address != "" {
		query += " address=?,"
		args = append(args, address)
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
	query += " WHERE branch_id = ?"
	args = append(args, branchID)

	// Thực hiện câu lệnh SQL cập nhật
	_, err := br.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (br *BranchRepository) DeleteBranch(branchID int) error {
	query := "DELETE FROM branch WHERE branch_id = ?"

	_, err := br.DB.Exec(query, branchID)
	if err != nil {
		return err
	}
	return nil
}
