package app

import (
	"log"
	"speakeasy/internal/pkg/authentication"
	"speakeasy/internal/pkg/trip"

	"github.com/gin-gonic/gin"
)

// Server object which contains router and services
type Server struct {
	router                *gin.Engine
	authenticationService authentication.Service
	tripService           trip.Service
}

// NewServer returns Server object
func NewServer(
	router *gin.Engine,
	authenticationService authentication.Service,
	tripService trip.Service,
) *Server {
	return &Server{
		router:                router,
		authenticationService: authenticationService,
		tripService:           tripService,
	}
}

// Run function to run server
func (s *Server) Run() error {
	err := s.Routes().Run()

	if err != nil {
		log.Printf("Server - there was an error calling Run on router: %v", err)
		return err
	}

	return nil
}
