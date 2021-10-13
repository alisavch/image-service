package apiserver

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strconv"
	"testing"

	"github.com/alisavch/image-service/internal/broker"

	"github.com/alisavch/image-service/internal/apiserver/mocks"

	"github.com/alisavch/image-service/internal/utils"

	"github.com/alisavch/image-service/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_findUserHistory(t *testing.T) {
	type fnBehavior func(mockSO *mocks.ServiceOperations, token string, r *http.Request)

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		userID               uuid.UUID
		fn                   fnBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "FindUserHistoryByID without errors",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			userID:      [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			fn: func(mockSO *mocks.ServiceOperations, token string, r *http.Request) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("FindUserHistoryByID", mock.Anything, s).Return([]models.History{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "[]\n",
		},
		{
			name:        "Users IDs do not match",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			userID:      [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			fn: func(mockSO *mocks.ServiceOperations, token string, r *http.Request) {
				asString := "00000000-0011-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"users IDs do not match\"}\n",
		},
		{
			name:        "Cannot complete request to get history",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			userID:      [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			fn: func(mockSO *mocks.ServiceOperations, token string, r *http.Request) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("FindUserHistoryByID", mock.Anything, s).Return([]models.History{}, fmt.Errorf("cannot complete request to get history"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot complete request to get history\"}\n",
		},
	}

	historyURL := "/api/user/%s/history"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSO := new(mocks.ServiceOperations)

			service := &Service{ServiceOperations: mockSO}
			mq := broker.NewAMQPBroker()

			s := NewServer(mq, service)

			s.router.HandleFunc(fmt.Sprintf(historyURL, "{userID}"),
				s.authorize(s.findUserHistory())).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf(historyURL, tt.userID), nil)

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

	uplImg := models.UploadedImage{
		ID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		Name:     "filename.jpeg",
		Location: "location",
	}

	type fnBehavior func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage, isRemoteStorage bool)

	tests := []struct {
		name                 string
		headerNames          []string
		headerValues         []string
		inputImage           models.UploadedImage
		params               params
		token                string
		userID               uuid.UUID
		isRemoteStorage      bool
		fn                   fnBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		// {
		//	name:            "Compress image without errors",
		//	headerNames:     []string{"Authorization", "Content-Type"},
		//	headerValues:    []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
		//	params:          params{name: "width", quantity: 100},
		//	token:           "token",
		//	userID:          [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		//	isRemoteStorage: true,
		//	fn: func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage, isRemoteStorage bool) {
		//		asString := "00000000-0000-0000-0000-000000000000"
		//		s := uuid.MustParse(asString)
		//		mockSO.On("ParseToken", token).Return(s, nil)
		//		if isRemoteStorage {
		//			mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(s, nil)
		//		}
		//		mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
		//		mockAMQP.On("QosQueue").Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
		//		mockSO.On("UpdateStatus", mock.Anything, 1, models.Processing).Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
		//		mockSO.On("CompressImage", 100, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(models.ResultedImage{ID: s, Name: "name", Location: "location"}, nil)
		//		mockSO.On("UpdateStatus", mock.Anything, 1, models.Done).Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Done)).Return(nil).Return(nil)
		//		mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
		//	},
		//	expectedStatusCode:   500,
		//	expectedResponseBody: "{\"Image ID\":1}\n",
		// },
		{
			name:         "Users IDs do not match",
			headerNames:  []string{"Authorization"},
			headerValues: []string{"Bearer token"},
			token:        "token",
			userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			fn: func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage, isRemoteStorage bool) {
				asString := "00000000-0000-0000-0000-000000000001"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"users IDs do not match\"}\n",
		},
		{
			name:         "Failed to load file",
			headerNames:  []string{"Authorization", "Content-Type"},
			headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
			params:       params{name: "width", quantity: 100},
			token:        "token",
			userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			fn: func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage, isRemoteStorage bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(s, utils.ErrUpload)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot upload the file\"}\n",
		},
		// {
		//	name:         "Failed update status",
		//	headerNames:  []string{"Authorization", "Content-Type"},
		//	headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
		//	params:       params{name: "width", quantity: 100},
		//	token:        "token",
		//	userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		//	fn: func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage, isRemoteStorage bool) {
		//		asString := "00000000-0000-0000-0000-000000000000"
		//		s := uuid.MustParse(asString)
		//		mockSO.On("ParseToken", token).Return(s, nil)
		//		mockSO.On("UploadImage", mock.Anything, uplImg).Return(s, nil)
		//		mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
		//		mockAMQP.On("QosQueue").Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil)
		//		mockSO.On("UpdateStatus", mock.Anything, s, models.Processing).Return(fmt.Errorf("cannot update image status"))
		//	},
		//	expectedStatusCode:   500,
		//	expectedResponseBody: "{\"error\":\"cannot update image status\"}\n",
		// },
		// {
		//	name:         "Failed compress image",
		//	headerNames:  []string{"Authorization", "Content-Type"},
		//	headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
		//	params:       params{name: "width", quantity: 100},
		//	token:        "token",
		//	userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		//	fn: func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage, isRemoteStorage bool) {
		//		asString := "00000000-0000-0000-0000-000000000000"
		//		s := uuid.MustParse(asString)
		//		mockSO.On("ParseToken", token).Return(s, nil)
		//		mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(s, nil)
		//		mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
		//		mockAMQP.On("QosQueue").Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
		//		mockSO.On("UpdateStatus", mock.Anything, 1, models.Processing).Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
		//		mockSO.On("CompressImage", 100, mock.Anything).Return(models.ResultedImage{}, fmt.Errorf("cannot open image"))
		//	},
		//	expectedStatusCode:   500,
		//	expectedResponseBody: "{\"error\":\"cannot open image\"}\n",
		// },
		// {
		//	name:         "Failed convert image 2",
		//	headerNames:  []string{"Authorization", "Content-Type"},
		//	headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
		//	token:        "token",
		//	userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		//	fn: func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage, isRemoteStorage bool) {
		//		asString := "00000000-0000-0000-0000-000000000000"
		//		s := uuid.MustParse(asString)
		//		mockSO.On("ParseToken", token).Return(s, nil)
		//		mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(s, nil)
		//		mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
		//		mockAMQP.On("QosQueue").Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
		//		mockSO.On("UpdateStatus", mock.Anything, 1, models.Processing).Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
		//		mockSO.On("CompressImage", 150, mock.Anything).Return(models.ResultedImage{ID: s, Name: "name", Location: "location"}, nil)
		//		mockSO.On("UpdateStatus", mock.Anything, 1, models.Done).Return(fmt.Errorf("cannot update image status"))
		//	},
		//	expectedStatusCode:   500,
		//	expectedResponseBody: "{\"error\":\"cannot update image status\"}\n",
		// },
		// {
		//	name:         "Failed create request",
		//	headerNames:  []string{"Authorization", "Content-Type"},
		//	headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
		//	params:       params{name: "width", quantity: 100},
		//	token:        "token",
		//	userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		//	fn: func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage, isRemoteStorage bool) {
		//		asString := "00000000-0000-0000-0000-000000000000"
		//		s := uuid.MustParse(asString)
		//		mockSO.On("ParseToken", token).Return(s, nil)
		//		mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(s, nil)
		//		mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
		//		mockAMQP.On("QosQueue").Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
		//		mockSO.On("UpdateStatus", mock.Anything, 1, models.Processing).Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
		//		mockSO.On("CompressImage", 100, mock.Anything).Return(models.ResultedImage{ID: s, Name: "name", Location: "location"}, nil)
		//		mockSO.On("UpdateStatus", mock.Anything, 1, models.Done).Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Done)).Return(nil).Return(nil)
		//		mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, fmt.Errorf("unable to insert resulted image into database"))
		//	},
		//	expectedStatusCode:   500,
		//	expectedResponseBody: "{\"error\":\"unable to insert resulted image into database\"}\n",
		// },
	}

	compressURL := "/api/user/%s/compress"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSO := new(mocks.ServiceOperations)
			mockAMQP := new(mocks.AMQP)

			tt.fn(mockSO, mockAMQP, tt.token, uplImg, tt.isRemoteStorage)

			service := &Service{ServiceOperations: mockSO}
			mq := broker.NewAMQPBroker()

			s := NewServer(mq, service)

			s.router.HandleFunc(fmt.Sprintf(compressURL, "{userID}"),
				s.authorize(s.compressImage())).Methods(http.MethodPost)

			imageBytes := []byte("uploadFile: undefined")
			buf := &bytes.Buffer{}
			writer := multipart.NewWriter(buf)
			header := make(textproto.MIMEHeader)
			header.Set("Content-Disposition", `form-data; name="uploadFile"; filename="filename.jpeg"`)
			header.Set("Content-Type", "image/jpeg")
			part, err := writer.CreatePart(header)
			require.NoError(t, err)
			_, err = part.Write(imageBytes)
			require.NoError(t, err)
			err = writer.Close()
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf(compressURL, tt.userID), buf)

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
		})
	}
}

func TestHandler_findCompressedImage(t *testing.T) {
	type fnBehavior func(mockSO *mocks.ServiceOperations, token string, compressedID uuid.UUID, isOriginal bool)

	resultedImage := models.ResultedImage{
		ID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		Name:     "filename",
		Location: "Location",
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
		userID               uuid.UUID
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
			userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			params:       params{name: "original", isOriginal: false},
			fn: func(mockSO *mocks.ServiceOperations, token string, compressedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("FindTheResultingImage", mock.Anything, compressedID, models.Compression).Return(resultedImage, nil)
				mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.Image{}, nil)
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
			userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			params:       params{name: "original", isOriginal: true},
			fn: func(mockSO *mocks.ServiceOperations, token string, compressedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				if isOriginal {
					mockSO.On("FindOriginalImage", mock.Anything, compressedID).Return(models.UploadedImage{}, nil)
					mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.Image{}, nil)
				}
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name:         "Inequality of identifiers",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			token:        "token",
			compressedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			fn: func(mockSO *mocks.ServiceOperations, token string, compressedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000001"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"users IDs do not match\"}\n",
		},
		{
			name:         "Wrong to find image",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			token:        "token",
			compressedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			fn: func(mockSO *mocks.ServiceOperations, token string, compressedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("FindTheResultingImage", mock.Anything, compressedID, models.Compression).Return(models.ResultedImage{}, utils.ErrFindImage)
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
			userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			fn: func(mockSO *mocks.ServiceOperations, token string, compressedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("FindTheResultingImage", mock.Anything, compressedID, models.Compression).Return(resultedImage, nil)
				mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.Image{}, utils.ErrSaveImage)
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
			userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},

			fn: func(mockSO *mocks.ServiceOperations, token string, compressedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				if isOriginal {
					mockSO.On("FindOriginalImage", mock.Anything, compressedID).Return(models.UploadedImage{}, utils.ErrFindImage)
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
			userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			fn: func(mockSO *mocks.ServiceOperations, token string, compressedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				if isOriginal {
					mockSO.On("FindOriginalImage", mock.Anything, compressedID).Return(models.UploadedImage{ID: s, Name: "filename", Location: "location"}, nil)
					mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.Image{}, utils.ErrSaveImage)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot save image:cannot save image\"}\n",
		},
	}

	getCompressedURL := "/api/user/%s/compress/%s"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSO := new(mocks.ServiceOperations)

			tt.fn(mockSO, tt.token, tt.compressedID, tt.params.isOriginal)

			service := &Service{ServiceOperations: mockSO}
			mq := broker.NewAMQPBroker()

			s := NewServer(mq, service)

			s.router.HandleFunc(fmt.Sprintf(getCompressedURL, "{userID}", "{compressedID}"),
				s.authorize(s.findCompressedImage())).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf(getCompressedURL, tt.userID, tt.compressedID), nil)

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
	type fnBehavior func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage)

	uplImg := models.UploadedImage{
		ID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		Name:     "filename.jpeg",
		Location: "location",
	}

	tests := []struct {
		name                 string
		headerNames          []string
		headerValues         []string
		inputImage           models.UploadedImage
		contentType          string
		token                string
		convertedID          string
		userID               uuid.UUID
		fn                   fnBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		// {
		//	name:         "Convert image without errors",
		//	headerNames:  []string{"Authorization", "Content-Type"},
		//	headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
		//	token:        "token",
		//	userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		//	fn: func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage) {
		//		asString := "00000000-0000-0000-0000-000000000000"
		//		s := uuid.MustParse(asString)
		//		mockSO.On("ParseToken", token).Return(s, nil)
		//		mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(s, nil)
		//		mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
		//		mockAMQP.On("QosQueue").Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
		//		mockSO.On("UpdateStatus", mock.Anything, s, models.Processing).Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
		//		mockSO.On("ConvertToType", mock.Anything).Return(models.ResultedImage{ID: s, Name: "name", Location: "location"}, nil)
		//		mockSO.On("UpdateStatus", mock.Anything, s, models.Done).Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Done)).Return(nil).Return(nil)
		//		mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
		//	},
		//	expectedStatusCode:   200,
		//	expectedResponseBody: "{\"Image ID\":1}\n",
		// },
		{
			name:         "Inequality of identifiers",
			headerNames:  []string{"Authorization"},
			headerValues: []string{"Bearer token"},
			token:        "token",
			userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			fn: func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage) {
				asString := "00000000-0000-0000-0000-000000000001"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"users IDs do not match\"}\n",
		},
		{
			name:         "Failed upload file",
			headerNames:  []string{"Authorization", "Content-Type"},
			headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
			token:        "token",
			userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			fn: func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(uplImg.ID, fmt.Errorf("unable to insert image into database")).Once()
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot upload the file\"}\n",
		},
		// {
		//	name:         "Failed update status",
		//	headerNames:  []string{"Authorization", "Content-Type"},
		//	headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
		//	token:        "token",
		//	userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		//	fn: func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage) {
		//		asString := "00000000-0000-0000-0000-000000000000"
		//		s := uuid.MustParse(asString)
		//		mockSO.On("ParseToken", token).Return(s, nil)
		//		mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(s, nil)
		//		//mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
		//		//mockAMQP.On("QosQueue").Return(nil)
		//		//mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
		//		//mockSO.On("UpdateStatus", mock.Anything, s, models.Processing).Return(fmt.Errorf("cannot update image status"))
		//	},
		//	expectedStatusCode:   500,
		//	expectedResponseBody: "{\"error\":\"cannot update image status\"}\n",
		// },
		// {
		//	name:         "Failed convert image",
		//	headerNames:  []string{"Authorization", "Content-Type"},
		//	headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
		//	token:        "token",
		//	userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		//	fn: func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage) {
		//		asString := "00000000-0000-0000-0000-000000000000"
		//		s := uuid.MustParse(asString)
		//		mockSO.On("ParseToken", token).Return(s, nil)
		//		mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(s, nil)
		//		mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
		//		mockAMQP.On("QosQueue").Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
		//		mockSO.On("UpdateStatus", mock.Anything, s, models.Processing).Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
		//		mockSO.On("ConvertToType", mock.Anything).Return(models.ResultedImage{}, fmt.Errorf("cannot convert"))
		//	},
		//	expectedStatusCode:   500,
		//	expectedResponseBody: "{\"error\":\"cannot convert\"}\n",
		// },
		// {
		//	name:         "Failed convert image 2",
		//	headerNames:  []string{"Authorization", "Content-Type"},
		//	headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
		//	token:        "token",
		//	userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		//	fn: func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage) {
		//		asString := "00000000-0000-0000-0000-000000000000"
		//		s := uuid.MustParse(asString)
		//		mockSO.On("ParseToken", token).Return(s, nil)
		//		mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(s, nil)
		//		mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
		//		mockAMQP.On("QosQueue").Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
		//		mockSO.On("UpdateStatus", mock.Anything, s, models.Processing).Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
		//		mockSO.On("ConvertToType", mock.Anything).Return(models.ResultedImage{ID: s, Name: "name", Location: "location"}, nil)
		//		mockSO.On("UpdateStatus", mock.Anything, s, models.Done).Return(fmt.Errorf("cannot update image status"))
		//	},
		//	expectedStatusCode:   500,
		//	expectedResponseBody: "{\"error\":\"cannot update image status\"}\n",
		// },
		// {
		//	name:         "Failed create request",
		//	headerNames:  []string{"Authorization", "Content-Type"},
		//	headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
		//	token:        "token",
		//	userID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		//	fn: func(mockSO *mocks.ServiceOperations, mockAMQP *mocks.AMQP, token string, uplImg models.UploadedImage) {
		//		asString := "00000000-0000-0000-0000-000000000000"
		//		s := uuid.MustParse(asString)
		//		mockSO.On("ParseToken", token).Return(s, nil)
		//		mockSO.On("UploadImage", mock.Anything, mock.Anything).Return(s, nil)
		//		mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
		//		mockAMQP.On("QosQueue").Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
		//		mockSO.On("UpdateStatus", mock.Anything, s, models.Processing).Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
		//		mockSO.On("ConvertToType", mock.Anything).Return(models.ResultedImage{ID: s, Name: "name", Location: "location"}, nil)
		//		mockSO.On("UpdateStatus", mock.Anything, s, models.Done).Return(nil)
		//		mockAMQP.On("Publish", "", q.Name, string(models.Done)).Return(nil).Return(nil)
		//		mockSO.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, fmt.Errorf("unable to insert resulted image into database"))
		//	},
		//	expectedStatusCode:   500,
		//	expectedResponseBody: "{\"error\":\"unable to insert resulted image into database\"}\n",
		// },
	}

	convertURL := "/api/user/%s/convert"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSO := new(mocks.ServiceOperations)
			mockAMQP := new(mocks.AMQP)

			tt.fn(mockSO, mockAMQP, tt.token, uplImg)

			service := &Service{ServiceOperations: mockSO}
			mq := broker.NewAMQPBroker()
			s := Server{router: mux.NewRouter(), service: service, mq: mq}

			s.router.HandleFunc(fmt.Sprintf(convertURL, "{userID}"),
				s.authorize(s.convertImage())).Methods(http.MethodPost)

			imageBytes := []byte("1111")
			buf := &bytes.Buffer{}
			writer := multipart.NewWriter(buf)
			header := make(textproto.MIMEHeader)
			header.Set("Content-Disposition", `form-data; name="uploadFile"; filename="filename.jpeg"`)
			header.Set("Content-Type", "image/jpeg")
			part, err := writer.CreatePart(header)
			require.NoError(t, err)
			_, err = part.Write(imageBytes)
			require.NoError(t, err)
			err = writer.Close()
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf(convertURL, tt.userID),
				buf)

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
		})
	}
}

func TestHandler_findConvertedImage(t *testing.T) {
	type fnBehavior func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool)

	resultedImage := models.ResultedImage{
		ID:       [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
		Name:     "filename",
		Location: "location",
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
			userID:      [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("FindTheResultingImage", mock.Anything, convertedID, models.Conversion).Return(resultedImage, nil)
				mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.Image{}, nil)
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
			userID:      [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				if isOriginal {
					mockSO.On("FindOriginalImage", mock.Anything, convertedID).Return(models.UploadedImage{}, nil)
					mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.Image{}, nil)
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
			userID:      [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
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
			name:        "Inequality of identifiers",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			convertedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			userID:      [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000001"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"users IDs do not match\"}\n",
		},
		{
			name:        "Wrong to find image",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			convertedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			userID:      [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("FindTheResultingImage", mock.Anything, mock.Anything, mock.Anything).Return(models.ResultedImage{}, utils.ErrFindImage)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot find image:cannot find image\"}\n",
		},
		{
			name:        "Incorrectly saved image",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			convertedID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			userID:      [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				mockSO.On("FindTheResultingImage", mock.Anything, convertedID, models.Conversion).Return(resultedImage, nil)
				mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.Image{}, utils.ErrSaveImage)
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
			userID:      [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				if isOriginal {
					mockSO.On("FindOriginalImage", mock.Anything, convertedID).Return(models.UploadedImage{}, utils.ErrFindImage)
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
			userID:      [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
			token:       "token",
			fn: func(mockSO *mocks.ServiceOperations, token string, convertedID uuid.UUID, isOriginal bool) {
				asString := "00000000-0000-0000-0000-000000000000"
				s := uuid.MustParse(asString)
				mockSO.On("ParseToken", token).Return(s, nil)
				if isOriginal {
					mockSO.On("FindOriginalImage", mock.Anything, convertedID).Return(models.UploadedImage{}, nil)
					mockSO.On("SaveImage", mock.Anything, mock.Anything, mock.Anything).Return(&models.Image{}, utils.ErrSaveImage)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot save image:cannot save image\"}\n",
		},
	}

	getConvertedURL := "/api/user/%s/convert/%s"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSO := new(mocks.ServiceOperations)

			tt.fn(mockSO, tt.token, tt.convertedID, tt.params.isOriginal)

			service := &Service{ServiceOperations: mockSO}
			s := Server{router: mux.NewRouter(), service: service}

			s.router.HandleFunc(fmt.Sprintf(getConvertedURL, "{userID}", "{convertedID}"),
				s.authorize(s.findConvertedImage())).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf(getConvertedURL, tt.userID, tt.convertedID), nil)

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
