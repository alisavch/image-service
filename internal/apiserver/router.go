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
	// swagger:operation POST /api/sign-up sign-up sign-up
	// ---
	// summary: Registers a user.
	// description: Could be any user.
	// parameters:
	// - name: user
	//   in: body
	//   description: the user to create
	//   schema:
	//     "$ref": "#/definitions/User"
	// responses:
	//   "201":
	//     description: user created successfully
	//   "401":
	//     description: unauthorized user
	//   "500":
	//     description: internal server error
	apiRouter.HandleFunc("/sign-up", s.signUp()).Methods(http.MethodPost)
	// swagger:operation POST /api/sign-in sign-in sign-in
	// ---
	// summary: Authorizes the user.
	// description: Only authorized user has access.
	// parameters:
	// - name: user
	//   in: body
	//   description: the user to create
	//   schema:
	//     "$ref": "#/definitions/User"
	// responses:
	//   "200":
	//     description: successful operation
	//   "403":
	//     description: not enough right
	//   "500":
	//     description: internal server error
	apiRouter.HandleFunc("/sign-in", s.signIn()).Methods(http.MethodPost)
}

func (s *Server) newUserRouter() {
	userRouter := s.router.PathPrefix("/api").Subrouter().PathPrefix("/user").Subrouter()
	// swagger:operation GET /api/user/{userID}/history history
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
	//     description: successful operation
	//   "400":
	//     description: bad request
	//   "500":
	//     description: internal server error
	userRouter.HandleFunc("/{userID}/history", s.authorize(s.findUserHistory())).Methods(http.MethodGet)
	// swagger:operation POST /api/user/{userID}/compress compress compress
	// ---
	// summary: Compresses the image.
	// description: Receives an image from an input form and compresses it, also allows you to enter width in the query string.
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
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/UploadedImage"
	// responses:
	//   "200":
	//     description: successful operation
	//   "400":
	//     description: bad request
	//   "500":
	//     description: internal server error
	userRouter.HandleFunc("/{userID}/compress", s.authorize(s.compressImage())).Methods(http.MethodPost)
	// swagger:operation GET /api/user/{userID}/compress/{compressedID} findCompressed
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
	//   in: query
	//   type: boolean
	//   required: false
	//   description: If true - the original will be saved
	// responses:
	//   "200":
	//     description: successful operation
	//   "400":
	//     description: bad request
	//   "500":
	//     description: internal server error
	userRouter.HandleFunc("/{userID}/compress/{compressedID}", s.authorize(s.findCompressedImage())).Methods(http.MethodGet)
	// swagger:operation POST /api/user/{userID}/convert convert convert
	// ---
	// summary: Converts the image.
	// description: Receives an image from an input form and converts it PNG to JPG and vice versa.
	// parameters:
	// - name: userID
	//   in: path
	//   description: userID to filter by id
	//   required: true
	//   type: integer
	// - name: uploadFile
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/UploadedImage"
	// responses:
	//   "200":
	//     description: successful operation
	//   "400":
	//     description: bad request
	//   "500":
	//     description: internal server error
	userRouter.HandleFunc("/{userID}/convert", s.authorize(s.convertImage())).Methods(http.MethodPost)
	// swagger:operation GET /api/user/{userID}/convert/{convertedID} findConverted
	// ---
	// summary: Finds the converted image.
	// description: Downloads the converted image and original if required.
	// parameters:
	// - name: userID
	//   in: path
	//   description: userID to filter by id
	//   required: true
	//   type: integer
	// - name: convertedID
	//   in: path
	//   description: convertedID to filter by id
	//   required: true
	//   type: integer
	// - name: original
	//   in: query
	//   type: boolean
	//   required: false
	//   description: If true - the original will be saved
	// responses:
	//   "200":
	//     description: successful operation
	//   "400":
	//     description: bad request
	//   "500":
	//     description: internal server error
	userRouter.HandleFunc("/{userID}/convert/{convertedID}", s.authorize(s.findConvertedImage())).Methods(http.MethodGet)
}
