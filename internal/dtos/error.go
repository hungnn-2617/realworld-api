package dtos

type ErrorResponse struct {
	Errors ErrorBody `json:"errors"`
}

type ErrorBody struct {
	Body []string `json:"body"`
}
