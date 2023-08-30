package domain

type HTTPResponse struct {
	Data   any    `json:"data,omitempty"`
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

const UserCtxKey = "user"
