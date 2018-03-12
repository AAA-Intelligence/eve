package db

import "errors"

// UserError are errors that are shown to the error
type UserError string

func (err UserError) Error() string {
	return string(err)
}

// ErrUserNameTaken is thrown when user can not be created because the name is allready taken
var ErrUserNameTaken = UserError("username allready taken")

// ErrWrongPassword is thrown if the user password combination is invalid
var ErrWrongPassword = UserError("wrong password")

// ErrUserNotExists is thrown if the user does not exist
var ErrUserNotExists = UserError("username not found")

// ErrInternalServerError is thrown if something went wrong which is not ment to be shown to the user
var ErrInternalServerError = UserError("internal server error")

// ErrNoUserCredentials is shown if no user credentials are provided
var ErrNoUserCredentials = UserError("missing user credentials")

// ErrConnectionClosed is thrown if a database request is made, but the connection to the database is closed
var ErrConnectionClosed = errors.New("connection is closed or not established jet")
