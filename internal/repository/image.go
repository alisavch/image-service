package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/utils"

	"github.com/google/uuid"
)

// ImageRepository provides access to the database.
type ImageRepository struct {
	db *sql.DB
}

// NewImageRepository configures ImageRepository.
func NewImageRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

// FindUserRequestHistory allows to get the history of interaction with the user's service.
func (i *ImageRepository) FindUserRequestHistory(ctx context.Context, id uuid.UUID) ([]models.History, error) {
	query := "SELECT i.uploaded_name, i.resulted_name, r.service_name, r.time_started, r.time_completed, r.status from image_service.request r INNER JOIN image_service.image i on r.image_id = i.id INNER JOIN image_service.user_account ua on ua.id = r.user_account_id where ua.id = $1"
	rows, err := i.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, utils.ErrCreateQuery
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var history []models.History

	for rows.Next() {
		var hist models.History
		if err := rows.Scan(&hist.UploadedName, &hist.ResultedName, &hist.ServiceName, &hist.TimeStarted, &hist.TimeCompleted, &hist.Status); err != nil {
			return history, nil
		}
		history = append(history, hist)
	}

	if err = rows.Err(); err != nil {
		return history, utils.ErrGetHistory
	}
	return history, nil
}

// UploadImage allows to upload an image.
func (i *ImageRepository) UploadImage(ctx context.Context, img models.Image) (uuid.UUID, error) {
	var id uuid.UUID
	query := "INSERT INTO image_service.image(uploaded_name, uploaded_location) VALUES($1, $2) RETURNING id"
	row := i.db.QueryRowContext(ctx, query, img.UploadedName, img.UploadedLocation)
	if err := row.Scan(&id); err != nil {
		return [16]byte{}, utils.ErrUploadImageToDB
	}

	return id, nil
}

// UploadResultedImage allows to upload a resulted image
func (i *ImageRepository) UploadResultedImage(ctx context.Context, img models.Image) error {
	query := "UPDATE image_service.image SET resulted_name = $1, resulted_location = $2 WHERE id = $3"
	result, err := i.db.ExecContext(ctx, query, img.ResultedName, img.ResultedLocation, img.ID)
	if err != nil {
		return utils.ErrUploadImageToDB
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return utils.ErrRowsAffected
	}
	if rows != 1 {
		return fmt.Errorf("%s:%d", utils.ErrExpectedAffected, rows)
	}

	return nil
}

// CreateRequest adds data to multiple tables and returns resulted image id.
func (i *ImageRepository) CreateRequest(ctx context.Context, user models.User, img models.Image, req models.Request) (uuid.UUID, error) {
	var id uuid.UUID
	query := "INSERT INTO image_service.request(user_account_id, image_id, service_name, status, time_started) VALUES($1, $2, $3, $4, $5) RETURNING id"
	row := i.db.QueryRowContext(ctx, query, user.ID, img.ID, req.ServiceName, req.Status, time.Now())
	if err := row.Scan(&id); err != nil {
		return [16]byte{}, utils.ErrCreateRequest
	}

	return id, nil
}

// FindResultedImage finds processed image by ID.WillReturnResult
func (i *ImageRepository) FindResultedImage(ctx context.Context, id uuid.UUID) (models.Image, error) {
	var filename, location string
	image := "SELECT i.resulted_name, i.resulted_location FROM image_service.image i INNER JOIN image_service.request r on i.id = r.image_id WHERE r.id=$1"
	row := i.db.QueryRowContext(ctx, image, id)
	if err := row.Scan(&filename, &location); err != nil {
		return models.Image{}, utils.ErrFindTheResultingImage
	}
	return models.Image{ResultedName: filename, ResultedLocation: location}, nil
}

// FindOriginalImage finds original image by ID.
func (i *ImageRepository) FindOriginalImage(ctx context.Context, id uuid.UUID) (models.Image, error) {
	var filename, location string
	image := "SELECT i.uploaded_name, i.uploaded_location FROM image_service.image i INNER JOIN image_service.request r on i.id = r.image_id WHERE r.id=$1"
	row := i.db.QueryRowContext(ctx, image, id)
	if err := row.Scan(&filename, &location); err != nil {
		return models.Image{}, utils.ErrFindOriginalImage
	}
	return models.Image{UploadedName: filename, UploadedLocation: location}, nil
}

// UpdateStatus updates the status of image processing.
func (i *ImageRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.Status) error {
	updated := "UPDATE image_service.request SET status = $1 WHERE id = $2"
	result, err := i.db.ExecContext(ctx, updated, status, id)
	if err != nil {
		return utils.ErrUpdateStatusRequest
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return utils.ErrRowsAffected
	}
	if rows != 1 {
		return fmt.Errorf("%s:%d", utils.ErrExpectedAffected, rows)
	}
	return nil
}

// SetCompletedTime sets the completion time to the database.
func (i *ImageRepository) SetCompletedTime(ctx context.Context, id uuid.UUID) error {
	updated := "UPDATE image_service.request SET time_completed = $1 WHERE id = $2"
	result, err := i.db.ExecContext(ctx, updated, time.Now(), id)
	if err != nil {
		return utils.ErrSetCompletedTime
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return utils.ErrRowsAffected
	}
	if rows != 1 {
		return fmt.Errorf("%s:%d", utils.ErrExpectedAffected, rows)
	}
	return nil
}
