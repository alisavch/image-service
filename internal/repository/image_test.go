package repository

import (
	"context"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alisavch/image-service/internal/models"
	"github.com/stretchr/testify/require"
)

// AnyTime ...
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestImageRepository_FindUserHistoryByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected wher opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewImageRepository(db)

	tests := []struct {
		name  string
		mock  func()
		input int
		want  []models.History
		isOk  bool
	}{
		{
			name:  "Test with correct values",
			input: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"uploaded_name", "resulted_name", "service", "time_start", "end_of_time", "status"}).AddRow("", "", "", "", "", "")
				mock.ExpectQuery("SELECT (.+) from image_service.request r INNER JOIN image_service.user_image ui on r.user_image_id = ui.id INNER JOIN image_service.uploaded_image upi on ui.uploaded_image_id = upi.id INNER JOIN image_service.resulted_image ri on ri.id = ui.resulting_image_id INNER JOIN image_service.user_account ua on ua.id = ui.user_account_id").
					WithArgs(1).WillReturnRows(rows)
			},
			want: []models.History(nil),
			isOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.FindUserHistoryByID(context.TODO(), tt.input)
			if tt.isOk {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			} else {
				require.Error(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestImageRepository_UploadImage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected wher opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewImageRepository(db)

	tests := []struct {
		name  string
		mock  func()
		input models.UploadedImage
		want  int
		isOk  bool
	}{
		{
			name: "Test with correct values",
			input: models.UploadedImage{
				Name:     "filename",
				Location: "location",
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO image_service.uploaded_image").
					WithArgs("filename", "location").WillReturnRows(rows)
			},
			want: 1,
			isOk: true,
		},
		{
			name: "Test with incorrect values",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO image_service.uploaded_image").
					WithArgs("", "location").WillReturnRows(rows)
			},
			input: models.UploadedImage{
				Name:     "",
				Location: "location",
			},
			isOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.UploadImage(context.TODO(), tt.input)
			if tt.isOk {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			} else {
				require.Error(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestImageRepository_CreateRequest(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected wher opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewImageRepository(db)

	type args struct {
		user          models.User
		uploadedImage models.UploadedImage
		resultedImage models.ResultedImage
		userImage     models.UserImage
		request       models.Request
	}

	tests := []struct {
		name  string
		mock  func()
		input args
		want  int
		isOk  bool
	}{
		{
			name: "Test with correct values",
			mock: func() {
				mock.ExpectBegin()

				resRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO image_service.resulted_image").
					WithArgs("filename", "location", models.Compression).WillReturnRows(resRows)
				uiRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO image_service.user_image").
					WithArgs(1, 1, 1, models.Queued).WillReturnRows(uiRows)
				mock.ExpectExec("INSERT INTO image_service.request").
					WithArgs(1, AnyTime{}, AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			input: args{
				user: models.User{
					ID: 1,
				},
				uploadedImage: models.UploadedImage{
					ID: 1,
				},
				resultedImage: models.ResultedImage{
					Name:     "filename",
					Location: "location",
					Service:  models.Compression,
				},
				userImage: models.UserImage{
					Status: models.Queued,
				},
				request: models.Request{
					TimeStart: time.Now(),
					EndOfTime: time.Now(),
				},
			},
			want: 1,
			isOk: true,
		},
		{
			name: "Test with incorrect values",
			mock: func() {
				mock.ExpectBegin()
				resRows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO image_service.resulted_image").
					WithArgs("", "", models.Compression).WillReturnRows(resRows)
				mock.ExpectRollback()
			},
			input: args{
				user: models.User{
					ID: 1,
				},
				uploadedImage: models.UploadedImage{
					ID: 1,
				},
				resultedImage: models.ResultedImage{
					Name:     "",
					Location: "",
					Service:  models.Compression,
				},
				userImage: models.UserImage{
					Status: models.Queued,
				},
				request: models.Request{
					TimeStart: time.Now(),
					EndOfTime: time.Now(),
				},
			},
			isOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.CreateRequest(context.TODO(), tt.input.user, tt.input.uploadedImage, tt.input.resultedImage, tt.input.userImage, tt.input.request)
			if tt.isOk {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			} else {
				require.Error(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestImageRepository_FindTheResultingImage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected wher opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewImageRepository(db)

	type args struct {
		id      int
		service models.Service
	}

	type mockBehavior func(args args)

	tests := []struct {
		name  string
		mock  mockBehavior
		input args
		want  models.ResultedImage
		isOk  bool
	}{
		{
			name: "Test with correct values",
			input: args{
				id:      1,
				service: models.Conversion,
			},
			mock: func(args args) {
				rows := sqlmock.NewRows([]string{"resulted_name", "resulted_location"}).
					AddRow("filename", "location")
				mock.ExpectQuery("SELECT (.+) FROM image_service.resulted_image").
					WithArgs(args.id, args.service).WillReturnRows(rows)
			},
			want: models.ResultedImage{
				Name:     "filename",
				Location: "location",
			},
			isOk: true,
		},
		{
			name: "Test with incorrect values",
			mock: func(args args) {
				rows := sqlmock.NewRows([]string{"resulted_name", "resulted_location"})
				mock.ExpectQuery("SELECT (.+) FROM image_service.resulted_image").
					WithArgs(args.id, args.service).WillReturnRows(rows)
			},
			input: args{
				id:      0,
				service: models.Compression,
			},
			isOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := repo.FindTheResultingImage(context.TODO(), tt.input.id, tt.input.service)
			if tt.isOk {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			} else {
				require.Error(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestImageRepository_FindOriginalImage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected wher opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewImageRepository(db)

	type args struct {
		id int
	}

	tests := []struct {
		name  string
		mock  func()
		input args
		want  models.UploadedImage
		isOk  bool
	}{
		{
			name: "Test with correct values",
			mock: func() {
				rows := sqlmock.NewRows([]string{"uploaded_name", "uploaded_location"}).
					AddRow("filename", "location")
				mock.ExpectQuery("SELECT (.+) FROM image_service.uploaded_image").
					WithArgs(1).WillReturnRows(rows)
			},
			input: args{
				id: 1,
			},
			want: models.UploadedImage{
				Name:     "filename",
				Location: "location",
			},
			isOk: true,
		},
		{
			name: "Test with incorrect values",
			mock: func() {
				rows := sqlmock.NewRows([]string{"uploaded_name", "uploaded_location"})
				mock.ExpectQuery("SELECT (.+) FROM image_service.uploaded_image").
					WithArgs(0).WillReturnRows(rows)

			},
			input: args{
				id: 0,
			},
			isOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.FindOriginalImage(context.TODO(), tt.input.id)
			if tt.isOk {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			} else {
				require.Error(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
