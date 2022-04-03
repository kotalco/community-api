package shared

type response struct {
	Data interface{} `json:"data"`
}

type SuccessMessage struct {
	Message string `json:"message"`
}

func NewResponse(data interface{}) interface{} {
	return response{data}
}
