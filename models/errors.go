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

	//ErrUsernamePassword = errors.New("Please Check again your email and Password")
)
//authValidation
var (
	ErrUsernamePassword = errors.New("Please Check again your email and Password")
	ErrInvalidOTP = errors.New("Your OTP invalid Please check again ")
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

var (
	BookingTypeRequired = errors.New("BookingType Required")
	BookingExpIdRequired = errors.New("BookingExpId Required")
	PromoIdRequired = errors.New("PromoId Required")
	PaymentMethodIdRequired = errors.New("PaymentMethodId Required")
	ExpPaymentIdRequired = errors.New("ExpPaymentId Required")
	StatusRequired = errors.New("Status Required")
	TotalPriceRequired = errors.New("TotalPrice Required")
	CurrencyRequired = errors.New("Currency Required")
)
