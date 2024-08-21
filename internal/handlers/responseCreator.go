package handlers

type SuccessResponse struct {
	data       map[string]interface{}
	success    bool
	statusCode string
}

type FailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
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
		Message: message,
		Success: false,
	}
}
