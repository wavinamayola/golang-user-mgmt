package models

type ValidationErrorResponse struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Message     string `json:"message"`
}

type Response struct {
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}
