package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	Server struct{}
)

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetOpenAPI(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func (s *Server) DummyAPI(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
