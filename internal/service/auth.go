package service

import (
	"context"
	"fmt"
	"time"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	signingKey = "QTicXLOhp5uxp80xTGrosP5Hpa9C"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"user_id"`
}

// AuthService provides access to repository.
type AuthService struct {
	repo AuthorizationRepo
}

// NewAuthService configures AuthService.
func NewAuthService(repo AuthorizationRepo) *AuthService {
	return &AuthService{repo: repo}
}

// CreateUser creates user.
func (s *AuthService) CreateUser(ctx context.Context, user models.User) (id uuid.UUID, err error) {
	user.Password, err = generatePasswordHash(user.Password)
	if err != nil {
		return [16]byte{}, fmt.Errorf("%s:%s", utils.ErrGenerateHash, err)
	}
	return s.repo.CreateUser(ctx, user)
}

// GenerateToken generates token.
func (s *AuthService) GenerateToken(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		return "get user error", err
	}

	match := checkPasswordHash(password, user.Password)
	if !match {
		return "password verification error", err
	}

	conf := utils.NewConfig()
	jwtTTL, err := time.ParseDuration(conf.TokenTTL)
	if err != nil {
		return "error getting jwt ttl", nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})
	return token.SignedString([]byte(signingKey))
}

// ParseToken parses token.
func (s *AuthService) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, utils.ErrSigningMethod
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return [16]byte{}, nil
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return [16]byte{}, utils.ErrInvalidToken
	}
	return claims.UserID, nil
}

func generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
