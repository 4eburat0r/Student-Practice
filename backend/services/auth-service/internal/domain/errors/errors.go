package errors

import "errors"

var (
    ErrInvalidCredentials = errors.New("invalid email or password")
    ErrUserAlreadyExists  = errors.New("user with this email already exists")
    ErrUserNotFound       = errors.New("user not found")
    ErrUnauthorized       = errors.New("unauthorized")
    
    ErrInvalidToken = errors.New("invalid or expired token")
    ErrTokenExpired = errors.New("token expired")

    ErrInternal      = errors.New("internal server error")
    ErrDatabaseError = errors.New("database error")
)
