package http

type Response struct {
	Data   any    `json:"data,omitempty"`
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

const UserCtxKey = "user"
