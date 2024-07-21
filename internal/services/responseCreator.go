package services

type SuccessResponse struct {
	data       map[string]interface{}
	success    bool
	statusCode string
}

type FailResponse struct {
	success bool
	message string
}

func NewSuccessResponse(data map[string]interface{}) *SuccessResponse {
	return &SuccessResponse{
		data:       data,
		success:    true,
		statusCode: "200",
	}
}

func NewFailResponse(message string) *FailResponse {
	return &FailResponse{
		message: message,
		success: false,
	}
}
