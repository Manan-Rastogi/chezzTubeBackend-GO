package models

type ErrorResponse struct {
	Code   int    `json:"code"`
	Message string `json:"msg"`
}

type SuccessResponse struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

type Response struct {
	Status          int16           `json:"status"`
	ErrorResponse   ErrorResponse   `json:"error,omitempty"`
	SuccessResponse SuccessResponse `json:"data,omitempty"`
}
