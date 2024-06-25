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

func (us *UserService) GetUserByEmail(email string) (*model.User, error) {
	return us.UserRepository.GetUserByEmail(email)
}

func (us *UserService) CreateUser(user model.UserRequest) error {
	return us.UserRepository.CreateUser(user)
}
