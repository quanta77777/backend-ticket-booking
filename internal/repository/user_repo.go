package repository

import (
	"database/sql"
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
	err := ur.DB.QueryRow("SELECT user_id, name, email, role, created_at FROM user WHERE user_id = ?", userID).
		Scan(&user.UserID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) AddUser(name, email, password string) error {
	query := "INSERT INTO user (name, email, password) VALUES (?, ?, ?)"
	_, err := ur.DB.Exec(query, name, email, password)

	if err != nil {
		return err
	}

	return nil
}
