package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk  = "ok"
	StatusErr = "error"
)

func Ok() Response {
	return Response{
		Status: StatusOk,
	}
}

func Err(msg string) Response {
	return Response{
		Status: StatusErr,
		Error:  msg,
	}
}
