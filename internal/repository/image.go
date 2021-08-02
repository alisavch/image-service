package repository

import (
	"context"
	"database/sql"
	_ "time" // Registers time.

	"github.com/alisavch/image-service/internal/model"
)

// ImageRepository provides access to the database.
type ImageRepository struct {
	db *sql.DB
}

// NewImageRepository is constructor of the ImageRepository.
func NewImageRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

// FindUserHistoryByID allows to get the history of interaction with the user's service.
func (i *ImageRepository) FindUserHistoryByID(ctx context.Context, id int) ([]model.History, error) {
	query := "SELECT upi.uploaded_name, ri.resulted_name, ri.service, r.time_start, r.end_of_time, ui.status from image_service.request r INNER JOIN image_service.user_image ui on r.user_image_id = ui.id INNER JOIN image_service.uploaded_image upi on ui.uploaded_image_id = upi.id INNER JOIN image_service.resulted_image ri on ri.id = ui.resulting_image_id INNER JOIN image_service.user_account ua on ua.id = ui.user_account_id where ua.id = $1"
	rows, err := i.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []model.History

	for rows.Next() {
		var hist model.History
		if err := rows.Scan(&hist.UploadedName, &hist.ResultedName, &hist.Service, &hist.TimeStart, &hist.EndOfTime, &hist.Status); err != nil {
			return history, nil
		}
		history = append(history, hist)
	}

	if err = rows.Err(); err != nil {
		return history, err
	}
	return history, nil
}

// UploadImage allows to upload an image.
func (i *ImageRepository) UploadImage(ctx context.Context, image model.UploadedImage) (int, error) {
	var id int
	query := "INSERT INTO image_service.uploaded_image(uploaded_name, uploaded_location) VALUES($1, $2) RETURNING id"
	row := i.db.QueryRowContext(ctx, query, image.Name, image.Location)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// CreateRequest adds data to multiple tables and returns resulted image id.
func (i *ImageRepository) CreateRequest(ctx context.Context, user model.User, uplImg model.UploadedImage, resImg model.ResultedImage, uI model.UserImage, r model.Request) (int, error) {
	tx, err := i.db.Begin()
	if err != nil {
		return 0, err
	}

	var resultedImageID, userImageID int
	createResultedImage := "INSERT INTO image_service.resulted_image(resulted_name, resulted_location, service) VALUES ($1, $2, $3) RETURNING id"
	row := tx.QueryRowContext(ctx, createResultedImage, resImg.Name, resImg.Location, resImg.Service)
	if err := row.Scan(&resultedImageID); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUserImage := "INSERT INTO image_service.user_image(user_account_id, uploaded_image_id, resulting_image_id, status) VALUES($1, $2, $3, $4) RETURNING id"
	row = tx.QueryRowContext(ctx, createUserImage, user.ID, uplImg.ID, resultedImageID, uI.Status)
	if err := row.Scan(&userImageID); err != nil {
		tx.Rollback()
		return 0, err
	}

	// TODO: request time
	createRequest := "INSERT INTO image_service.request(user_image_id, time_start, end_of_time) VALUES($1, $2, $3)"
	_, err = tx.ExecContext(ctx, createRequest, userImageID, r.TimeStart, r.EndOfTime)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return resultedImageID, tx.Commit()
}

// FindTheResultingImage finds processed image by ID.
func (i *ImageRepository) FindTheResultingImage(ctx context.Context, id int, service model.Service) (model.ResultedImage, error) {
	var filename, location string
	image := "SELECT ri.resulted_name, ri.resulted_location FROM image_service.resulted_image ri WHERE ri.id=$1 and ri.service=$2"
	row := i.db.QueryRowContext(ctx, image, id, service)
	if err := row.Scan(&filename, &location); err != nil {
		return model.ResultedImage{}, err
	}
	return model.ResultedImage{Name: filename, Location: location}, nil
}

// FindOriginalImage finds original image by ID.
func (i *ImageRepository) FindOriginalImage(ctx context.Context, id int) (model.UploadedImage, error) {
	var filename, location string
	image := "SELECT ui.uploaded_name, ui.uploaded_location FROM image_service.uploaded_image ui INNER JOIN image_service.user_image usi on ui.id = usi.uploaded_image_id INNER JOIN image_service.resulted_image ri on ri.id = usi.resulting_image_id WHERE ri.id =$1"
	row := i.db.QueryRowContext(ctx, image, id)
	if err := row.Scan(&filename, &location); err != nil {
		return model.UploadedImage{}, err
	}
	return model.UploadedImage{Name: filename, Location: location}, nil
}
