package service

import (
	"database/sql"
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/repository"
)

type UserService struct {
	UserRepository *repository.UserRepository
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		UserRepository: repository.NewUserRepository(db),
	}
}

func (us *UserService) GetAllUsers() ([]model.User, error) {
	return us.UserRepository.GetAllUsers()
}

func (us *UserService) GetUserByID(userID string) (*model.User, error) {
	return us.UserRepository.GetUserByID(userID)
}

func (us *UserService) AddUser(name, email, password string) error {
	return us.UserRepository.AddUser(name, email, password)
}
