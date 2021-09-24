package apiserver

import (
	"testing"
)

func TestHandler_findUserHistory(t *testing.T) {
	//type fnBehavior func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, r *http.Request)
	//
	//tests := []struct {
	//	name                 string
	//	headerName           string
	//	headerValue          string
	//	token                string
	//	userID               uuid.UUID
	//	fn                   fnBehavior
	//	expectedStatusCode   int
	//	expectedResponseBody string
	//}{
	//	{
	//		name:        "Test with correct values",
	//		headerName:  "Authorization",
	//		headerValue: "Bearer token",
	//		token:       "token",
	//		userID:      [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, r *http.Request) {
	//			asString := "00000000-0000-0000-0000-000000000000"
	//			s := uuid.MustParse(asString)
	//			mockAuthorization.On("ParseToken", token).Return(asString, nil)
	//			mockImage.On("FindUserHistoryByID", mock.Anything, s).Return([]models.History{}, nil)
	//		},
	//		expectedStatusCode:   200,
	//		expectedResponseBody: "[]\n",
	//	},
	//	{
	//		name:        "Inequality of identifiers",
	//		headerName:  "Authorization",
	//		headerValue: "Bearer token",
	//		token:       "token",
	//		//userID:      "11",
	//		userID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, r *http.Request) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//		},
	//		expectedStatusCode:   401,
	//		expectedResponseBody: "{\"error\":\"you can only view your data\"}\n",
	//	},
	//	{
	//		name:        "Invalid user ID",
	//		headerName:  "Authorization",
	//		headerValue: "Bearer token",
	//		token:       "token",
	//		//userID:      "a",
	//		userID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, r *http.Request) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//		},
	//		expectedStatusCode:   401,
	//		expectedResponseBody: "{\"error\":\"int conversion error\"}\n",
	//	},
	//	{
	//		name:        "Cannot complete request to get history",
	//		headerName:  "Authorization",
	//		headerValue: "Bearer token",
	//		token:       "token",
	//		//userID:      "1",
	//		userID: [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000},
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, r *http.Request) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("FindUserHistoryByID", mock.Anything, 1).Return([]models.History{}, fmt.Errorf("cannot complete request to get history"))
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"cannot complete request to get history\"}\n",
	//	},
	//}
	//
	//historyURL := "/api/user/%s/history"
	//
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		mockAuthorization := new(mocks.Authorization)
	//		mockImage := new(mocks.Image)
	//
	//		services := &service.Service{Authorization: mockAuthorization, Image: mockImage}
	//		s := Server{router: mux.NewRouter(), service: services}
	//
	//		s.router.HandleFunc(fmt.Sprintf(historyURL, "{userID}"),
	//			s.authorize(s.findUserHistory())).Methods(http.MethodGet)
	//
	//		w := httptest.NewRecorder()
	//		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf(historyURL, tt.userID), nil)
	//
	//		tt.fn(mockAuthorization, mockImage, tt.token, req)
	//
	//		req.Header.Set(tt.headerName, tt.headerValue)
	//		s.ServeHTTP(w, req)
	//		mockAuthorization.AssertExpectations(t)
	//		mockImage.AssertExpectations(t)
	//		require.Equal(t, tt.expectedStatusCode, w.Code)
	//		require.Equal(t, tt.expectedResponseBody, w.Body.String())
	//	})
	//}
}

func TestHandler_compressImage(t *testing.T) {
	//type params struct {
	//	name     string
	//	quantity int
	//}
	//
	//uplImg := models.UploadedImage{
	//	ID:       1,
	//	Name:     "filename.jpeg",
	//	Location: "location",
	//}
	//
	//q := amqp.Queue{Name: "", Messages: 1, Consumers: 1}
	//
	//type fnBehavior func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, mockAWS *mocks3.AmazonS3, token string, uplImg models.UploadedImage, isRemoteStorage bool)
	//
	//tests := []struct {
	//	name                 string
	//	headerNames          []string
	//	headerValues         []string
	//	inputImage           models.UploadedImage
	//	params               params
	//	token                string
	//	userID               string
	//	isRemoteStorage      bool
	//	fn                   fnBehavior
	//	expectedStatusCode   int
	//	expectedResponseBody string
	//}{
	//	{
	//		name:            "Test with correct values",
	//		headerNames:     []string{"Authorization", "Content-Type"},
	//		headerValues:    []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
	//		params:          params{name: "width", quantity: 100},
	//		token:           "token",
	//		userID:          "1",
	//		isRemoteStorage: true,
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, mockAWS *mocks3.AmazonS3, token string, uplImg models.UploadedImage, isRemoteStorage bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			if isRemoteStorage {
	//				mockAWS.On("UploadToAWS", mock.Anything, mock.Anything).Return(mock.Anything, nil)
	//				mockImage.On("UploadImage", mock.Anything, mock.Anything).Return(1, nil)
	//			}
	//			mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
	//			mockAMQP.On("QosQueue").Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Processing).Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
	//			mockImage.On("CompressImage", 100, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(models.ResultedImage{ID: 1, Name: "name", Location: "location"}, nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Done).Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Done)).Return(nil).Return(nil)
	//			mockImage.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	//		},
	//		expectedStatusCode:   200,
	//		expectedResponseBody: "{\"Image ID\":1}\n",
	//	},
	//	{
	//		name:         "Invalid user ID",
	//		headerNames:  []string{"Authorization"},
	//		headerValues: []string{"Bearer token"},
	//		token:        "token",
	//		userID:       "a",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, mockAWS *mocks3.AmazonS3, token string, uplImg models.UploadedImage, isRemoteStorage bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//		},
	//		expectedStatusCode:   401,
	//		expectedResponseBody: "{\"error\":\"int conversion error\"}\n",
	//	},
	//	{
	//		name:         "Inequality of identifiers",
	//		headerNames:  []string{"Authorization"},
	//		headerValues: []string{"Bearer token"},
	//		token:        "token",
	//		userID:       "2",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, mockAWS *mocks3.AmazonS3, token string, uplImg models.UploadedImage, isRemoteStorage bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//		},
	//		expectedStatusCode:   401,
	//		expectedResponseBody: "{\"error\":\"you can only view your data\"}\n",
	//	},
	//	{
	//		name:            "Error upload file",
	//		headerNames:     []string{"Authorization", "Content-Type"},
	//		headerValues:    []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
	//		params:          params{name: "width", quantity: 100},
	//		token:           "token",
	//		userID:          "1",
	//		isRemoteStorage: true,
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, mockAWS *mocks3.AmazonS3, token string, uplImg models.UploadedImage, isRemoteStorage bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			if isRemoteStorage {
	//				mockAWS.On("UploadToAWS", mock.Anything, mock.Anything).Return(mock.Anything, nil).Once()
	//				//mockImage.On("UploadImage", mock.Anything, mock.Anything).Return(uplImg.ID, fmt.Errorf("unable to insert image into database")).Once()
	//			}
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"error upload file\"}\n",
	//	},
	//	{
	//		name:         "Failed update status",
	//		headerNames:  []string{"Authorization", "Content-Type"},
	//		headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
	//		params:       params{name: "width", quantity: 100},
	//		token:        "token",
	//		userID:       "1",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, mockAWS *mocks3.AmazonS3, token string, uplImg models.UploadedImage, isRemoteStorage bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("UploadImage", mock.Anything, mock.Anything).Return(1, nil)
	//			mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
	//			mockAMQP.On("QosQueue").Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Processing).Return(fmt.Errorf("cannot update image status"))
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"cannot update image status\"}\n",
	//	},
	//	{
	//		name:         "Failed compress image",
	//		headerNames:  []string{"Authorization", "Content-Type"},
	//		headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
	//		params:       params{name: "width", quantity: 100},
	//		token:        "token",
	//		userID:       "1",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, mockAWS *mocks3.AmazonS3, token string, uplImg models.UploadedImage, isRemoteStorage bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("UploadImage", mock.Anything, mock.Anything).Return(1, nil)
	//			mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
	//			mockAMQP.On("QosQueue").Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Processing).Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
	//			mockImage.On("CompressImage", 100, mock.Anything).Return(models.ResultedImage{}, fmt.Errorf("cannot open image"))
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"cannot open image\"}\n",
	//	},
	//	{
	//		name:         "Failed convert image 2",
	//		headerNames:  []string{"Authorization", "Content-Type"},
	//		headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
	//		token:        "token",
	//		userID:       "1",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, mockAWS *mocks3.AmazonS3, token string, uplImg models.UploadedImage, isRemoteStorage bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("UploadImage", mock.Anything, mock.Anything).Return(1, nil)
	//			mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
	//			mockAMQP.On("QosQueue").Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Processing).Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
	//			mockImage.On("CompressImage", 150, mock.Anything).Return(models.ResultedImage{ID: 1, Name: "name", Location: "location"}, nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Done).Return(fmt.Errorf("cannot update image status"))
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"cannot update image status\"}\n",
	//	},
	//	{
	//		name:         "Failed create request",
	//		headerNames:  []string{"Authorization", "Content-Type"},
	//		headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
	//		params:       params{name: "width", quantity: 100},
	//		token:        "token",
	//		userID:       "1",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, mockAWS *mocks3.AmazonS3, token string, uplImg models.UploadedImage, isRemoteStorage bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("UploadImage", mock.Anything, mock.Anything).Return(1, nil)
	//			mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
	//			mockAMQP.On("QosQueue").Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Processing).Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
	//			mockImage.On("CompressImage", 100, mock.Anything).Return(models.ResultedImage{ID: 1, Name: "name", Location: "location"}, nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Done).Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Done)).Return(nil).Return(nil)
	//			mockImage.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, fmt.Errorf("unable to insert resulted image into database"))
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"unable to insert resulted image into database\"}\n",
	//	},
	//}
	//
	//compressURL := "/api/user/%s/compress"
	//
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		mockAuthorization := new(mocks.Authorization)
	//		mockImage := new(mocks.Image)
	//		mockAMQP := new(mocks2.AMQP)
	//		mockAWS := new(mocks3.AmazonS3)
	//
	//		tt.fn(mockAuthorization, mockImage, mockAMQP, mockAWS, tt.token, uplImg, tt.isRemoteStorage)
	//
	//		services := &service.Service{Authorization: mockAuthorization, Image: mockImage}
	//		mq := &broker.AMQPBroker{AMQP: mockAMQP}
	//		s := Server{router: mux.NewRouter(), service: services, mq: mq}
	//
	//		s.router.HandleFunc(fmt.Sprintf(compressURL, "{userID}"),
	//			s.authorize(s.compressImage())).Methods(http.MethodPost)
	//
	//		imageBytes := []byte("uploadFile: undefined")
	//		buf := &bytes.Buffer{}
	//		writer := multipart.NewWriter(buf)
	//		header := make(textproto.MIMEHeader)
	//		header.Set("Content-Disposition", `form-data; name="uploadFile"; filename="filename.jpeg"`)
	//		header.Set("Content-Type", "image/jpeg")
	//		part, err := writer.CreatePart(header)
	//		require.NoError(t, err)
	//		_, err = part.Write(imageBytes)
	//		require.NoError(t, err)
	//		err = writer.Close()
	//		require.NoError(t, err)
	//
	//		w := httptest.NewRecorder()
	//		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf(compressURL, tt.userID), buf)
	//
	//		body, err := ioutil.ReadAll(req.Body)
	//		require.NoError(t, err)
	//		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	//
	//		req.Header.Set("Authorization", "Bearer token")
	//		req.Header.Set("Content-Type", writer.FormDataContentType())
	//		req.Header.Set("Content-Length", "1000")
	//
	//		q := req.URL.Query()
	//		q.Add(tt.params.name, strconv.Itoa(tt.params.quantity))
	//		req.URL.RawQuery = q.Encode()
	//
	//		err = req.Body.Close()
	//		require.NoError(t, err)
	//		s.ServeHTTP(w, req)
	//		mockAuthorization.AssertExpectations(t)
	//		mockImage.AssertExpectations(t)
	//		require.Equal(t, tt.expectedStatusCode, w.Code)
	//		require.Equal(t, tt.expectedResponseBody, w.Body.String())
	//	})
	//}
}

func TestHandler_findCompressedImage(t *testing.T) {
	//type fnBehavior func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool)
	//
	//resultedImage := models.ResultedImage{
	//	ID:       1,
	//	Name:     "filename",
	//	Location: "Location",
	//}
	//
	//type params struct {
	//	name       string
	//	isOriginal bool
	//}
	//
	//tests := []struct {
	//	name                 string
	//	headerName           []string
	//	headerValue          []string
	//	params               params
	//	token                string
	//	compressedID         string
	//	userID               string
	//	fn                   fnBehavior
	//	expectedStatusCode   int
	//	expectedResponseBody string
	//}{
	//	{
	//		name:         "Test with correct saving the image",
	//		headerName:   []string{"Authorization", "Content-Type"},
	//		headerValue:  []string{"Bearer token"},
	//		compressedID: "1",
	//		token:        "token",
	//		userID:       "1",
	//		params:       params{name: "original", isOriginal: false},
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("FindTheResultingImage", mock.Anything, compressedID, models.Compression).Return(resultedImage, nil)
	//			mockImage.On("SaveImage", resultedImage.Name, "/results/").Return(&models.Image{}, nil)
	//		},
	//		expectedStatusCode:   200,
	//		expectedResponseBody: "",
	//	},
	//	{
	//		name:         "Invalid user ID\"",
	//		headerName:   []string{"Authorization", "Content-Type"},
	//		headerValue:  []string{"Bearer token"},
	//		token:        "token",
	//		userID:       "a",
	//		compressedID: "1",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//		},
	//		expectedStatusCode:   401,
	//		expectedResponseBody: "{\"error\":\"int conversion error\"}\n",
	//	},
	//	{
	//		name:         "Inequality of identifiers",
	//		headerName:   []string{"Authorization", "Content-Type"},
	//		headerValue:  []string{"Bearer token"},
	//		token:        "token",
	//		userID:       "12",
	//		compressedID: "1",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//		},
	//		expectedStatusCode:   401,
	//		expectedResponseBody: "{\"error\":\"you can only view your data\"}\n",
	//	},
	//	{
	//		name:         "Wrong to find image",
	//		headerName:   []string{"Authorization", "Content-Type"},
	//		headerValue:  []string{"Bearer token"},
	//		compressedID: "1",
	//		token:        "token",
	//		userID:       "1",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("FindTheResultingImage", mock.Anything, compressedID, models.Compression).Return(models.ResultedImage{}, utils.ErrFindImage)
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"cannot find image\"}\n",
	//	},
	//	{
	//		name:         "Incorrectly saved image",
	//		headerName:   []string{"Authorization", "Content-Type"},
	//		headerValue:  []string{"Bearer token"},
	//		compressedID: "1",
	//		token:        "token",
	//		userID:       "1",
	//		params:       params{name: "original", isOriginal: false},
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("FindTheResultingImage", mock.Anything, compressedID, models.Compression).Return(resultedImage, nil)
	//			mockImage.On("SaveImage", resultedImage.Name, "/results/").Return(&models.Image{}, utils.ErrSaveImage)
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"cannot save image\"}\n",
	//	},
	//	{
	//		name:         "Incorrectly saved original image",
	//		headerName:   []string{"Authorization", "Content-Type"},
	//		headerValue:  []string{"Bearer token"},
	//		compressedID: "1",
	//		token:        "token",
	//		userID:       "1",
	//		params:       params{name: "original", isOriginal: true},
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			if isOriginal {
	//				mockImage.On("FindOriginalImage", mock.Anything, compressedID).Return(models.UploadedImage{}, utils.ErrFindImage)
	//			}
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"cannot find image\"}\n",
	//	},
	//	{
	//		name:         "Wrong to find original image",
	//		headerName:   []string{"Authorization", "Content-Type"},
	//		headerValue:  []string{"Bearer token"},
	//		compressedID: "1",
	//		token:        "token",
	//		userID:       "1",
	//		params:       params{name: "original", isOriginal: true},
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			if isOriginal {
	//				mockImage.On("FindOriginalImage", mock.Anything, compressedID).Return(models.UploadedImage{ID: 1, Name: "filename", Location: "location"}, nil)
	//				mockImage.On("SaveImage", mock.Anything, "/uploads/").Return(&models.Image{}, nil)
	//			}
	//		},
	//		expectedStatusCode:   200,
	//		expectedResponseBody: "",
	//	},
	//}
	//
	//getCompressedURL := "/api/user/%s/compress/%s"
	//
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		mockAuthorization := new(mocks.Authorization)
	//		mockImage := new(mocks.Image)
	//		compressInt, err := strconv.Atoi(tt.compressedID)
	//		require.NoError(t, err)
	//
	//		services := &service.Service{Authorization: mockAuthorization, Image: mockImage}
	//		s := Server{router: mux.NewRouter(), service: services}
	//
	//		s.router.HandleFunc(fmt.Sprintf(getCompressedURL, "{userID}", "{compressedID}"),
	//			s.authorize(s.findCompressedImage())).Methods(http.MethodGet)
	//
	//		w := httptest.NewRecorder()
	//		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf(getCompressedURL, tt.userID, tt.compressedID), nil)
	//
	//		tt.fn(mockAuthorization, mockImage, tt.token, compressInt, tt.params.isOriginal)
	//
	//		q := req.URL.Query()
	//		q.Add(tt.params.name, strconv.FormatBool(tt.params.isOriginal))
	//		req.URL.RawQuery = q.Encode()
	//
	//		req.Header.Set(tt.headerName[0], tt.headerValue[0])
	//
	//		s.ServeHTTP(w, req)
	//		mockAuthorization.AssertExpectations(t)
	//		mockImage.AssertExpectations(t)
	//		require.Equal(t, tt.expectedStatusCode, w.Code)
	//		require.Equal(t, tt.expectedResponseBody, w.Body.String())
	//	})
	//}
}

func TestHandler_convertImage(t *testing.T) {
	//type fnBehavior func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, token string, uplImg models.UploadedImage)
	//
	//uplImg := models.UploadedImage{
	//	ID:       0,
	//	Name:     "filename.jpeg",
	//	Location: "location",
	//}
	//
	//q := amqp.Queue{Name: "", Messages: 1, Consumers: 1}
	//
	//tests := []struct {
	//	name                 string
	//	headerNames          []string
	//	headerValues         []string
	//	inputImage           models.UploadedImage
	//	contentType          string
	//	token                string
	//	convertedID          string
	//	userID               string
	//	fn                   fnBehavior
	//	expectedStatusCode   int
	//	expectedResponseBody string
	//}{
	//	{
	//		name:         "Test with correct values",
	//		headerNames:  []string{"Authorization", "Content-Type"},
	//		headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
	//		token:        "token",
	//		userID:       "1",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, token string, uplImg models.UploadedImage) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("UploadImage", mock.Anything, mock.Anything).Return(1, nil)
	//			mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
	//			mockAMQP.On("QosQueue").Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Processing).Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
	//			mockImage.On("ConvertToType", mock.Anything).Return(models.ResultedImage{ID: 1, Name: "name", Location: "location"}, nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Done).Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Done)).Return(nil).Return(nil)
	//			mockImage.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	//		},
	//		expectedStatusCode:   200,
	//		expectedResponseBody: "{\"Image ID\":1}\n",
	//	},
	//	{
	//		name:         "Test with invalid token",
	//		headerNames:  []string{"Authorization"},
	//		headerValues: []string{"Bearer token"},
	//		token:        "token",
	//		userID:       "1",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, token string, uplImg models.UploadedImage) {
	//			mockAuthorization.On("ParseToken", token).Return(1, utils.ErrFailedConvert)
	//		},
	//		expectedStatusCode:   401,
	//		expectedResponseBody: "{\"error\":\"failed convert to int userID\"}\n",
	//	},
	//	{
	//		name:         "Invalid user ID",
	//		headerNames:  []string{"Authorization"},
	//		headerValues: []string{"Bearer token"},
	//		token:        "token",
	//		userID:       "a",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, token string, uplImg models.UploadedImage) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//		},
	//		expectedStatusCode:   401,
	//		expectedResponseBody: "{\"error\":\"int conversion error\"}\n",
	//	},
	//	{
	//		name:         "Inequality of identifiers",
	//		headerNames:  []string{"Authorization"},
	//		headerValues: []string{"Bearer token"},
	//		token:        "token",
	//		userID:       "2",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, token string, uplImg models.UploadedImage) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//		},
	//		expectedStatusCode:   401,
	//		expectedResponseBody: "{\"error\":\"you can only view your data\"}\n",
	//	},
	//	{
	//		name:         "Failed upload file",
	//		headerNames:  []string{"Authorization", "Content-Type"},
	//		headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
	//		token:        "token",
	//		userID:       "1",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, token string, uplImg models.UploadedImage) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("UploadImage", mock.Anything, mock.Anything).Return(uplImg.ID, fmt.Errorf("unable to insert image into database")).Once()
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"error upload file\"}\n",
	//	},
	//	{
	//		name:         "Failed update status",
	//		headerNames:  []string{"Authorization", "Content-Type"},
	//		headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
	//		token:        "token",
	//		userID:       "1",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, token string, uplImg models.UploadedImage) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("UploadImage", mock.Anything, mock.Anything).Return(1, nil)
	//			mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
	//			mockAMQP.On("QosQueue").Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Processing).Return(fmt.Errorf("cannot update image status"))
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"cannot update image status\"}\n",
	//	},
	//	{
	//		name:         "Failed convert image",
	//		headerNames:  []string{"Authorization", "Content-Type"},
	//		headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
	//		token:        "token",
	//		userID:       "1",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, token string, uplImg models.UploadedImage) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("UploadImage", mock.Anything, mock.Anything).Return(1, nil)
	//			mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
	//			mockAMQP.On("QosQueue").Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Processing).Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
	//			mockImage.On("ConvertToType", mock.Anything).Return(models.ResultedImage{}, fmt.Errorf("cannot convert"))
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"cannot convert\"}\n",
	//	},
	//	{
	//		name:         "Failed convert image 2",
	//		headerNames:  []string{"Authorization", "Content-Type"},
	//		headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
	//		token:        "token",
	//		userID:       "1",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, token string, uplImg models.UploadedImage) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("UploadImage", mock.Anything, mock.Anything).Return(1, nil)
	//			mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
	//			mockAMQP.On("QosQueue").Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Processing).Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
	//			mockImage.On("ConvertToType", mock.Anything).Return(models.ResultedImage{ID: 1, Name: "name", Location: "location"}, nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Done).Return(fmt.Errorf("cannot update image status"))
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"cannot update image status\"}\n",
	//	},
	//	{
	//		name:         "Failed create request",
	//		headerNames:  []string{"Authorization", "Content-Type"},
	//		headerValues: []string{"Bearer token", "image/jpeg", `multipart/form-data; boundary="foo123"`},
	//		token:        "token",
	//		userID:       "1",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, mockAMQP *mocks2.AMQP, token string, uplImg models.UploadedImage) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("UploadImage", mock.Anything, mock.Anything).Return(1, nil)
	//			mockAMQP.On("DeclareQueue", "publisher").Return(q, nil)
	//			mockAMQP.On("QosQueue").Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Queued)).Return(nil).Return(nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Processing).Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Processing)).Return(nil).Return(nil)
	//			mockImage.On("ConvertToType", mock.Anything).Return(models.ResultedImage{ID: 1, Name: "name", Location: "location"}, nil)
	//			mockImage.On("UpdateStatus", mock.Anything, 1, models.Done).Return(nil)
	//			mockAMQP.On("Publish", "", q.Name, string(models.Done)).Return(nil).Return(nil)
	//			mockImage.On("CreateRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, fmt.Errorf("unable to insert resulted image into database"))
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"unable to insert resulted image into database\"}\n",
	//	},
	//}
	//
	//convertURL := "/api/user/%s/convert"
	//
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		mockAuthorization := new(mocks.Authorization)
	//		mockImage := new(mocks.Image)
	//		mockAMQP := new(mocks2.AMQP)
	//
	//		tt.fn(mockAuthorization, mockImage, mockAMQP, tt.token, uplImg)
	//
	//		services := &service.Service{Authorization: mockAuthorization, Image: mockImage}
	//		mq := &broker.AMQPBroker{AMQP: mockAMQP}
	//		s := Server{router: mux.NewRouter(), service: services, mq: mq}
	//
	//		s.router.HandleFunc(fmt.Sprintf(convertURL, "{userID}"),
	//			s.authorize(s.convertImage())).Methods(http.MethodPost)
	//
	//		imageBytes := []byte("1111")
	//		buf := &bytes.Buffer{}
	//		writer := multipart.NewWriter(buf)
	//		header := make(textproto.MIMEHeader)
	//		header.Set("Content-Disposition", `form-data; name="uploadFile"; filename="filename.jpeg"`)
	//		header.Set("Content-Type", "image/jpeg")
	//		part, err := writer.CreatePart(header)
	//		require.NoError(t, err)
	//		_, err = part.Write(imageBytes)
	//		require.NoError(t, err)
	//		err = writer.Close()
	//		require.NoError(t, err)
	//
	//		w := httptest.NewRecorder()
	//		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf(convertURL, tt.userID),
	//			buf)
	//
	//		body, err := ioutil.ReadAll(req.Body)
	//		require.NoError(t, err)
	//		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	//
	//		req.Header.Set("Authorization", "Bearer token")
	//		req.Header.Set("Content-Type", writer.FormDataContentType())
	//
	//		err = req.Body.Close()
	//		require.NoError(t, err)
	//
	//		s.ServeHTTP(w, req)
	//		mockAuthorization.AssertExpectations(t)
	//		mockImage.AssertExpectations(t)
	//		require.Equal(t, tt.expectedStatusCode, w.Code)
	//		require.Equal(t, tt.expectedResponseBody, w.Body.String())
	//	})
	//}
}

func TestHandler_findConvertedImage(t *testing.T) {
	//type fnBehavior func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool)
	//
	//resultedImage := models.ResultedImage{
	//	ID:       1,
	//	Name:     "filename",
	//	Location: "location",
	//}
	//
	//type params struct {
	//	name       string
	//	isOriginal bool
	//}
	//
	//tests := []struct {
	//	name                 string
	//	headerName           []string
	//	headerValue          []string
	//	convertedID          string
	//	params               params
	//	token                string
	//	userID               string
	//	fn                   fnBehavior
	//	expectedStatusCode   int
	//	expectedResponseBody string
	//}{
	//	{
	//		name:        "Test with all correct values",
	//		headerName:  []string{"Authorization", "Content-Type"},
	//		headerValue: []string{"Bearer token"},
	//		convertedID: "1",
	//		userID:      "1",
	//		token:       "token",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("FindTheResultingImage", mock.Anything, convertedID, models.Conversion).Return(resultedImage, nil)
	//			mockImage.On("SaveImage", resultedImage.Name, "/results/").Return(&models.Image{}, nil)
	//		},
	//		expectedStatusCode:   200,
	//		expectedResponseBody: "",
	//	},
	//	{
	//		name:        "Test with all correct values for original",
	//		headerName:  []string{"Authorization", "Content-Type"},
	//		headerValue: []string{"Bearer token"},
	//		convertedID: "1",
	//		userID:      "1",
	//		token:       "token",
	//		params:      params{name: "original", isOriginal: true},
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			if isOriginal {
	//				mockImage.On("FindOriginalImage", mock.Anything, convertedID).Return(models.UploadedImage{}, nil)
	//				mockImage.On("SaveImage", mock.Anything, "/uploads/").Return(&models.Image{}, nil)
	//			}
	//		},
	//		expectedStatusCode:   200,
	//		expectedResponseBody: "",
	//	},
	//	{
	//		name:        "Test with invalid token",
	//		headerName:  []string{"Authorization", "Content-Type"},
	//		headerValue: []string{"Bearer token"},
	//		token:       "token",
	//		convertedID: "1",
	//		userID:      "1",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, utils.ErrEmptyToken)
	//		},
	//		expectedStatusCode:   401,
	//		expectedResponseBody: "{\"error\":\"token is empty\"}\n",
	//	},
	//	{
	//		name:        "Invalid user ID",
	//		headerName:  []string{"Authorization", "Content-Type"},
	//		headerValue: []string{"Bearer token"},
	//		token:       "token",
	//		convertedID: "1",
	//		userID:      "a",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//		},
	//		expectedStatusCode:   401,
	//		expectedResponseBody: "{\"error\":\"int conversion error\"}\n",
	//	},
	//	{
	//		name:        "Inequality of identifiers",
	//		headerName:  []string{"Authorization", "Content-Type"},
	//		headerValue: []string{"Bearer token"},
	//		token:       "token",
	//		convertedID: "1",
	//		userID:      "2",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, compressedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//		},
	//		expectedStatusCode:   401,
	//		expectedResponseBody: "{\"error\":\"you can only view your data\"}\n",
	//	},
	//	{
	//		name:        "Failed to find image",
	//		headerName:  []string{"Authorization", "Content-Type"},
	//		headerValue: []string{"Bearer token"},
	//		convertedID: "1",
	//		userID:      "1",
	//		token:       "token",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("FindTheResultingImage", mock.Anything, convertedID, models.Conversion).Return(models.ResultedImage{}, utils.ErrFindImage)
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"cannot find image\"}\n",
	//	},
	//	{
	//		name:        "Failed to save image",
	//		headerName:  []string{"Authorization", "Content-Type"},
	//		headerValue: []string{"Bearer token"},
	//		convertedID: "1",
	//		userID:      "1",
	//		token:       "token",
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			mockImage.On("FindTheResultingImage", mock.Anything, convertedID, models.Conversion).Return(resultedImage, nil)
	//			mockImage.On("SaveImage", resultedImage.Name, "/results/").Return(&models.Image{}, utils.ErrSaveImage)
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"cannot save image\"}\n",
	//	},
	//	{
	//		name:        "Failed to find original image",
	//		headerName:  []string{"Authorization", "Content-Type"},
	//		headerValue: []string{"Bearer token"},
	//		convertedID: "1",
	//		userID:      "1",
	//		token:       "token",
	//		params:      params{name: "original", isOriginal: true},
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			if isOriginal {
	//				mockImage.On("FindOriginalImage", mock.Anything, convertedID).Return(models.UploadedImage{}, utils.ErrFindImage)
	//			}
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"cannot find image\"}\n",
	//	},
	//	{
	//		name:        "Failed to save original image",
	//		headerName:  []string{"Authorization", "Content-Type"},
	//		headerValue: []string{"Bearer token"},
	//		convertedID: "1",
	//		userID:      "1",
	//		token:       "token",
	//		params:      params{name: "original", isOriginal: true},
	//		fn: func(mockAuthorization *mocks.Authorization, mockImage *mocks.Image, token string, convertedID int, isOriginal bool) {
	//			mockAuthorization.On("ParseToken", token).Return(1, nil)
	//			if isOriginal {
	//				mockImage.On("FindOriginalImage", mock.Anything, convertedID).Return(models.UploadedImage{}, nil)
	//				mockImage.On("SaveImage", mock.Anything, "/uploads/").Return(&models.Image{}, utils.ErrSaveImage)
	//			}
	//		},
	//		expectedStatusCode:   500,
	//		expectedResponseBody: "{\"error\":\"cannot save image\"}\n",
	//	},
	//}
	//
	//getConvertedURL := "/api/user/%s/convert/%s"
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		mockAuthorization := new(mocks.Authorization)
	//		mockImage := new(mocks.Image)
	//
	//		convertedInt, err := strconv.Atoi(tt.convertedID)
	//		require.NoError(t, err)
	//		tt.fn(mockAuthorization, mockImage, tt.token, convertedInt, tt.params.isOriginal)
	//
	//		services := &service.Service{Authorization: mockAuthorization, Image: mockImage}
	//		s := Server{router: mux.NewRouter(), service: services}
	//
	//		s.router.HandleFunc(fmt.Sprintf(getConvertedURL, "{userID}", "{convertedID}"),
	//			s.authorize(s.findConvertedImage())).Methods(http.MethodGet)
	//
	//		w := httptest.NewRecorder()
	//		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf(getConvertedURL, tt.userID, tt.convertedID), nil)
	//
	//		q := req.URL.Query()
	//		q.Add(tt.params.name, strconv.FormatBool(tt.params.isOriginal))
	//		req.URL.RawQuery = q.Encode()
	//
	//		req.Header.Set(tt.headerName[0], tt.headerValue[0])
	//
	//		s.ServeHTTP(w, req)
	//		mockAuthorization.AssertExpectations(t)
	//		mockImage.AssertExpectations(t)
	//		require.Equal(t, tt.expectedStatusCode, w.Code)
	//		require.Equal(t, tt.expectedResponseBody, w.Body.String())
	//	})
	//}
}
