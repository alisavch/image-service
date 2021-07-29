package service

import (
	"time"

	"github.com/alisavch/image-service/internal/utils"

	"github.com/alisavch/image-service/internal/model"
	"github.com/alisavch/image-service/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	signingKey = "QTicXLOhp5uxp80xTGrosP5Hpa9C"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID int64 `json:"user_id"`
}

// AuthService provides access to repository.
type AuthService struct {
	repo repository.Authorization
}

// NewAuthService is constructor of AuthService.
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// CreateUser creates user.
func (s *AuthService) CreateUser(user model.User) (int, error) {
	user.Password, _ = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

// GenerateToken generates token.
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return "get user error", err
	}
	match := s.checkPasswordHash(password, user.Password)
	if !match {
		return "err check passwords", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		int64(user.ID),
	})
	return token.SignedString([]byte(signingKey))
}

// ParseToken parses token.
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, utils.ErrSigningMethod
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, nil
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, utils.ErrInvalidToken
	}
	return int(claims.UserID), nil
}

func (s *AuthService) generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *AuthService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
