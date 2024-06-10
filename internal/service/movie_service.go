package service

import (
	"database/sql"

	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/repository"
)

type MovieService struct {
	MovieRepository *repository.MovieRepository
}

func NewMovieService(db *sql.DB) *MovieService {
	return &MovieService{MovieRepository: repository.NewMovieRepository(db)}
}

func (ms *MovieService) AddMovie(Title, Director, Genre string, Duration int, ImageURL, ImageID, Description, ReleaseDate, EndDate string) (*model.Movie, error) {

	// Thêm bộ phim vào cơ sở dữ liệu với URL ảnh đã tải lên
	movie, err := ms.MovieRepository.AddMovie(Title, Director, Genre, Duration, ImageURL, ImageID, Description, ReleaseDate, EndDate)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (ms *MovieService) GetAllMovie() ([]model.Movie, error) {
	return ms.MovieRepository.GetAllMovie()
}

func (ms *MovieService) GetMovieByID(movieID int) (*model.Movie, error) {
	return ms.MovieRepository.GetMovieByID(movieID)
}

func (ms *MovieService) UpdateMovie(movieID int, Title, Director, Genre string, Duration int, ImageURL, ImageID string) error {
	return ms.MovieRepository.UpdateMovie(movieID, Title, Director, Genre, Duration, ImageURL, ImageID)
}

func (ms *MovieService) DeleteMovie(movieID int) error {
	return ms.MovieRepository.DeleteMovie(movieID)
}
