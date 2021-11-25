package apiserver

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/alisavch/image-service/internal/apiserver/mocks"
	"github.com/alisavch/image-service/internal/broker"
	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/utils"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createImage(t *testing.T, filename string) ([]byte, multipart.File) {
	file, err := os.Create(filename)
	require.NoError(t, err)

	alpha := image.NewAlpha(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			alpha.Set(x, y, color.Alpha{A: uint8(x % 256)})
		}
	}
	err = jpeg.Encode(file, alpha, nil)
	require.NoError(t, err)

	content, err := ioutil.ReadFile(filename)
	require.NoError(t, err)

	return content, file
}

func cleanAfterTest(t *testing.T) {
	err := os.RemoveAll("./uploads/")
	require.NoError(t, err)
	err = os.RemoveAll("./results/")
	require.NoError(t, err)

	d, err := os.Open(".")
	require.NoError(t, err)

	defer func(d *os.File) {
		err := d.Close()
		require.NoError(t, err)
	}(d)

	files, err := d.Readdir(-1)
	require.NoError(t, err)

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".jpeg" || filepath.Ext(file.Name()) == ".Anything" {
				err := os.Remove(file.Name())
				require.NoError(t, err)
			}
		}
	}
}

func TestHandler_findUserHistory(t *testing.T) {
	type fnBehavior func(mockSO *mocks.ServiceOperations, token string, r *http.Request)

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		fn                   fnBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "FindUserHistory without errors",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, r *http.Request) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("FindUserRequestHistory", mock.Anything, s).Return([]models.History{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "[]\n",
		},
		{
			name:        "Cannot complete request to get history",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, r *http.Request) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("FindUserRequestHistory", mock.Anything, s).Return([]models.History{}, fmt.Errorf("cannot complete request to get history"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot complete request to get history\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSO := new(mocks.ServiceOperations)
			mockAWS := new(mocks.S3Bucket)

			currentService := NewAPI(mockSO, mockAWS)
			mq := broker.NewAMQPBroker(mockSO, mockAWS)

			s := NewServer(mq, currentService)

			s.router.HandleFunc("/api/history",
				s.authorize(s.findUserHistory())).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/history", nil)

			tt.fn(mockSO, tt.token, req)

			req.Header.Set(tt.headerName, tt.headerValue)
			s.ServeHTTP(w, req)
			mockSO.AssertExpectations(t)
			require.Equal(t, tt.expectedStatusCode, w.Code)
			require.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_compressImage(t *testing.T) {
	type params struct {
		name     string
		quantity int
	}

	type model struct {
		image models.Image
		req   models.Request
		user  models.User
	}

	uplImg := models.Image{
		ID:               [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		UploadedName:     "filename.jpeg",
		UploadedLocation: "location",
		ResultedName:     "name",
		ResultedLocation: "location",
	}

	reqImg := models.Request{
		ID:            [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		UserAccountID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		ImageID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		ServiceName:   models.Compression,
		Status:        models.Queued,
	}

	userImg := models.User{
		ID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
	}

	modelStruct := model{
		image: uplImg,
		req:   reqImg,
		user:  userImg,
	}

	q := amqp.Queue{Name: "", Messages: 1, Consumers: 1}

	type fnBehavior func(mockSO *mocks.ServiceOperations, mockBucket *mocks.S3Bucket, mockAMQP *mocks.AMQP, token string, model model, storage string)

	tests := []struct {
		name                 string
		headerNames          []string
		headerValues         []string
		params               params
		token                string
		userID               uuid.UUID
		fn                   fnBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "Compress image without errors",
			headerNames:  []string{"Authorization", "Content-Type"},
			headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
			params:       params{name: "width", quantity: 100},
			token:        "token",
			fn: func(mockSO *mocks.ServiceOperations, mockBucket *mocks.S3Bucket, mockAMQP *mocks.AMQP, token string, model model, storage string) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				switch storage {
				case aws:
					mockBucket.On("UploadToS3Bucket", mock.Anything, mock.Anything).Return(mock.Anything, nil)
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(model.image.ID, nil)
					mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.req.ID, nil)
					mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
					mockAMQP.On("QosQueue").Return(nil)
					mockSO.On("UpdateStatus", mock.Anything, model.req.ID, models.Processing).Return(nil)
					mockAMQP.On("Publish", "", q.Name, mock.Anything).Return(nil)
				case local:
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(model.image.ID, nil)
					mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.req.ID, nil)
					mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
					mockAMQP.On("QosQueue").Return(nil)
					mockSO.On("UpdateStatus", mock.Anything, model.req.ID, models.Processing).Return(nil)
					mockAMQP.On("Publish", "", q.Name, mock.Anything).Return(nil)
				}
			},
			expectedStatusCode:   202,
			expectedResponseBody: "{\"Image ID\":\"00000000-0000-0000-0000-000000000000\"}\n",
		},
		{
			name:         "Failed to load file",
			headerNames:  []string{"Authorization", "Content-Type"},
			headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
			params:       params{name: "width", quantity: 100},
			token:        "token",
			fn: func(mockSO *mocks.ServiceOperations, mockBucket *mocks.S3Bucket, mockAMQP *mocks.AMQP, token string, model model, storage string) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				switch storage {
				case aws:
					mockBucket.On("UploadToS3Bucket", mock.Anything, mock.Anything).Return(mock.Anything, nil)
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(uplImg.ID, utils.ErrUpload)

				case local:
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(uplImg.ID, utils.ErrUpload)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot upload the file:cannot upload the file\"}\n",
		},
		{
			name:         "Failed update status",
			headerNames:  []string{"Authorization", "Content-Type"},
			headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
			params:       params{name: "width", quantity: 100},
			token:        "token",
			fn: func(mockSO *mocks.ServiceOperations, mockBucket *mocks.S3Bucket, mockAMQP *mocks.AMQP, token string, model model, storage string) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				switch storage {
				case aws:
					mockBucket.On("UploadToS3Bucket", mock.Anything, mock.Anything).Return(mock.Anything, nil)
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(model.image.ID, nil)
					mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.req.ID, nil)
					mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
					mockAMQP.On("QosQueue").Return(nil)
					mockSO.On("UpdateStatus", mock.Anything, model.req.ID, models.Processing).Return(utils.ErrUpdateStatusRequest)
				case local:
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(model.image.ID, nil)
					mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.req.ID, nil)
					mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
					mockAMQP.On("QosQueue").Return(nil)
					mockSO.On("UpdateStatus", mock.Anything, model.req.ID, models.Processing).Return(utils.ErrUpdateStatusRequest)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot update image status\"}\n",
		},
		{
			name:         "Failed create request",
			headerNames:  []string{"Authorization", "Content-Type"},
			headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
			params:       params{name: "width", quantity: 100},
			token:        "token",
			fn: func(mockSO *mocks.ServiceOperations, mockBucket *mocks.S3Bucket, mockAMQP *mocks.AMQP, token string, model model, storage string) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				switch storage {
				case aws:
					mockBucket.On("UploadToS3Bucket", mock.Anything, mock.Anything).Return(mock.Anything, nil)
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(model.image.ID, nil)
					mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(uplImg.ID, fmt.Errorf("unable to insert resulted image into database"))

				case local:
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(uplImg.ID, nil)
					mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(uplImg.ID, fmt.Errorf("unable to insert resulted image into database"))
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"unable to insert resulted image into database\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := utils.NewConfig()

			mockAMQP := new(mocks.AMQP)
			mockBucket := new(mocks.S3Bucket)
			mockSO := new(mocks.ServiceOperations)

			currentService := NewAPI(mockSO, mockBucket)
			s := NewServer(mockAMQP, currentService)

			s.router.HandleFunc("/api/compress",
				s.authorize(s.compressImage())).Methods(http.MethodPost)

			buf := &bytes.Buffer{}
			writer := multipart.NewWriter(buf)
			header := make(textproto.MIMEHeader)
			header.Set("Content-Disposition", `form-data; name="uploadFile"; filename="filename.jpeg"`)
			header.Set("Content-Type", "image/jpeg")

			content, file := createImage(t, "filename.jpeg")

			part, err := writer.CreatePart(header)
			require.NoError(t, err)
			_, err = io.Copy(part, bytes.NewReader(content))
			require.NoError(t, err)
			err = writer.Close()
			require.NoError(t, err)
			err = file.Close()
			require.NoError(t, err)

			tt.fn(mockSO, mockBucket, mockAMQP, tt.token, modelStruct, conf.Storage)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/compress", buf)

			body, err := ioutil.ReadAll(req.Body)
			require.NoError(t, err)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

			req.Header.Set("Authorization", "Bearer token")
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("Content-Length", "1000")

			q := req.URL.Query()
			q.Add(tt.params.name, strconv.Itoa(tt.params.quantity))
			req.URL.RawQuery = q.Encode()

			err = req.Body.Close()

			require.NoError(t, err)
			s.ServeHTTP(w, req)
			mockSO.AssertExpectations(t)
			mockAMQP.AssertExpectations(t)
			require.Equal(t, tt.expectedStatusCode, w.Code)
			require.Equal(t, tt.expectedResponseBody, w.Body.String())

			cleanAfterTest(t)
		})
	}
}

func TestHandler_findCompressedImage(t *testing.T) {
	type fnBehavior func(mockSO *mocks.ServiceOperations, token string, compressedID uuid.UUID, isOriginal bool)

	resultedImage := models.Image{
		ID:               [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		ResultedName:     "filename",
		ResultedLocation: "Location",
	}

	type params struct {
		name       string
		isOriginal bool
	}

	tests := []struct {
		name                 string
		headerName           []string
		headerValue          []string
		params               params
		token                string
		compressedID         uuid.UUID
		fn                   fnBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "Find compressed image without errors",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			token:        "token",
			compressedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			params:       params{name: "original", isOriginal: false},
			fn: func(mockSO *mocks.ServiceOperations, token string, compressedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("CheckStatus", mock.Anything, compressedID).Return(nil)
				mockSO.On("FindResultedImage", mock.Anything, compressedID).Return(resultedImage, nil)
				mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.SavedImage{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name:         "Find original image without errors",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			token:        "token",
			compressedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			params:       params{name: "original", isOriginal: true},
			fn: func(mockSO *mocks.ServiceOperations, token string, compressedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				if isOriginal {
					mockSO.On("FindOriginalImage", mock.Anything, compressedID).Return(models.Image{}, nil)
					mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.SavedImage{}, nil)
				}
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name:         "Wrong to find image",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			token:        "token",
			compressedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			fn: func(mockSO *mocks.ServiceOperations, token string, compressedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("CheckStatus", mock.Anything, compressedID).Return(nil)
				mockSO.On("FindResultedImage", mock.Anything, compressedID).Return(models.Image{}, utils.ErrFindImage)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot find image:cannot find image\"}\n",
		},
		{
			name:         "Incorrectly saved image",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			token:        "token",
			params:       params{name: "original", isOriginal: false},
			compressedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			fn: func(mockSO *mocks.ServiceOperations, token string, compressedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("CheckStatus", mock.Anything, compressedID).Return(nil)
				mockSO.On("FindResultedImage", mock.Anything, compressedID).Return(resultedImage, nil)
				mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.SavedImage{}, utils.ErrSaveImage)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot save image:cannot save image\"}\n",
		},
		{
			name:         "Wrong to find original image",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			token:        "token",
			params:       params{name: "original", isOriginal: true},
			compressedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			fn: func(mockSO *mocks.ServiceOperations, token string, compressedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				if isOriginal {
					mockSO.On("FindOriginalImage", mock.Anything, compressedID).Return(models.Image{}, utils.ErrFindImage)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot find image:cannot find image\"}\n",
		},
		{
			name:         "Incorrectly saved original image",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			token:        "token",
			params:       params{name: "original", isOriginal: true},
			compressedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			fn: func(mockSO *mocks.ServiceOperations, token string, compressedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				if isOriginal {
					mockSO.On("FindOriginalImage", mock.Anything, compressedID).Return(models.Image{ID: s, UploadedName: "filename", UploadedLocation: "location"}, nil)
					mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.SavedImage{}, utils.ErrSaveImage)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot save image:cannot save image\"}\n",
		},
	}

	getCompressedURL := "/api/compress/%s"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBucket := new(mocks.S3Bucket)
			mockSO := new(mocks.ServiceOperations)

			currentService := NewAPI(mockSO, mockBucket)
			mq := broker.NewAMQPBroker(mockSO, mockBucket)

			s := NewServer(mq, currentService)

			tt.fn(mockSO, tt.token, tt.compressedID, tt.params.isOriginal)

			s.router.HandleFunc(fmt.Sprintf(getCompressedURL, "{compressedID}"),
				s.authorize(s.findCompressedImage())).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf(getCompressedURL, tt.compressedID), nil)

			q := req.URL.Query()
			q.Add(tt.params.name, strconv.FormatBool(tt.params.isOriginal))
			req.URL.RawQuery = q.Encode()

			req.Header.Set(tt.headerName[0], tt.headerValue[0])

			s.ServeHTTP(w, req)
			mockSO.AssertExpectations(t)
			require.Equal(t, tt.expectedStatusCode, w.Code)
			require.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_convertImage(t *testing.T) {
	type model struct {
		image models.Image
		req   models.Request
		user  models.User
	}

	uplImg := models.Image{
		ID:               [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		UploadedName:     "filename.jpeg",
		UploadedLocation: "location",
		ResultedName:     "name",
		ResultedLocation: "location",
	}

	reqImg := models.Request{
		ID:            [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		UserAccountID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		ImageID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		ServiceName:   models.Conversion,
		Status:        models.Queued,
	}

	userImg := models.User{
		ID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
	}

	modelStruct := model{
		image: uplImg,
		req:   reqImg,
		user:  userImg,
	}

	q := amqp.Queue{Name: "", Messages: 1, Consumers: 1}

	type fnBehavior func(mockSO *mocks.ServiceOperations, mockBucket *mocks.S3Bucket, mockAMQP *mocks.AMQP, token string, model model, storage string)

	tests := []struct {
		name                 string
		headerNames          []string
		headerValues         []string
		inputImage           models.Image
		contentType          string
		token                string
		convertedID          string
		userID               uuid.UUID
		fn                   fnBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "Convert image without errors",
			headerNames:  []string{"Authorization", "Content-Type"},
			headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
			token:        "token",
			fn: func(mockSO *mocks.ServiceOperations, mockBucket *mocks.S3Bucket, mockAMQP *mocks.AMQP, token string, model model, storage string) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				switch storage {
				case aws:
					mockBucket.On("UploadToS3Bucket", mock.Anything, mock.Anything).Return(mock.Anything, nil)
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(model.image.ID, nil)
					mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.req.ID, nil)
					mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
					mockAMQP.On("QosQueue").Return(nil)
					mockSO.On("UpdateStatus", mock.Anything, model.req.ID, models.Processing).Return(nil)
					mockAMQP.On("Publish", "", q.Name, mock.Anything).Return(nil)
				case local:
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(model.image.ID, nil)
					mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.req.ID, nil)
					mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
					mockAMQP.On("QosQueue").Return(nil)
					mockSO.On("UpdateStatus", mock.Anything, model.req.ID, models.Processing).Return(nil)
					mockAMQP.On("Publish", "", q.Name, mock.Anything).Return(nil)
				}
			},
			expectedStatusCode:   202,
			expectedResponseBody: "{\"Image ID\":\"00000000-0000-0000-0000-000000000000\"}\n",
		},
		{
			name:         "Failed upload file",
			headerNames:  []string{"Authorization", "Content-Type"},
			headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
			token:        "token",
			fn: func(mockSO *mocks.ServiceOperations, mockBucket *mocks.S3Bucket, mockAMQP *mocks.AMQP, token string, model model, storage string) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				switch storage {
				case aws:
					mockBucket.On("UploadToS3Bucket", mock.Anything, mock.Anything).Return(mock.Anything, nil)
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(uplImg.ID, utils.ErrUploadImageToDB)
				case local:
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(uplImg.ID, utils.ErrUploadImageToDB)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot upload the file:unable to insert image into database\"}\n",
		},
		{
			name:         "Failed update status",
			headerNames:  []string{"Authorization", "Content-Type"},
			headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
			token:        "token",
			fn: func(mockSO *mocks.ServiceOperations, mockBucket *mocks.S3Bucket, mockAMQP *mocks.AMQP, token string, model model, storage string) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				switch storage {
				case aws:
					mockBucket.On("UploadToS3Bucket", mock.Anything, mock.Anything).Return(mock.Anything, nil)
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(model.image.ID, nil)
					mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.req.ID, nil)
					mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
					mockAMQP.On("QosQueue").Return(nil)
					mockSO.On("UpdateStatus", mock.Anything, uplImg.ID, models.Processing).Return(utils.ErrUpdateStatusRequest)
				case local:
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(uplImg.ID, nil)
					mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.req.ID, nil)
					mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
					mockAMQP.On("QosQueue").Return(nil)
					mockSO.On("UpdateStatus", mock.Anything, model.req.ID, models.Processing).Return(utils.ErrUpdateStatusRequest)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot update image status\"}\n",
		},
		{
			name:         "Failed create request",
			headerNames:  []string{"Authorization", "Content-Type"},
			headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
			token:        "token",
			fn: func(mockSO *mocks.ServiceOperations, mockBucket *mocks.S3Bucket, mockAMQP *mocks.AMQP, token string, model model, storage string) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				switch storage {
				case aws:
					mockBucket.On("UploadToS3Bucket", mock.Anything, mock.Anything).Return(mock.Anything, nil)
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(model.image.ID, nil)
					mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(s, utils.ErrCreateRequest)
				case local:
					mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(model.image.ID, nil)
					mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(s, utils.ErrCreateRequest)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot create request\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := utils.NewConfig()

			mockAMQP := new(mocks.AMQP)
			mockBucket := new(mocks.S3Bucket)
			mockSO := new(mocks.ServiceOperations)

			currentService := NewAPI(mockSO, mockBucket)
			s := NewServer(mockAMQP, currentService)

			s.router.HandleFunc("/api/convert",
				s.authorize(s.convertImage())).Methods(http.MethodPost)

			content, file := createImage(t, "filename.jpeg")

			buf := &bytes.Buffer{}
			writer := multipart.NewWriter(buf)
			header := make(textproto.MIMEHeader)
			header.Set("Content-Disposition", `form-data; name="uploadFile"; filename="filename.jpeg"`)
			header.Set("Content-Type", "image/jpeg")
			part, err := writer.CreatePart(header)
			require.NoError(t, err)
			_, err = io.Copy(part, bytes.NewReader(content))
			require.NoError(t, err)
			err = writer.Close()
			require.NoError(t, err)
			err = file.Close()
			require.NoError(t, err)

			tt.fn(mockSO, mockBucket, mockAMQP, tt.token, modelStruct, conf.Storage)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/convert", buf)

			body, err := ioutil.ReadAll(req.Body)
			require.NoError(t, err)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

			req.Header.Set("Authorization", "Bearer token")
			req.Header.Set("Content-Type", writer.FormDataContentType())

			err = req.Body.Close()
			require.NoError(t, err)

			s.ServeHTTP(w, req)
			mockSO.AssertExpectations(t)
			require.Equal(t, tt.expectedStatusCode, w.Code)
			require.Equal(t, tt.expectedResponseBody, w.Body.String())

			cleanAfterTest(t)
		})
	}
}

func TestHandler_findConvertedImage(t *testing.T) {
	type fnBehavior func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool)

	resultedImage := models.Image{
		ID:               [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		ResultedName:     "filename",
		ResultedLocation: "location",
	}

	type params struct {
		name       string
		isOriginal bool
	}

	tests := []struct {
		name                 string
		headerName           []string
		headerValue          []string
		convertedID          uuid.UUID
		params               params
		token                string
		userID               uuid.UUID
		fn                   fnBehavior
		isRemoteStorage      bool
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Find converted image without errors",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			convertedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("CheckStatus", mock.Anything, convertedID).Return(nil)
				mockSO.On("FindResultedImage", mock.Anything, convertedID).Return(resultedImage, nil)
				mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.SavedImage{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name:        "Find original image without errors",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			params:      params{name: "original", isOriginal: true},
			convertedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				if isOriginal {
					mockSO.On("FindOriginalImage", mock.Anything, convertedID).Return(models.Image{}, nil)
					mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.SavedImage{}, nil)
				}
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name:        "Test with invalid token",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			convertedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, utils.ErrEmptyToken)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"token is empty\"}\n",
		},
		{
			name:        "Wrong to find image",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			convertedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("CheckStatus", mock.Anything, convertedID).Return(nil)
				mockSO.On("FindResultedImage", mock.Anything, mock.Anything).Return(models.Image{}, utils.ErrFindImage)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot find image:cannot find image\"}\n",
		},
		{
			name:        "Incorrectly saved image",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			convertedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("CheckStatus", mock.Anything, convertedID).Return(nil)
				mockSO.On("FindResultedImage", mock.Anything, convertedID).Return(resultedImage, nil)
				mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.SavedImage{}, utils.ErrSaveImage)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot save image:cannot save image\"}\n",
		},
		{
			name:        "Wrong to find original image",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			params:      params{name: "original", isOriginal: true},
			convertedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				if isOriginal {
					mockSO.On("FindOriginalImage", mock.Anything, convertedID).Return(models.Image{}, utils.ErrFindImage)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot find image:cannot find image\"}\n",
		},
		{
			name:        "Incorrectly saved original image",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			params:      params{name: "original", isOriginal: true},
			convertedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				if isOriginal {
					mockSO.On("FindOriginalImage", mock.Anything, convertedID).Return(models.Image{}, nil)
					mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.SavedImage{}, utils.ErrSaveImage)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot save image:cannot save image\"}\n",
		},
	}

	getConvertedURL := "/api/convert/%s"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBucket := new(mocks.S3Bucket)
			mockSO := new(mocks.ServiceOperations)

			currentService := NewAPI(mockSO, mockBucket)
			mq := broker.NewAMQPBroker(mockSO, mockBucket)

			s := NewServer(mq, currentService)

			tt.fn(mockSO, tt.token, tt.convertedID, tt.params.isOriginal)

			s.router.HandleFunc(fmt.Sprintf(getConvertedURL, tt.convertedID),
				s.authorize(s.findConvertedImage())).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf(getConvertedURL, tt.convertedID), nil)

			q := req.URL.Query()
			q.Add(tt.params.name, strconv.FormatBool(tt.params.isOriginal))
			req.URL.RawQuery = q.Encode()

			req.Header.Set(tt.headerName[0], tt.headerValue[0])

			s.ServeHTTP(w, req)
			mockSO.AssertExpectations(t)
			require.Equal(t, tt.expectedStatusCode, w.Code)
			require.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
