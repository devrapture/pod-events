package apperrors

import "errors"

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrSpotifyTokenNotFound is returned when a spotify token is not found
	ErrorSpotifyTokenNotFound = errors.New("spotify token not found")
	// ErrInvalidToken is returned when a token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrChannelIDNotFound is returned when a channel is not found
	ErrChannelIDNotFound = errors.New("channel id not found")
	// ErrInvalidTelegramBotToken is returned when a telegram bot token is invalid
	ErrInvalidTelegramBotToken = errors.New("invalid telegram bot token")
	// ErrTelegramBotTokenNotProvided is returned when a telegram bot token is not provided
	ErrTelegramBotTokenNotProvided = errors.New("telegram bot token is not provided")
	// ErrInvalidTelegramBotFormat is returned when a telegram bot token is invalid
	ErrInvalidTelegramBotFormat = errors.New("Telegram bot tokens look like: 7123456789:AAFxxxxxxxxxxxxx — get yours from @BotFather")
	// ErrNoExpiredTelegramConnectionsFound is returned when no expired or consumed telegram connections are found for cleanup
	ErrNoExpiredTelegramConnectionsFound = errors.New("no expired telegram connections found")
	// ErrTelegramChannelAlreadyExists is returned when a user already has a Telegram notification channel
	ErrTelegramChannelAlreadyExists = errors.New("telegram channel already exists")
	// ErrPodcastShowNotFound is returned when a podcast show cannot be found
	ErrPodcastShowNotFound = errors.New("podcast show not found")
	// ErrEpisodeNotFound is returned when an episode cannot be found
	ErrEpisodeNotFound = errors.New("episode not found")
)
