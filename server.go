package genesis

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

type Server struct {
	app *fiber.App
}

func (s *Server) FiberConfig() fiber.Config {
	readTimeoutSecondsCount := 10

	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
	}
}

func (s *Server) Run(port string, app *fiber.App) error {
	s.app = app
	return s.app.Listen(":" + port)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
