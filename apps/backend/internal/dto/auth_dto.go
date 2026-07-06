package dto

type AuthExchangeRequest struct {
	Code string `json:"code" binding:"required" example:"abc123"`
}
