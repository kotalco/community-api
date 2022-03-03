package shared

type response struct {
	Data interface{} `json:"data"`
}

func NewResponse(data interface{}) interface{} {
	return response{data}
}
