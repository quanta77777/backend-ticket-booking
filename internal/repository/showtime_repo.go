package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"movie-ticket-booking/internal/model"
	"time"
)

type ShowTimeRepository struct {
	DB *sql.DB
}

func NewShowTimeRepository(db *sql.DB) *ShowTimeRepository {
	return &ShowTimeRepository{DB: db}
}

func (mr *ShowTimeRepository) TheaterHasShowtimeInTimeRange(theaterID int, startTime, endTime time.Time) (bool, error) {
	// Chuyển đổi startTime và endTime sang múi giờ UTC
	startTimeUTC := startTime.UTC()
	endTimeUTC := endTime.UTC()

	// Thêm khoảng thời gian 15 phút trước và sau suất chiếu
	cleaningBuffer := 5 * time.Minute
	bufferedStartTime := startTimeUTC.Add(-cleaningBuffer)
	bufferedEndTime := endTimeUTC.Add(cleaningBuffer)

	var count int
	query := `
		SELECT COUNT(*) 
		FROM showtime 
		WHERE theater_id = ? 
		AND (
			(start_time BETWEEN ? AND ?) 
			OR (end_time BETWEEN ? AND ?)
			OR (? BETWEEN start_time AND end_time)
			OR (? BETWEEN start_time AND end_time)
			OR (start_time < ? AND end_time > ?)
		)
	`
	err := mr.DB.QueryRow(query, theaterID, bufferedStartTime, bufferedEndTime, bufferedStartTime, bufferedEndTime, startTimeUTC, endTimeUTC, startTimeUTC, endTimeUTC).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (mr *ShowTimeRepository) AddShowtime(branch_id, movie_id, theater_id, cinema_id int, start_time, end_time time.Time) error {

	hasShowtime, err := mr.TheaterHasShowtimeInTimeRange(theater_id, start_time, end_time)

	if err != nil {
		return err
	}
	const errTheaterHasShowtime = "theater already has showtime in the given time range"
	if hasShowtime {

		return errors.New(errTheaterHasShowtime)
	}

	query := "INSERT INTO showtime (branch_id, movie_id, theater_id,cinema_id, start_time, end_time) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = mr.DB.Exec(query, branch_id, movie_id, theater_id, cinema_id, start_time, end_time)

	if err != nil {
		return err
	}

	return nil
}

func (mr *ShowTimeRepository) GetAllShowtimeByDay(day time.Time) ([]model.Showtime, error) {
	rows, err := mr.DB.Query("SELECT * FROM showtime WHERE DATE(start_time) = ?", day.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var showtimes []model.Showtime
	for rows.Next() {
		var showtime model.Showtime
		var startTimeStr, endTimeStr string
		var createdAtBytes []byte
		err := rows.Scan(&showtime.ShowtimeID, &showtime.CinemaID, &showtime.BranchID, &showtime.TheaterID, &showtime.MovieID, &startTimeStr, &endTimeStr, &createdAtBytes)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		showtime.StartTime, err = time.Parse("2006-01-02 15:04:05", startTimeStr)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		showtime.EndTime, err = time.Parse("2006-01-02 15:04:05", endTimeStr)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
		if err != nil {
			return nil, err
		}
		showtime.CreatedAt = createdAt
		showtimes = append(showtimes, showtime)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New(err.Error())
	}
	return showtimes, nil

}

func (mr *ShowTimeRepository) GetShowtimeByDayAndMovieID(day time.Time, movieID int) ([]model.Showtime, error) {

	currentTimePlus5 := time.Now().Add(5 * time.Minute)
	fmt.Printf("Current Time Plus 5 Minutes: %s\n", currentTimePlus5)

	// Thực hiện truy vấn SQL với thời gian đã tính toán
	query := `
        SELECT showtime_id, cinema_id, branch_id, theater_id, movie_id, start_time, end_time, created_at
        FROM showtime
        WHERE start_time >= ? AND movie_id = ? AND DATE(start_time) = ? ORDER BY start_time
    `
	rows, err := mr.DB.Query(query, currentTimePlus5.Format("2006-01-02 15:04:05"), movieID, day.Format("2006-01-02"))
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	var showtimes []model.Showtime
	for rows.Next() {
		var showtime model.Showtime
		var startTimeStr, endTimeStr string
		var createdAtBytes []byte
		err := rows.Scan(&showtime.ShowtimeID, &showtime.CinemaID, &showtime.BranchID, &showtime.TheaterID, &showtime.MovieID, &startTimeStr, &endTimeStr, &createdAtBytes)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		showtime.StartTime, err = time.Parse("2006-01-02 15:04:05", startTimeStr)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		showtime.EndTime, err = time.Parse("2006-01-02 15:04:05", endTimeStr)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		showtime.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
		if err != nil {
			return nil, err
		}
		showtimes = append(showtimes, showtime)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New(err.Error())
	}
	return showtimes, nil
}

func (mr *ShowTimeRepository) GetShowtimWithBranch(branchID int, day time.Time) ([]model.Showtime, error) {
	currentTimePlus5 := time.Now().Add(5 * time.Minute)
	fmt.Printf("Current Time Plus 5 Minutes: %s\n", currentTimePlus5)
	var showtimes []model.Showtime
	rows, err := mr.DB.Query("SELECT showtime_id, cinema_id, branch_id,  theater_id, movie_id, start_time, end_time, created_at FROM showtime WHERE start_time >= ? AND branch_id = ? AND DATE(start_time) = ?  ORDER BY start_time", currentTimePlus5.Format("2006-01-02 15:04:05"), branchID, day.Format("2006-01-02"))
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var showtime model.Showtime
		var startTimeStr, endTimeStr string
		var createdAtBytes []byte
		err := rows.Scan(&showtime.ShowtimeID, &showtime.CinemaID, &showtime.BranchID, &showtime.TheaterID, &showtime.MovieID, &startTimeStr, &endTimeStr, &createdAtBytes)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		showtime.StartTime, err = time.Parse("2006-01-02 15:04:05", startTimeStr)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		showtime.EndTime, err = time.Parse("2006-01-02 15:04:05", endTimeStr)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
		if err != nil {
			return nil, err
		}
		showtime.CreatedAt = createdAt
		showtimes = append(showtimes, showtime)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New(err.Error())
	}

	return showtimes, nil
}

func (mr *ShowTimeRepository) GetShowtimWithCinema(cinemaID int, day time.Time) ([]model.Showtime, error) {
	var showtimes []model.Showtime
	rows, err := mr.DB.Query("SELECT showtime_id, cinema_id, branch_id,  theater_id, movie_id, start_time, end_time, created_at FROM showtime WHERE cinema_id = ? AND DATE(start_time) = ?", cinemaID, day.Format("2006-01-02"))
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var showtime model.Showtime
		var startTimeStr, endTimeStr string
		var createdAtBytes []byte
		err := rows.Scan(&showtime.ShowtimeID, &showtime.CinemaID, &showtime.BranchID, &showtime.TheaterID, &showtime.MovieID, &startTimeStr, &endTimeStr, &createdAtBytes)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		showtime.StartTime, err = time.Parse("2006-01-02 15:04:05", startTimeStr)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		showtime.EndTime, err = time.Parse("2006-01-02 15:04:05", endTimeStr)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
		if err != nil {
			return nil, err
		}
		showtime.CreatedAt = createdAt
		showtimes = append(showtimes, showtime)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New(err.Error())
	}

	return showtimes, nil
}

func (mr *ShowTimeRepository) GetShowtimByID(showtimeID int) ([]model.Showtime, error) {
	var showtimes []model.Showtime
	rows, err := mr.DB.Query("SELECT showtime_id, cinema_id, branch_id,  theater_id, movie_id, start_time, end_time, created_at FROM showtime WHERE showtime_id = ?", showtimeID)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var showtime model.Showtime
		var startTimeStr, endTimeStr string
		var createdAtBytes []byte
		err := rows.Scan(&showtime.ShowtimeID, &showtime.CinemaID, &showtime.BranchID, &showtime.TheaterID, &showtime.MovieID, &startTimeStr, &endTimeStr, &createdAtBytes)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		showtime.StartTime, err = time.Parse("2006-01-02 15:04:05", startTimeStr)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		showtime.EndTime, err = time.Parse("2006-01-02 15:04:05", endTimeStr)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
		if err != nil {
			return nil, err
		}
		showtime.CreatedAt = createdAt
		showtimes = append(showtimes, showtime)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New(err.Error())
	}

	return showtimes, nil
}
