package service

import (
	"database/sql"
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/repository"
	"time"
)

type ShowtimeService struct {
	ShowtimeRepository *repository.ShowTimeRepository
}

func NewShowtimeService(db *sql.DB) *ShowtimeService {
	return &ShowtimeService{ShowtimeRepository: repository.NewShowTimeRepository(db)}
}

func (ms *ShowtimeService) AddShowtime(branch_id, movie_id, theater_id, cinema_id int, start_time, end_time time.Time) error {
	return ms.ShowtimeRepository.AddShowtime(branch_id, movie_id, theater_id, cinema_id, start_time, end_time)
}

func (ms *ShowtimeService) GetAllShowtimeByDay(day time.Time) ([]model.Showtime, error) {
	return ms.ShowtimeRepository.GetAllShowtimeByDay(day)
}

func (ms *ShowtimeService) GetShowtimeByDayAndMovieID(day time.Time, movieID int) ([]model.Showtime, error) {
	return ms.ShowtimeRepository.GetShowtimeByDayAndMovieID(day, movieID)
}

func (ms *ShowtimeService) GetShowtimeWithBranch(branchID int, day time.Time) ([]model.Showtime, error) {
	return ms.ShowtimeRepository.GetShowtimWithBranch(branchID, day)
}

func (ms *ShowtimeService) GetShowtimeWithCinema(cinemaID int, day time.Time) ([]model.Showtime, error) {
	return ms.ShowtimeRepository.GetShowtimWithCinema(cinemaID, day)
}

func (ms *ShowtimeService) GetShowtimByID(showtimeID int) ([]model.Showtime, error) {
	return ms.ShowtimeRepository.GetShowtimByID(showtimeID)
}
