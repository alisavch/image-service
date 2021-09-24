package repository

import (
	"context"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alisavch/image-service/internal/models"
	"github.com/stretchr/testify/require"
)

// AnyTime if structure for the right time.
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface.
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
		input uuid.UUID
		want  []models.History
		isOk  bool
	}{
		{
			name:  "Test with correct values",
			input: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			mock: func() {
				asString := "00000000-0000-0000-0000-000000000000"
				rows := sqlmock.NewRows([]string{"uploaded_name", "resulted_name", "service", "time_start", "end_of_time", "status"}).AddRow("", "", "", "", "", "")
				mock.ExpectQuery("SELECT (.+) from image_service.request r INNER JOIN image_service.user_image ui on r.user_image_id = ui.id INNER JOIN image_service.uploaded_image upi on ui.uploaded_image_id = upi.id INNER JOIN image_service.resulted_image ri on ri.id = ui.resulting_image_id INNER JOIN image_service.user_account ua on ua.id = ui.user_account_id").
					WithArgs(asString).WillReturnRows(rows)
			},
			want: []models.History(nil),
			isOk: true,
		},
		{
			name: "Test with incorrect values",
			mock: func() {
			},
			input: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
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
		want  uuid.UUID
		isOk  bool
	}{
		{
			name: "Test with correct values",
			input: models.UploadedImage{
				Name:     "filename",
				Location: "location",
			},
			mock: func() {
				asString := "00000000-0000-0000-0000-000000000000"
				rows := sqlmock.NewRows([]string{"id"}).AddRow(asString)
				mock.ExpectQuery("INSERT INTO image_service.uploaded_image").
					WithArgs("filename", "location").WillReturnRows(rows)
			},
			want: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
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
		want  uuid.UUID
		isOk  bool
	}{
		{
			name: "Test with correct values",
			mock: func() {
				mock.ExpectBegin()
				asString := "00000000-0000-0000-0000-000000000000"
				asStringUserID := "00000000-0000-0000-0000-000000000000"
				asStringImageID := "00000000-0000-0000-0000-000000000000"
				resRows := sqlmock.NewRows([]string{"id"}).AddRow(asStringUserID)
				mock.ExpectQuery("INSERT INTO image_service.resulted_image").
					WithArgs("filename", "location", models.Compression).WillReturnRows(resRows)
				uiRows := sqlmock.NewRows([]string{"id"}).AddRow(asStringImageID)
				mock.ExpectQuery("INSERT INTO image_service.user_image").
					WithArgs(asStringUserID, asString, asString, models.Queued).WillReturnRows(uiRows)
				mock.ExpectExec("INSERT INTO image_service.request").
					WithArgs(asStringImageID, AnyTime{}, AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			input: args{
				user: models.User{
					ID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
				},
				uploadedImage: models.UploadedImage{
					ID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
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
			want: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
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
					ID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
				},
				uploadedImage: models.UploadedImage{
					ID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
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
		id      uuid.UUID
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
				id:      [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
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
				id:      [16]byte{},
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
		id uuid.UUID
	}

	type mockBehavior func(args args)

	tests := []struct {
		name  string
		mock  mockBehavior
		input args
		want  models.UploadedImage
		isOk  bool
	}{
		{
			name: "Test with correct values",
			mock: func(args2 args) {
				rows := sqlmock.NewRows([]string{"uploaded_name", "uploaded_location"}).
					AddRow("filename", "location")
				mock.ExpectQuery("SELECT (.+) FROM image_service.uploaded_image").
					WithArgs(args2.id).WillReturnRows(rows)
			},
			input: args{
				id: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			},
			want: models.UploadedImage{
				Name:     "filename",
				Location: "location",
			},
			isOk: true,
		},
		{
			name: "Test with incorrect values",
			mock: func(args2 args) {
				rows := sqlmock.NewRows([]string{"uploaded_name", "uploaded_location"})
				mock.ExpectQuery("SELECT (.+) FROM image_service.uploaded_image").
					WithArgs(args2.id).WillReturnRows(rows)

			},
			input: args{
				id: [16]byte{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

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

func TestImageRepository_UpdateStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected wher opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewImageRepository(db)

	type args struct {
		id     uuid.UUID
		status string
	}

	type mockBehavior func(args args)

	tests := []struct {
		name  string
		mock  mockBehavior
		input args
		want  error
		isOk  bool
	}{
		{
			name: "Test with correct values",
			mock: func(args2 args) {
				mock.ExpectExec("UPDATE image_service.user_image SET status").
					WithArgs(args2.status, args2.id).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				id:     [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
				status: string(models.Processing),
			},
			want: nil,
			isOk: true,
		},
		{
			name: "Test with incorrect values",
			mock: func(args2 args) {
				mock.ExpectExec("UPDATE image_service.user_image SET status").
					WithArgs(args2.status, args2.id).WillReturnError(fmt.Errorf("cannot update image status"))
			},
			input: args{
				id:     [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
				status: string(models.Processing),
			},
			want: fmt.Errorf("cannot update image status"),
			isOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			err := repo.UpdateStatus(context.TODO(), tt.input.id, models.Status(tt.input.status))
			if tt.isOk {
				fmt.Println(err)
				require.NoError(t, err)
				require.Equal(t, tt.want, err)
			} else {
				require.Error(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
