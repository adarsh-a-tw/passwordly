package common

type ErrorResponse struct {
	Message string `json:"message"`
}

func InternalServerError() ErrorResponse {
	return ErrorResponse{Message: "Something went wrong. Try again."}
}
