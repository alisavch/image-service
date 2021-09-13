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

	"github.com/alisavch/image-service/internal/utils"

	"github.com/stretchr/testify/mock"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/service"
	"github.com/alisavch/image-service/internal/service/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestHandler_findUserHistory(t *testing.T) {
	type fnBehavior func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, r *http.Request)

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
			name:        "Test with correct values",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, r *http.Request) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				mockImage.On("FindUserHistoryByID", mock.Anything, 1).Return([]models.History{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "[]\n",
		},
		{
			name:        "Test with failed conversion",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, r *http.Request) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrFailedConvert)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"failed convert to int userID\"}\n",
		},
		{
			name:        "Test with missing user ID",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, r *http.Request) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrMissingParams)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"id is missing in parameters\"}\n",
		},
		{
			name:        "Test with failed conversion userID",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, r *http.Request) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrAtoi)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"int conversion error\"}\n",
		},
		{
			name:        "Inequality test",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, r *http.Request) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrPrivacy)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"you can only view your data\"}\n",
		},
		{
			name:        "Test with incorrect values",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, r *http.Request) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				mockImage.On("FindUserHistoryByID", mock.Anything, 1).Return([]models.History{}, fmt.Errorf(""))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthorization := new(mocks.Authorization)
			mockImage := new(mocks.Image)

			services := &service.Service{Authorization: mockAuthorization, Image: mockImage}
			s := Server{router: mux.NewRouter(), service: services}

			s.router.HandleFunc("/api/user/{userID}/history",
				s.authorize(s.findUserHistory())).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/user/1/history", nil)

			tt.fn(mockAuthorization, mockImage, tt.token, req)

			req.Header.Set(tt.headerName, tt.headerValue)
			s.ServeHTTP(w, req)
			mockAuthorization.AssertExpectations(t)
			mockImage.AssertExpectations(t)
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
	//
	uplImg := models.UploadedImage{
		ID:       0,
		Name:     "filename.jpeg",
		Location: "location",
	}

	type fnBehavior func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, uplImg models.UploadedImage)

	tests := []struct {
		name                 string
		headerNames          []string
		headerValues         []string
		inputImage           models.UploadedImage
		params               params
		token                string
		fn                   fnBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "Test with invalid token",
			headerNames:  []string{"Authorization"},
			headerValues: []string{"Bearer token"},
			token:        "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, uplImg models.UploadedImage) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrFailedConvert)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"failed convert to int userID\"}\n",
		},
		{
			name:         "Test with missing user ID",
			headerNames:  []string{"Authorization"},
			headerValues: []string{"Bearer token"},
			token:        "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, uplImg models.UploadedImage) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrMissingParams)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"id is missing in parameters\"}\n",
		},
		{
			name:         "Test with failed conversion userID",
			headerNames:  []string{"Authorization"},
			headerValues: []string{"Bearer token"},
			token:        "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, uplImg models.UploadedImage) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrAtoi)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"int conversion error\"}\n",
		},
		{
			name:         "Inequality test",
			headerNames:  []string{"Authorization"},
			headerValues: []string{"Bearer token"},
			token:        "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, uplImg models.UploadedImage) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrPrivacy)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"you can only view your data\"}\n",
		},
		{
			name:         "Transform width error test",
			headerNames:  []string{"Authorization"},
			headerValues: []string{"Bearer token"},
			token:        "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, uplImg models.UploadedImage) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrAtoi)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"int conversion error\"}\n",
		},
		//{
		//	name:         "Test with correct values",
		//	headerNames:  []string{"Authorization", "Content-Type"},
		//	headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
		//	params:       params{name: "width", quantity: 100},
		//	token:        "token",
		//	fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, uplImg models.UploadedImage) {
		//		mockAuthorization.On("ParseToken", token).Return(1, nil)
		//		mockImage.On("UploadImage", mock.Anything, uplImg).Return(1, nil)
		//		//mockImage.On("UpdateStatus", mock.Anything, 1, models.Processing).Return(nil)
		//		//mockImage.On("CompressImage", 100, models.UploadedImage{Name: "filename", Location: "location"}).Return(models.ResultedImage{Name: "nnn", Location: "lll"}, nil)
		//		//
		//		//mockImage.On("UpdateStatus", mock.Anything, 1, models.Done).Return(nil)
		//		//mockImage.On("CreateRequest", models.User{Username: "u", Password: "p"},
		//		//	models.UploadedImage{ID: 1, Name: "n", Location: "l"},
		//		//	models.ResultedImage{ID: 1, Name: "n", Location: "l", Service: "s"},
		//		//	models.UserImage{ID: 1, UserAccountID: 1, UploadedImageID: 1, ResultedImageID: 1, Status: "s"},
		//		//	models.Request{ID: 1, UserImageID: 1}).Return(1, nil)
		//
		//	},
		//	expectedStatusCode:   200,
		//	expectedResponseBody: "{\"error\":\"token claims is invalid\"}\n",
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthorization := new(mocks.Authorization)
			mockImage := new(mocks.Image)

			tt.fn(mockAuthorization, mockImage, tt.token, uplImg)

			services := &service.Service{Authorization: mockAuthorization, Image: mockImage}
			s := Server{router: mux.NewRouter(), service: services}

			s.router.HandleFunc("/api/user/{userID}/compress",
				s.authorize(s.compressImage())).Methods(http.MethodPost)

			imageBytes := []byte("1111")
			buf := &bytes.Buffer{}
			writer := multipart.NewWriter(buf)
			header := make(textproto.MIMEHeader)
			header.Set("Content-Disposition", `form-data; name="uploadFile"; filename="filename.jpeg"`)
			header.Set("Content-Type", "image/jpeg")
			part, err := writer.CreatePart(header)
			require.NoError(t, err)
			_, err = part.Write(imageBytes)
			err = writer.Close()
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/user/1/compress", buf)

			body, err := ioutil.ReadAll(req.Body)
			require.NoError(t, err)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

			req.Header.Set("Authorization", "Bearer token")
			req.Header.Set("Content-Type", writer.FormDataContentType())
			err = req.Body.Close()
			require.NoError(t, err)
			s.ServeHTTP(w, req)
			mockAuthorization.AssertExpectations(t)
			mockImage.AssertExpectations(t)
			require.Equal(t, tt.expectedStatusCode, w.Code)
			require.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_findCompressedImage(t *testing.T) {
	type fnBehavior func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool)

	resultedImage := models.ResultedImage{
		ID:       1,
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
		compressedID         int
		params               params
		token                string
		fn                   fnBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "Test with correct saving the image",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			compressedID: 1,
			token:        "token",
			params:       params{name: "original", isOriginal: false},
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				mockImage.On("FindTheResultingImage", mock.Anything, compressedID, models.Compression).Return(resultedImage, nil)
				mockImage.On("SaveImage", resultedImage.Name, "/results/").Return(&models.Image{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name:        "Test with failed conversion userID",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrFailedConvert)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"failed convert to int userID\"}\n",
		},
		{
			name:        "Test with missing userID",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrMissingParams)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"id is missing in parameters\"}\n",
		},
		{
			name:        "Test with missing compressedID",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrMissingParams)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"id is missing in parameters\"}\n",
		},
		{
			name:        "Test with failed conversion userID",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrAtoi)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"int conversion error\"}\n",
		},
		{
			name:        "Inequality test",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrPrivacy)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"you can only view your data\"}\n",
		},
		{
			name:        "Test with failed conversion compressedImageID",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrAtoi)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"int conversion error\"}\n",
		},

		{
			name:         "Test with incorrect finding the image",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			compressedID: 1,
			token:        "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				mockImage.On("FindTheResultingImage", mock.Anything, compressedID, models.Compression).Return(models.ResultedImage{}, utils.ErrFindImage)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot find image\"}\n",
		},
		{
			name:         "Test with incorrect saving the image",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			compressedID: 1,
			token:        "token",
			params:       params{name: "original", isOriginal: false},
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				mockImage.On("FindTheResultingImage", mock.Anything, compressedID, models.Compression).Return(resultedImage, nil)
				mockImage.On("SaveImage", resultedImage.Name, "/results/").Return(&models.Image{}, utils.ErrSaveImage)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot save image\"}\n",
		},
		{
			name:         "Test with incorrect finding the original image",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			compressedID: 1,
			token:        "token",
			params:       params{name: "original", isOriginal: true},
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				if isOriginal {
					mockImage.On("FindOriginalImage", mock.Anything, 1).Return(models.UploadedImage{
						ID:       1,
						Name:     "filename",
						Location: "location",
					}, nil)
					mockImage.On("SaveImage", mock.Anything, "/uploads/").Return(&models.Image{}, utils.ErrSaveImage)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot save image\"}\n",
		},
		{
			name:         "Test with incorrect saving the original image",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			compressedID: 1,
			token:        "token",
			params:       params{name: "original", isOriginal: true},
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				if isOriginal {
					mockImage.On("FindOriginalImage", mock.Anything, compressedID).Return(models.UploadedImage{}, utils.ErrFindImage)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot find image\"}\n",
		},
		{
			name:         "Test with correct saving the original image",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			compressedID: 1,
			token:        "token",
			params:       params{name: "original", isOriginal: true},
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				if isOriginal {
					mockImage.On("FindOriginalImage", mock.Anything, compressedID).Return(models.UploadedImage{ID: 1, Name: "filename", Location: "location"}, nil)
					mockImage.On("SaveImage", mock.Anything, "/uploads/").Return(&models.Image{}, nil)
				}
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthorization := new(mocks.Authorization)
			mockImage := new(mocks.Image)
			tt.fn(mockAuthorization, mockImage, tt.token, tt.compressedID, tt.params.isOriginal)

			services := &service.Service{Authorization: mockAuthorization, Image: mockImage}
			s := Server{router: mux.NewRouter(), service: services}

			s.router.HandleFunc("/api/user/{userID}/compress/{compressedID}",
				s.authorize(s.findCompressedImage())).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/user/1/compress/1", nil)

			q := req.URL.Query()
			q.Add(tt.params.name, strconv.FormatBool(tt.params.isOriginal))
			req.URL.RawQuery = q.Encode()

			req.Header.Set(tt.headerName[0], tt.headerValue[0])

			s.ServeHTTP(w, req)
			mockAuthorization.AssertExpectations(t)
			mockImage.AssertExpectations(t)
			require.Equal(t, tt.expectedStatusCode, w.Code)
			require.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_convertImage(t *testing.T) {
	type fnBehavior func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string)

	tests := []struct {
		name                 string
		headerNames          []string
		headerValues         []string
		inputImage           models.UploadedImage
		contentType          string
		token                string
		fn                   fnBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		//{
		//	name:        "Test with incorrect converting",
		//	headerName:  []string{"Authorization", "Content-Type"},
		//	headerValue: []string{"Bearer token"},
		//	inputImage:  models.UploadedImage{Name: "filename", Location: "location"},
		//	contentType: `multipart/form-data; boundary="foo123"`,
		//	token:       "token",
		//	fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, image models.UploadedImage) {
		//		mockAuthorization.On("ParseToken", token).Return(1, nil)
		//		mockImage.On("ConvertToType", mock.Anything).Return(models.ResultedImage{}, fmt.Errorf("cannot convert"))
		//	},
		//	expectedStatusCode:   500,
		//	expectedResponseBody: "{\"error\":\"cannot convert\"}\n",
		//},
		//{
		//	name:        "Test with incorrect creating request",
		//	headerName:  []string{"Authorization", "Content-Type"},
		//	headerValue: []string{"Bearer token"},
		//	inputImage:  models.UploadedImage{Name: "filename", Location: "location"},
		//	contentType: `multipart/form-data; boundary="foo123"`,
		//	token:       "token",
		//	fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, image models.UploadedImage) {
		//		mockAuthorization.On("ParseToken", token).Return(1, nil)
		//		mockImage.On("ConvertToType", mock.Anything).Return(models.ResultedImage{ID:1, Name: "filename", Location: "location", Service: model.Conversion}, nil)
		//		mockImage.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, models.Request{}).Return(0, fmt.Errorf("cannot create request with image"))
		//	},
		//	expectedStatusCode:   500,
		//	expectedResponseBody: "{\"error\":\"cannot create request with image\"}\n",
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthorization := new(mocks.Authorization)
			mockImage := new(mocks.Image)

			services := &service.Service{Authorization: mockAuthorization, Image: mockImage}
			s := Server{router: mux.NewRouter(), service: services}

			s.router.HandleFunc("/api/user/{userID}/convert",
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
			err = writer.Close()
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/user/1/convert",
				buf)

			body, err := ioutil.ReadAll(req.Body)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

			req.Header.Set("Authorization", "Bearer token")
			req.Header.Set("Content-Type", writer.FormDataContentType())

			tt.fn(mockAuthorization, mockImage, tt.token)

			err = req.Body.Close()
			require.NoError(t, err)

			s.ServeHTTP(w, req)
			mockAuthorization.AssertExpectations(t)
			mockImage.AssertExpectations(t)
			require.Equal(t, tt.expectedStatusCode, w.Code)
			require.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_findConvertedImage(t *testing.T) {
	type fnBehavior func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool)

	resultedImage := models.ResultedImage{
		ID:       1,
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
		convertedID          int
		params               params
		token                string
		fn                   fnBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Test with all correct values",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			convertedID: 1,
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				mockImage.On("FindTheResultingImage", mock.Anything, convertedID, models.Conversion).Return(resultedImage, nil)
				mockImage.On("SaveImage", resultedImage.Name, "/results/").Return(&models.Image{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name:        "Test with all correct values for original",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			convertedID: 1,
			token:       "token",
			params:      params{name: "original", isOriginal: true},
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				if isOriginal {
					mockImage.On("FindOriginalImage", mock.Anything, convertedID).Return(models.UploadedImage{}, nil)
					mockImage.On("SaveImage", mock.Anything, "/uploads/").Return(&models.Image{}, nil)
				}
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name:        "Test with failed conversion userID",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrFailedConvert)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"failed convert to int userID\"}\n",
		},
		{
			name:        "Test with missing query param convertedImageID",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrRequest)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"incorrect request\"}\n",
		},
		{
			name:        "Test with missing userID",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrMissingParams)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"id is missing in parameters\"}\n",
		},
		{
			name:        "Test with failed conversion userID",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrAtoi)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"int conversion error\"}\n",
		},
		{
			name:        "Inequality test",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrPrivacy)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"you can only view your data\"}\n",
		},
		{
			name:        "Test with failed conversion convertedImageID",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, utils.ErrAtoi)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"int conversion error\"}\n",
		},
		{
			name:        "Test with incorrect finding the image",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			convertedID: 1,
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				mockImage.On("FindTheResultingImage", mock.Anything, convertedID, models.Conversion).Return(models.ResultedImage{}, utils.ErrFindImage)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot find image\"}\n",
		},
		{
			name:        "Test with incorrect saving the image",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			convertedID: 1,
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				mockImage.On("FindTheResultingImage", mock.Anything, convertedID, models.Conversion).Return(resultedImage, nil)
				mockImage.On("SaveImage", resultedImage.Name, "/results/").Return(&models.Image{}, utils.ErrSaveImage)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot save image\"}\n",
		},
		{
			name:        "Test with incorrect finding the original image",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			convertedID: 1,
			token:       "token",
			params:      params{name: "original", isOriginal: true},
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				if isOriginal {
					mockImage.On("FindOriginalImage", mock.Anything, convertedID).Return(models.UploadedImage{}, utils.ErrFindImage)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot find image\"}\n",
		},
		{
			name:        "Test with incorrect saving the original image",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			convertedID: 1,
			token:       "token",
			params:      params{name: "original", isOriginal: true},
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				if isOriginal {
					mockImage.On("FindOriginalImage", mock.Anything, convertedID).Return(models.UploadedImage{}, nil)
					mockImage.On("SaveImage", mock.Anything, "/uploads/").Return(&models.Image{}, utils.ErrSaveImage)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot save image\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthorization := new(mocks.Authorization)
			mockImage := new(mocks.Image)
			tt.fn(mockAuthorization, mockImage, tt.token, tt.convertedID, tt.params.isOriginal)

			services := &service.Service{Authorization: mockAuthorization, Image: mockImage}
			s := Server{router: mux.NewRouter(), service: services}

			s.router.HandleFunc("/api/user/{userID}/convert/{convertedID}",
				s.authorize(s.findConvertedImage())).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/user/1/convert/1", nil)

			q := req.URL.Query()
			q.Add(tt.params.name, strconv.FormatBool(tt.params.isOriginal))
			req.URL.RawQuery = q.Encode()

			req.Header.Set(tt.headerName[0], tt.headerValue[0])

			s.ServeHTTP(w, req)
			mockAuthorization.AssertExpectations(t)
			mockImage.AssertExpectations(t)
			require.Equal(t, tt.expectedStatusCode, w.Code)
			require.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
