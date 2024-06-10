package service

import (
	"database/sql"
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/repository"
)

type TheaterService struct {
	TheaterRepository *repository.TheaterRepository
}

func NewTheaterService(db *sql.DB) *TheaterService {
	return &TheaterService{TheaterRepository: repository.NewTheaterRepository(db)}
}

func (ts *TheaterService) AddTheater(branchID int, name string) (int64, error) {
	return ts.TheaterRepository.AddTheater(branchID, name)
}

func (ts *TheaterService) GetTheatersByBranchID(branchID int) ([]model.Theater, error) {
	return ts.TheaterRepository.GetAllTheaterByBranchID(branchID)
}

func (ts *TheaterService) GetTheaterByID(theaterID int) (*model.Theater, error) {
	return ts.TheaterRepository.GetTheaterByID(theaterID)
}
