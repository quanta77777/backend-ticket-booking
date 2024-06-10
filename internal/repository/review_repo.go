package repository

import (
	"database/sql"
	"movie-ticket-booking/internal/model"
	"time"
)

type ReviewRepository interface {
	CreateReview(review model.ReviewRequest) error
	GetReviewsByMovieID(movieID int) ([]model.Review, error)
	UserHasTicketForMovie(userID, movieID int) (bool, error)
	HasUserReviewedMovie(userID, movieID int) (bool, error)
	GetAverageRatingAndCountByMovieID(movieID int) (float64, int, error)
}

type reviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) ReviewRepository {
	return &reviewRepository{db}
}

func (r *reviewRepository) CreateReview(review model.ReviewRequest) error {
	query := `
        INSERT INTO review (user_id, movie_id, rating, comment, image_url, image_id) 
        VALUES (?, ?, ?, ?,?, ?)
    `
	var imageURL, publicID interface{}
	if review.ImageURL != "" {
		imageURL = review.ImageURL
	} else {
		imageURL = nil
	}
	if review.ImageID != "" {
		publicID = review.ImageID
	} else {
		publicID = nil
	}
	_, err := r.db.Exec(query, review.UserID, review.MovieID, review.Rating, review.Comment, imageURL, publicID)
	return err
}

func (r *reviewRepository) GetReviewsByMovieID(movieID int) ([]model.Review, error) {
	query := `
        SELECT review_id, user_id, movie_id, rating, comment, COALESCE(image_url, ''), COALESCE(image_id, ''), created_at 
        FROM review
        WHERE movie_id = ?
    `
	rows, err := r.db.Query(query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []model.Review
	for rows.Next() {
		var review model.Review
		var createdAtBytes []byte
		if err := rows.Scan(&review.ReviewID, &review.UserID, &review.MovieID, &review.Rating, &review.Comment, &review.ImageURL, &review.ImageID, &createdAtBytes); err != nil {
			return nil, err
		}
		createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
		if err != nil {
			return nil, err
		}
		review.CreatedAt = createdAt
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (r *reviewRepository) UserHasTicketForMovie(userID, movieID int) (bool, error) {
	query := `
        SELECT COUNT(*) 
        FROM ticket
        WHERE user_id = ? AND movie_id = ?
    `
	var count int
	err := r.db.QueryRow(query, userID, movieID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *reviewRepository) HasUserReviewedMovie(userID, movieID int) (bool, error) {
	query := `
        SELECT COUNT(*)
        FROM review
        WHERE user_id = ? AND movie_id = ?
    `
	var count int
	err := r.db.QueryRow(query, userID, movieID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *reviewRepository) GetAverageRatingAndCountByMovieID(movieID int) (float64, int, error) {
	query := `
        SELECT COALESCE(AVG(rating), 0), COUNT(*)
        FROM review
        WHERE movie_id = ?
    `
	var avgRating float64
	var count int
	err := r.db.QueryRow(query, movieID).Scan(&avgRating, &count)
	if err != nil {
		return 0, 0, err
	}
	return avgRating, count, nil
}
