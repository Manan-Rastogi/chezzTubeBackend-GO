package models

type ErrorResponse struct {
	Code   int    `json:"code,omitempty"`
	Message string `json:"msg,omitempty"`
}

type SuccessResponse struct {
	Code int `json:"code,omitempty"`
	Data any `json:"data,omitempty"`
}

type Response struct {
	Status          int16           `json:"status"`
	ErrorResponse   *ErrorResponse   `json:"error,omitempty"`
	SuccessResponse *SuccessResponse `json:"data,omitempty"`
}
