package controllers

type GenericResponse struct {
	StatusCode  int //Application Status code
	Message     string
	MessageType string
	Data        interface{}
}

func CreateResponse(status int, message, messageType string, data interface{}) *GenericResponse {
	return &GenericResponse{
		StatusCode:  status,
		Message:     message,
		MessageType: messageType,
		Data:        data,
	}
}
func CreateSuccessResponse(message, messageType string, data interface{}) *GenericResponse {
	return CreateResponse(200, message, messageType, data)
}
func CreateFailureResponse(status int, message, messageType string, data interface{}) *GenericResponse {
	return CreateResponse(200, message, messageType, data)
}
