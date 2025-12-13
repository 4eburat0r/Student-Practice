package errors

import "errors"

var (
	ErrResumeNotFound     = errors.New("resume not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden: you don't have permission to perform this action")
	ErrInvalidInput       = errors.New("invalid input data")
	ErrInternalServer     = errors.New("internal server error")
	ErrStudentNotFound    = errors.New("student not found")
	ErrInvalidResumeOwner = errors.New("you are not the owner of this resume")
)
