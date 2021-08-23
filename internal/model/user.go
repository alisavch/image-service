package model

// User contains information about user.
//
// A user is the security principal for this application.
//
// swagger:model
type User struct {
	// the ID for this user
	//
	// required: false
	ID int `json:"id,omitempty"`

	// the username for this user
	//
	// required: true
	Username string `json:"username,omitempty"`

	// the password for this user
	//
	// required: true
	Password string `json:"password,omitempty"`
}
