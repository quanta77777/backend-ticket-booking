package service

import (
	"errors"
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	accessSecret  = []byte("access_secret")
	refreshSecret = []byte("refresh_secret")
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *AuthService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) CreateToken(email, role string) (*model.TokenDetails, error) {
	accessTokenExpiry := time.Now().Add(time.Minute * 15).Unix()
	refreshTokenExpiry := time.Now().Add(time.Hour * 24 * 7).Unix()

	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["email"] = email
	accessTokenClaims["role"] = role
	accessTokenClaims["exp"] = accessTokenExpiry
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	accessTokenString, err := accessToken.SignedString(accessSecret)
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["email"] = email
	refreshTokenClaims["role"] = role
	refreshTokenClaims["exp"] = refreshTokenExpiry
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	refreshTokenString, err := refreshToken.SignedString(refreshSecret)
	if err != nil {
		return nil, err
	}

	return &model.TokenDetails{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, errors.New("token expired")
		}
		return token, nil
	}

	return nil, err
}

func (s *AuthService) RefreshToken(refreshTokenString string) (*model.TokenDetails, error) {
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := claims["email"].(string)
		role := claims["role"].(string)
		return s.CreateToken(email, role)
	}
	return nil, errors.New("invalid refresh token")
}
