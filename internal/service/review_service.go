package service

import (
	"fmt"
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/repository"
)

type ReviewService interface {
	CreateReview(review model.ReviewRequest) error
	GetReviewsByMovieID(movieID int) ([]model.Review, error)
	GetAverageRatingAndCountByMovieID(movieID int) (float64, int, error)
}

type reviewService struct {
	reviewRepo repository.ReviewRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository) ReviewService {
	return &reviewService{reviewRepo}
}

func (s *reviewService) CreateReview(review model.ReviewRequest) error {
	// Kiểm tra xem user đã mua vé cho phim chưa
	hasTicket, err := s.reviewRepo.UserHasTicketForMovie(review.UserID, review.MovieID)
	if err != nil {
		return err
	}

	if !hasTicket {
		return fmt.Errorf("user has not purchased a ticket for this movie")
	}

	// Kiểm tra xem user đã đánh giá phim chưa
	hasReviewed, err := s.reviewRepo.HasUserReviewedMovie(review.UserID, review.MovieID)
	if err != nil {
		return err
	}

	if hasReviewed {
		return fmt.Errorf("user has already reviewed this movie")
	}

	return s.reviewRepo.CreateReview(review)
}

func (s *reviewService) GetReviewsByMovieID(movieID int) ([]model.Review, error) {
	return s.reviewRepo.GetReviewsByMovieID(movieID)
}

func (s *reviewService) GetAverageRatingAndCountByMovieID(movieID int) (float64, int, error) {
	return s.reviewRepo.GetAverageRatingAndCountByMovieID(movieID)
}
