package apiserver

import (
	"net/http"
)

// ConfigureRouter registers a couple of URL paths and handlers.
func (s *Server) ConfigureRouter() {
	s.newAPIRouter()
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
	//     description: user registered successfully
	//   "400":
	//     description: bad request
	//   "409":
	//     description: user already exists
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
	//   "401":
	//     description: unauthorized user
	apiRouter.HandleFunc("/sign-in", s.signIn()).Methods(http.MethodPost)
	// swagger:operation GET /api/history history history
	// ---
	// summary: Finds users history.
	// description: Lists all queries created by user.
	// responses:
	//   "200":
	//     description: successful operation
	//   "401":
	//     description: unauthorized user
	//   "500":
	//     description: internal server error
	apiRouter.HandleFunc("/history", s.authorize(s.findUserHistory())).Methods(http.MethodGet)
	// swagger:operation POST /api/compress compress compress
	// ---
	// summary: Compresses the image.
	// description: Receives an image from an input form and compresses it, also allows you to enter width in the query string.
	// parameters:
	// - name: width
	//   in: query
	//   type: integer
	//   required: false
	// - name: uploadFile
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/Image"
	// responses:
	//   "202":
	//     description: request accepted
	//   "401":
	//     description: unauthorized user
	//   "500":
	//     description: internal server error
	apiRouter.HandleFunc("/compress", s.authorize(s.compressImage())).Methods(http.MethodPost)
	// swagger:operation GET /api/compress/{compressedID} findCompressedImage findCompressedImage
	// ---
	// summary: Finds the compressed image.
	// description: Downloads the compressed image and original if required.
	// parameters:
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
	//   "401":
	//     description: unauthorized user
	//   "404":
	//     description: image is being processed
	//   "500":
	//     description: internal server error
	apiRouter.HandleFunc("/compress/{compressedID}", s.authorize(s.findCompressedImage())).Methods(http.MethodGet)
	// swagger:operation POST /api/convert convert convert
	// ---
	// summary: Converts the image.
	// description: Receives an image from an input form and converts it PNG to JPG and vice versa.
	// parameters:
	// - name: uploadFile
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/Image"
	// responses:
	//   "202":
	//     description: request accepted
	//   "401":
	//     description: unauthorized user
	//   "500":
	//     description: internal server error
	apiRouter.HandleFunc("/convert", s.authorize(s.convertImage())).Methods(http.MethodPost)
	// swagger:operation GET /api/convert/{convertedID} findConvertedImage findConvertedImage
	// ---
	// summary: Finds the converted image.
	// description: Downloads the converted image and original if required.
	// parameters:
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
	//   "401":
	//     description: unauthorized user
	//   "404":
	//     description: image is being processed
	//   "500":
	//     description: internal server error
	apiRouter.HandleFunc("/convert/{convertedID}", s.authorize(s.findConvertedImage())).Methods(http.MethodGet)
}
