package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alisavch/image-service/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAuthRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected wher opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewAuthRepository(db)

	tests := []struct {
		name  string
		mock  func()
		input models.User
		want  int
		isOk  bool
	}{
		{
			name: "Test with correct values",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)
				mock.ExpectQuery("INSERT INTO image_service.user_account").
					WithArgs("mock", "12345").WillReturnRows(rows)
			},
			input: models.User{
				Username: "mock",
				Password: "12345",
			},
			want: 1,
			isOk: true,
		},
		{
			name: "Test with incorrect values",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO image_service.user_account").
					WithArgs("mock", "").WillReturnRows(rows)
			},
			input: models.User{
				Username: "mock",
				Password: "",
			},
			isOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.CreateUser(context.TODO(), tt.input)
			if tt.isOk {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.Error(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthRepository_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected wher opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewAuthRepository(db)

	tests := []struct {
		name  string
		mock  func()
		input string
		want  models.User
		isOk  bool
	}{
		{
			name: "Test with correct values",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "password"}).
					AddRow(1, "12345")
				mock.ExpectQuery("SELECT (.+) FROM image_service.user_account").
					WithArgs("mock").WillReturnRows(rows)
			},
			input: "mock",
			want: models.User{
				ID:       1,
				Password: "12345",
			},
			isOk: true,
		},
		{
			name: "Test with incorrect values",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "password"})
				mock.ExpectQuery("SELECT (.+) FROM image_service.user_account").
					WithArgs("not_found").WillReturnRows(rows)
			},
			input: "not_found",
			isOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.GetUser(context.TODO(), tt.input)
			if tt.isOk {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.Error(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
