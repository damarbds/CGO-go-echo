package models

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("Your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrUnAuthorize = errors.New("Unauthorize")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("Your Item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("Bad Request")

	ErrUsernamePassword = errors.New("Please Check again your email and Password")
)

var	(
ValidationExpId = errors.New("ExpId Required")
// ErrNotFound will throw if the requested item is not exists
ValidationStatus = errors.New("Status Required")
// ErrConflict will throw if the current action already exists
ValidationBookedBy = errors.New("BookedBy Required")
// ErrConflict will throw if the current action already exists
ValidationBookedDate = errors.New("Booking Date Required")
)
