package server

import (
	"context"
	"fmt"
	"github.com/AskaryanKarine/BMSTU-ds-1/pkg/validation"
	"github.com/go-playground/validator/v10"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo *echo.Echo
	pr   personRepository
}

const gracefulShutdownDeadline = 10 * time.Second

func New(pr personRepository) *Server {
	e := echo.New()
	s := &Server{
		echo: e,
		pr:   pr,
	}

	s.echo.Validator = validation.MustRegisterCustomValidator(validator.New())

	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}))

	api := s.echo.Group("/api/v1")

	persons := api.Group("/persons")
	persons.POST("", s.createPerson)
	persons.GET("", s.getPersons)
	persons.GET("/:id", s.getPersonByID)
	persons.PATCH("/:id", s.updatePerson)
	persons.DELETE("/:id", s.deletePersonByID)

	return s
}

func (s *Server) Run(port int) {
	portStr := fmt.Sprintf(":%d", port)
	go func() {
		log.Info("server starting on", "port", portStr)
		if err := s.echo.Start(portStr); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownDeadline)
	defer cancel()

	log.Info("server shutting down")
	if err := s.echo.Server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
