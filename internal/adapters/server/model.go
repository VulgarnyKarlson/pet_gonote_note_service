package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Data   any    `json:"data,omitempty"`
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type Endpoint struct {
	Method      string
	Path        string
	Middlewares []string
	Handler     func(r *http.Request) (*Response, error)
}

type Server interface {
	Run() error
	Stop() error
	AddEndpoint(e Endpoint)
	Use(m mux.MiddlewareFunc, id string)
}
