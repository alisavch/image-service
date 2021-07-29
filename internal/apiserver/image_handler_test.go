package apiserver

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/alisavch/image-service/internal/utils"
	"github.com/stretchr/testify/mock"

	"github.com/alisavch/image-service/internal/model"
	"github.com/alisavch/image-service/internal/service"
	"github.com/alisavch/image-service/internal/service/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestHandler_findUserHistory(t *testing.T) {
	type fnBehavior func(mockAuthorization *mocks.Authorization, token string)

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
			fn: func(mockAuthorization *mocks.Authorization, token string) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1\n",
		},
		{
			name:        "Test with invalid header name",
			headerName:  "",
			headerValue: "Bearer token",
			token:       "token",
			fn: func(r *mocks.Authorization, token string) {
				r.On("ParseToken", token).Return(1, nil)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"auth header is empty\"}\n",
		},
		{
			name:        "Test with invalid header value",
			headerName:  "Authorization",
			headerValue: "",
			token:       "token",
			fn: func(r *mocks.Authorization, token string) {
				r.On("ParseToken", token).Return(1, nil)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"auth header is empty\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := new(mocks.Authorization)
			tt.fn(auth, tt.token)

			services := &service.Service{Authorization: auth}
			s := Server{router: mux.NewRouter(), service: services}

			s.router.HandleFunc("/api/{userID}/history",
				s.authorize(s.findUserHistory())).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/1/history", nil)
			req.Header.Set(tt.headerName, tt.headerValue)
			s.ServeHTTP(w, req)
			require.Equal(t, tt.expectedStatusCode, w.Code)
			require.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_compressImage(t *testing.T) {
	//type fnBehavior func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, quality int, image model.UploadedImage)
	//
	//type params struct {
	//	name    string
	//	quality int
	//}
	//
	//tests := []struct {
	//	name                 string
	//	headerName           []string
	//	headerValue          []string
	//	inputImage           model.UploadedImage
	//	contentType          string
	//	params               params
	//	token                string
	//	fn                   fnBehavior
	//	expectedStatusCode   int
	//	expectedResponseBody string
	//}{
	//	{
	//		name:        "Test with correct values",
	//		headerName:  []string{"Authorization", "Content-Type"},
	//		headerValue: []string{"Bearer token"},
	//		inputImage:  model.UploadedImage{Name: "filename", Location: "location"},
	//		contentType: `multipart/form-data; boundary="foo123"`,
	//		params:      params{name: "ratio", quality: 95},
	//		token:       "token",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, quality int, image model.UploadedImage) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("CompressImage", mock.Anything, mock.Anything).Return(model.ResultedImage{ID: 1, Name: "filename", Location: "location"}, nil)
	//			mockImage.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	//			mockImage.On("UploadImage", mock.Anything).Return(1, nil)
	//
	//		},
	//		expectedStatusCode:   200,
	//		expectedResponseBody: "1\n",
	//	},
	//
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		auth := new(mocks.Authorization)
	//		image := new(mocks.Image)
	//		tt.fn(auth, image, tt.token, tt.params.quality, tt.inputImage)
	//
	//		services := &service.Service{Authorization: auth, Image: image}
	//		s := Server{router: mux.NewRouter(), service: services}
	//
	//		s.router.HandleFunc("/api/{userID}/compress",
	//			s.authorize(s.compressImage())).Methods(http.MethodPost)
	//
	//		w := httptest.NewRecorder()
	//		req := httptest.NewRequest(http.MethodPost, "/api/1/compress",
	//			ioutil.NopCloser(new(bytes.Buffer)))
	//
	//		q := req.URL.Query()
	//		q.Add(tt.params.name, string(rune(tt.params.quality)))
	//		req.URL.RawQuery = q.Encode()
	//
	//		req.Header.Set(tt.headerName[0], tt.headerValue[0])
	//		req.Header.Set(tt.headerName[1], tt.contentType)
	//
	//		if _, err := req.MultipartReader(); err != nil {
	//			t.Fatalf("MultipartReades: %v", err)
	//		}
	//
	//		s.ServeHTTP(w, req)
	//		require.Equal(t, tt.expectedStatusCode, w.Code)
	//		require.Equal(t, tt.expectedResponseBody, w.Body.String())
	//	})
	//}
}

func TestHandler_findCompressedImage(t *testing.T) {
	type fnBehavior func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool)

	resultedImage := model.ResultedImage{
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
			name:         "Test with all correct values with original",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			compressedID: 1,
			token:        "token",
			params:       params{name: "original", isOriginal: true},
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				mockImage.On("FindTheResultingImage", compressedID, model.Compression).Return(resultedImage, nil)
				mockImage.On("SaveImage", resultedImage.Name, "\\results\\", resultedImage.Name).Return(nil)
				if isOriginal {
					mockImage.On("FindOriginalImage", compressedID).Return(model.UploadedImage{ID: 1, Name: "filename", Location: "location"}, nil)
					mockImage.On("SaveImage", mock.Anything, "\\uploads\\", mock.Anything).Return(nil)
				}
			},
			expectedStatusCode:   200,
			expectedResponseBody: "\"successfully saved\"\n",
		},
		{
			name:         "Test with incorrect finding the image",
			headerName:   []string{"Authorization", "Content-Type"},
			headerValue:  []string{"Bearer token"},
			compressedID: 1,
			token:        "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				mockImage.On("FindTheResultingImage", compressedID, model.Compression).Return(model.ResultedImage{}, utils.ErrFindImage)
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
			params:       params{name: "original", isOriginal: true},
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				mockImage.On("FindTheResultingImage", compressedID, model.Compression).Return(resultedImage, nil)
				mockImage.On("SaveImage", resultedImage.Name, "\\results\\", resultedImage.Name).Return(utils.ErrSaveImage)
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
				mockImage.On("FindTheResultingImage", compressedID, model.Compression).Return(resultedImage, nil)
				mockImage.On("SaveImage", resultedImage.Name, "\\results\\", resultedImage.Name).Return(nil)
				if isOriginal {
					mockImage.On("FindOriginalImage", compressedID).Return(model.UploadedImage{}, utils.ErrFindImage)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot find image\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthorization := new(mocks.Authorization)
			mockImage := new(mocks.Image)
			tt.fn(mockAuthorization, mockImage, tt.token, tt.compressedID, tt.params.isOriginal)

			services := &service.Service{Authorization: mockAuthorization, Image: mockImage}
			s := Server{router: mux.NewRouter(), service: services}

			s.router.HandleFunc("/api/{userID}/compress/{id}",
				s.authorize(s.findCompressedImage())).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/1/compress/1", nil)

			q := req.URL.Query()
			q.Add(tt.params.name, strconv.FormatBool(tt.params.isOriginal))
			req.URL.RawQuery = q.Encode()

			req.Header.Set(tt.headerName[0], tt.headerValue[0])

			s.ServeHTTP(w, req)
			require.Equal(t, tt.expectedStatusCode, w.Code)
			require.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_convertImage(t *testing.T) {
	//type fnBehavior func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, image model.UploadedImage)
	//
	//tests := []struct {
	//	name                 string
	//	headerName           []string
	//	headerValue          []string
	//	inputImage           model.UploadedImage
	//	contentType          string
	//	token                string
	//	fn                   fnBehavior
	//	expectedStatusCode   int
	//	expectedResponseBody string
	//}{
	//	{
	//		name:        "Test with correct values",
	//		headerName:  []string{"Authorization", "Content-Type"},
	//		headerValue: []string{"Bearer token"},
	//		inputImage:  model.UploadedImage{Name: "filename", Location: "location"},
	//		contentType: `multipart/form-data; boundary="foo123"`,
	//		token:       "token",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, image model.UploadedImage) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("UploadImage", image).Return(1, nil)
	//			mockImage.On("ConvertToType", mock.Anything).Return(model.ResultedImage{ID:1, Name: "filename", Location: "location", Service: model.Conversion}, nil)
	//			mockImage.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, model.Request{}).Return(1, nil)
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "1\n",
	//	},
	//	//{
	//	//	name:        "Test with incorrect converting",
	//	//	headerName:  []string{"Authorization", "Content-Type"},
	//	//	headerValue: []string{"Bearer token"},
	//	//	inputImage:  model.UploadedImage{Name: "filename", Location: "location"},
	//	//	contentType: `multipart/form-data; boundary="foo123"`,
	//	//	token:       "token",
	//	//	fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, image model.UploadedImage) {
	//	//		mockAuthorization.On("ParseToken", token).Return(1, nil)
	//	//		mockImage.On("ConvertToType", mock.Anything).Return(model.ResultedImage{}, utils.ErrConvert)
	//	//	},
	//	//	expectedStatusCode:   500,
	//	//	expectedResponseBody: "{\"error\":\"cannot convert\"}\n",
	//	//},
	//	//{
	//	//	name:        "Test with incorrect creating request",
	//	//	headerName:  []string{"Authorization", "Content-Type"},
	//	//	headerValue: []string{"Bearer token"},
	//	//	inputImage:  model.UploadedImage{Name: "filename", Location: "location"},
	//	//	contentType: `multipart/form-data; boundary="foo123"`,
	//	//	token:       "token",
	//	//	fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, image model.UploadedImage) {
	//	//		mockAuthorization.On("ParseToken", token).Return(1, nil)
	//	//		mockImage.On("ConvertToType", mock.Anything).Return(model.ResultedImage{ID:1, Name: "filename", Location: "location", Service: model.Conversion}, nil)
	//	//		mockImage.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, model.Request{}).Return(0, utils.ErrCreateRequest)
	//	//	},
	//	//	expectedStatusCode:   500,
	//	//	expectedResponseBody: "{\"error\":\"cannot create request with image\"}\n",
	//	//},
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		auth := new(mocks.Authorization)
	//		image := new(mocks.Image)
	//		tt.fn(auth, image, tt.token, tt.inputImage)
	//
	//		services := &service.Service{Authorization: auth, Image: image}
	//		s := Server{router: mux.NewRouter(), service: services}
	//
	//		s.router.HandleFunc("/api/{userID}/convert",
	//			s.authorize(s.convertImage())).Methods(http.MethodPost)
	//
	//		w := httptest.NewRecorder()
	//		req := httptest.NewRequest(http.MethodPost, "/api/1/convert",
	//			ioutil.NopCloser(new(bytes.Buffer)))
	//
	//		req.Header.Set(tt.headerName[0], tt.headerValue[0])
	//		req.Header.Set(tt.headerName[1], tt.contentType)
	//
	//		s.ServeHTTP(w, req)
	//		require.Equal(t, tt.expectedStatusCode, w.Code)
	//		require.Equal(t, tt.expectedResponseBody, w.Body.String())
	//	})
	//}
}

func TestHandler_findConvertedImage(t *testing.T) {
	type fnBehavior func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool)

	resultedImage := model.ResultedImage{
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
				mockImage.On("FindTheResultingImage", convertedID, model.Conversion).Return(resultedImage, nil)
				mockImage.On("SaveImage", resultedImage.Name, "\\results\\", resultedImage.Name).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "\"successfully saved\"\n",
		},
		{
			name:        "Test with incorrect finding the image",
			headerName:  []string{"Authorization", "Content-Type"},
			headerValue: []string{"Bearer token"},
			convertedID: 1,
			token:       "token",
			fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool) {
				mockAuthorization.On("ParseToken", token).Return(1, nil)
				mockImage.On("FindTheResultingImage", convertedID, model.Conversion).Return(model.ResultedImage{}, utils.ErrFindImage)
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
				mockImage.On("FindTheResultingImage", convertedID, model.Conversion).Return(resultedImage, nil)
				mockImage.On("SaveImage", resultedImage.Name, "\\results\\", resultedImage.Name).Return(utils.ErrSaveImage)
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
				mockImage.On("FindTheResultingImage", convertedID, model.Conversion).Return(resultedImage, nil)
				mockImage.On("SaveImage", resultedImage.Name, "\\results\\", resultedImage.Name).Return(nil)
				if isOriginal {
					mockImage.On("FindOriginalImage", convertedID).Return(model.UploadedImage{}, utils.ErrFindImage)
				}
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"cannot find image\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthorization := new(mocks.Authorization)
			mockImage := new(mocks.Image)
			tt.fn(mockAuthorization, mockImage, tt.token, tt.convertedID, tt.params.isOriginal)

			services := &service.Service{Authorization: mockAuthorization, Image: mockImage}
			s := Server{router: mux.NewRouter(), service: services}

			s.router.HandleFunc("/api/{userID}/convert/{id}",
				s.authorize(s.findConvertedImage())).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/1/convert/1", nil)

			q := req.URL.Query()
			q.Add(tt.params.name, strconv.FormatBool(tt.params.isOriginal))
			req.URL.RawQuery = q.Encode()

			req.Header.Set(tt.headerName[0], tt.headerValue[0])

			s.ServeHTTP(w, req)
			require.Equal(t, tt.expectedStatusCode, w.Code)
			require.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
