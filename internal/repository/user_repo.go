package repository

import (
	"database/sql"
	"errors"
	"movie-ticket-booking/internal/model"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) GetAllUsers() ([]model.User, error) {
	rows, err := ur.DB.Query("SELECT user_id, name, email, role, created_at FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.UserID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (ur *UserRepository) GetUserByID(userID string) (*model.User, error) {
	var user model.User
	err := ur.DB.QueryRow("SELECT user_id, name, email, role, COALESCE(image_url, ''), COALESCE(image_id, ''), created_at FROM user WHERE user_id = ?", userID).
		Scan(&user.UserID, &user.Name, &user.Email, &user.Role, &user.ImageUrl, &user.ImageID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) CreateUser(user model.UserRequest) error {
	query := "INSERT INTO user (name, email, password, image_url, image_id) VALUES (?, ?, ?, ?, ?)"
	var imageURL, publicID interface{}
	if user.ImageURL != "" {
		imageURL = user.ImageURL
	} else {
		imageURL = nil
	}
	if user.ImageID != "" {
		publicID = user.ImageID
	} else {
		publicID = nil
	}
	_, err := ur.DB.Exec(query, &user.Name, &user.Email, &user.Password, imageURL, publicID)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	query := `
        SELECT user_id, name, email, password, role, COALESCE(image_url, ''), COALESCE(image_id, '')
        FROM user
        WHERE email = ?
    `
	row := r.DB.QueryRow(query, email)

	var user model.User
	if err := row.Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.Role, &user.ImageUrl, &user.ImageID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
