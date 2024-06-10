package repository

import (
	"database/sql"
	"movie-ticket-booking/internal/model"
	"strings"
	"time"
)

type MovieRepository struct {
	DB *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{DB: db}
}

func (mr *MovieRepository) AddMovie(Title, Director, Genre string, Duration int, ImageURL, ImageID, Description, ReleaseDate, EndDate string) (*model.Movie, error) {

	query := "INSERT INTO movie (title, director, genre,duration, image_url, image_id, description, release_date, end_date ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := mr.DB.Exec(query, Title, Director, Genre, Duration, ImageURL, ImageID, Description, ReleaseDate, EndDate)
	if err != nil {
		return nil, err
	}
	movieID, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	movie := &model.Movie{
		MovieID:     int(movieID),
		Title:       Title,
		Director:    Director,
		Genre:       Genre,
		Duration:    Duration,
		ImageURL:    ImageURL,
		ImageID:     ImageID,
		Discription: Description,
		ReleaseDate: ReleaseDate,
		EndDate:     EndDate,
	}

	return movie, nil
}

func (mr *MovieRepository) GetAllMovie() ([]model.Movie, error) {
	rows, err := mr.DB.Query("SELECT * FROM movie")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var movies []model.Movie
	for rows.Next() {
		var movie model.Movie
		var createdAtBytes []byte
		err := rows.Scan(&movie.MovieID, &movie.Title, &movie.Director, &movie.Genre, &movie.Duration, &movie.ImageURL, &movie.ImageID, &movie.Discription, &movie.ReleaseDate, &movie.EndDate, &createdAtBytes)
		if err != nil {
			return nil, err
		}
		createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
		if err != nil {
			return nil, err
		}
		movie.CreatedAt = createdAt
		movies = append(movies, movie)
	}
	return movies, nil
}

func (mr *MovieRepository) GetMovieByID(movieID int) (*model.Movie, error) {
	var movie model.Movie
	var createdAtBytes []byte
	err := mr.DB.QueryRow("SELECT * FROM movie WHERE movie_id = ?", movieID).
		Scan(&movie.MovieID, &movie.Title, &movie.Director, &movie.Genre, &movie.Duration, &movie.ImageURL, &movie.ImageID, &movie.Discription, &movie.ReleaseDate, &movie.EndDate, &createdAtBytes)
	if err != nil {
		return nil, err
	}
	createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
	if err != nil {
		return nil, err
	}
	movie.CreatedAt = createdAt
	return &movie, nil
}

func (mr *MovieRepository) UpdateMovie(movieID int, Title, Director, Genre string, Duration int, ImageURL, ImageID string) error {

	query := "UPDATE movie SET"
	var args []interface{}
	if Title != "" {
		query += " title=?,"
		args = append(args, Title)
	}
	if Director != "" {
		query += " director=?,"
		args = append(args, Director)
	}
	if Genre != "" {
		query += " genre=?,"
		args = append(args, Genre)
	}
	if Duration != 0 {
		query += " duration=?,"
		args = append(args, Duration)
	}
	if ImageURL != "" {
		query += "image_url=?,"
		args = append(args, ImageURL)
	}
	if ImageID != "" {
		query += "image_id=?,"
		args = append(args, ImageID)
	}

	// Loại bỏ dấu phẩy cuối cùng nếu có
	query = strings.TrimSuffix(query, ",")

	// Thêm điều kiện WHERE cho movieID
	query += " WHERE movie_id = ?"
	args = append(args, movieID)

	// Thực hiện câu lệnh SQL cập nhật
	_, err := mr.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (mr *MovieRepository) DeleteMovie(movieID int) error {
	// Tạo câu lệnh SQL DELETE để xóa bộ phim với ID tương ứng
	query := "DELETE FROM movie WHERE movie_id = ?"

	// Thực thi câu lệnh SQL
	_, err := mr.DB.Exec(query, movieID)
	if err != nil {
		return err
	}

	return nil
}
