package service

import (
	"database/sql"

	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/repository"
)

type CinemaService struct {
	CinemaRepository *repository.CinemaRepository
}

func NewCinemaService(db *sql.DB) *CinemaService {
	return &CinemaService{CinemaRepository: repository.NewCinemaRepository(db)}
}

func (cs *CinemaService) AddCinema(name, ImageURL, PublicID string) error {
	return cs.CinemaRepository.AddCinema(name, ImageURL, PublicID)
}

func (cs *CinemaService) GetAllCinema() ([]model.Cinema, error) {
	return cs.CinemaRepository.GetAllCinema()
}
func (cs *CinemaService) GetCinemaByID(cinemaID int) (*model.Cinema, error) {
	return cs.CinemaRepository.GetCinemaByID(cinemaID)
}

func (cs *CinemaService) UpdateCinema(cinemaID int, name, imageURL, imageID string) error {
	return cs.CinemaRepository.UpdateCinema(cinemaID, name, imageURL, imageID)
}

func (cs *CinemaService) DeleteCinema(cinemaID int) error {
	return cs.CinemaRepository.DeleteCinema(cinemaID)
}
