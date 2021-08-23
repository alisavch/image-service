package apiserver

import (
	"net/http"
)

// ConfigureRouter registers a couple of URL paths and handlers.
func (s *Server) ConfigureRouter() {
	s.newAPIRouter()
	s.newUserRouter()
}

func (s *Server) newAPIRouter() {
	apiRouter := s.router.PathPrefix("/api").Subrouter()
	// swagger:route POST /api/sign-up sign-up sign-up
	// ---
	// summary: Registers a user.
	// description: Could be any user.
	// parameters:
	// - name: username
	//   in: body
	//   description: user object
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/User"
	// - name: password
	//   in: body
	//   description: user object
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/User"
	// responses:
	//   "200":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/badReq"
	//   "500":
	//     "$ref": "#/responses/internal"
	apiRouter.HandleFunc("/sign-up", s.signUp()).Methods(http.MethodPost)
	// swagger:route POST /api/sign-in sign-in sign-in
	// ---
	// summary: Authorizes the user.
	// description: Only authorized user has access.
	// parameters:
	// - name: username
	//   in: body
	//   description: user object
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/User"
	// - name: password
	//   in: body
	//   description: user object
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/User"
	// responses:
	//   "200":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/badReq"
	//   "500":
	//     "$ref": "#/responses/internal"
	apiRouter.HandleFunc("/sign-in", s.signIn()).Methods(http.MethodPost)
}

func (s *Server) newUserRouter() {
	userRouter := s.router.PathPrefix("/api").Subrouter().PathPrefix("/user").Subrouter()
	// swagger:route GET /api/user/{userID}/history history history
	// ---
	// summary: Finds users history.
	// description: Lists all queries created by user.
	// parameters:
	// - name: userID
	//   in: path
	//   description: userID to filter by id
	//   required: true
	//   type: integer
	// responses:
	//   "200":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/badReq"
	//   "500":
	//     "$ref": "#/responses/internal"
	userRouter.HandleFunc("/{userID}/history", s.authorize(s.findUserHistory())).Methods(http.MethodGet)
	// swagger:route POST /api/user/{userID}/compress compress compress
	// ---
	// summary: Compresses the image.
	// description: Receives an image from an input form and compresses it, also allows you to enter width in the query string.
	// consumes:
	// - multipart/form-data
	// parameters:
	// - name: userID
	//   in: path
	//   description: userID to filter by id
	//   required: true
	//   type: integer
	// - name: width
	//   in: query
	//   type: integer
	//   required: false
	// - name: uploadFile
	//   in: formData
	//	 type: file
	//   required: true
	//   description: The file to upload.
	// responses:
	//   "200":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/badReq"
	//   "500":
	//     "$ref": "#/responses/internal"
	userRouter.HandleFunc("/{userID}/compress", s.authorize(s.compressImage())).Methods(http.MethodPost)
	// swagger:route GET /api/user/{userID}/compress/{compressedID} findCompressed
	// ---
	// summary: Finds the compressed image.
	// description: Downloads the compressed image and original if required.
	// parameters:
	// - name: userID
	//   in: path
	//   description: userID to filter by id
	//   required: true
	//   type: integer
	// - name: compressedID
	//   in: path
	//   description: compressedID to filter by id
	//   required: true
	//   type: integer
	// - name: original
	//   in: path
	//   description: compressedID to filter by id
	//   required: true
	//   type: integer
	// - name: uploadFile
	//   in: formData
	//	 type: file
	//   required: true
	// - name: original
	//   in: query
	//   schema:
	//     type: boolean
	//   description: If true - the original will be saved
	// responses:
	//   "200":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/badReq"
	//   "500":
	//     "$ref": "#/responses/internal"
	userRouter.HandleFunc("/{userID}/compress/{compressedID}", s.authorize(s.findCompressedImage())).Methods(http.MethodGet)
	// swagger:route POST /api/user/{userID}/convert convert convert
	// ---
	// summary: Converts the image.
	// description: Receives an image from an input form and converts it PNG to JPG and vice versa.
	// consumes:
	// - multipart/form-data
	// parameters:
	// - name: userID
	//   in: path
	//   description: userID to filter by id
	//   required: true
	//   type: integer
	// - name: uploadFile
	//   in: formData
	//	 type: file
	//   required: true
	//   description: The file to upload.
	// responses:
	//   "200":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/badReq"
	//   "500":
	//     "$ref": "#/responses/internal"
	userRouter.HandleFunc("/{userID}/convert", s.authorize(s.convertImage())).Methods(http.MethodPost)
	// swagger:route POST /api/user/{userID}/convert/{convertedID} findConverted
	// ---
	// summary: Finds the converted image.
	// description: Downloads the converted image and original if required.
	// consumes:
	// - multipart/form-data
	// parameters:
	// - name: userID
	//   in: path
	//   description: userID to filter by id
	//   required: true
	//   type: integer
	// - name: uploadFile
	//   in: formData
	//	 type: file
	//   required: true
	//   description: The file to upload.
	// - name: original
	//   in: query
	//   schema:
	//     type: boolean
	//   description: If true - the original will be saved
	// responses:
	//   "200":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/badReq"
	//   "500":
	//     "$ref": "#/responses/internal"
	userRouter.HandleFunc("/{userID}/convert/{convertedID}", s.authorize(s.findConvertedImage())).Methods(http.MethodGet)
}
