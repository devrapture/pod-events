package apperrors

import "errors"

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrSpotifyTokenNotFound is returned when a spotify token is not found
	ErrorSpotifyTokenNotFound = errors.New("spotify token not found")
	// ErrInvalidToken is returned when a token is invalid
	ErrInvalidToken = errors.New("invalid token")
)
