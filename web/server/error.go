package server

type APIError struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func (e APIError) Error() string {
	return e.Msg
}
