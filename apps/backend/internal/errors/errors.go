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
	// ErrDeleteSubscriptionFailed is returned when a subscription cannot be deleted
	ErrDeleteSubscriptionFailed = errors.New("failed to delete subscription")
	// ErrSubscriptionAlreadyExists is returned when a subscription already exists
	ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
	// ErrSubscriptionNotFound is returned when a subscription cannot be found
	ErrSubscriptionNotFound = errors.New("subscription not found")
	// ErrSpotifyResourceNotFound is returned when a spotify resource cannot be found
	ErrSpotifyResourceNotFound = errors.New("spotify resource not found")
	// ErrSpotifyAuthorizationRequired is returned when Spotify authorization is missing or expired
	ErrSpotifyAuthorizationRequired = errors.New("spotify authorization required")
	// ErrSpotifyUnavailable is returned when Spotify cannot process a request
	ErrSpotifyUnavailable = errors.New("spotify is unavailable")
	// ErrInvalidSpotifyShowIDs is returned when no valid Spotify show IDs are provided
	ErrInvalidSpotifyShowIDs = errors.New("at least one valid spotify show id is required")
	// ErrTooManySpotifyShowIDs is returned when a subscription batch exceeds the allowed size
	ErrTooManySpotifyShowIDs = errors.New("too many spotify show ids")
)
