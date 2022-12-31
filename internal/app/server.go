package app

import (
	"log"
	"speakeasy/internal/pkg/authentication"
	"speakeasy/internal/pkg/profile"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router                *gin.Engine
	profileService        profile.ProfileService
	authenticationService authentication.AuthenticationService
}

func NewServer(
	router *gin.Engine,
	profileService profile.ProfileService,
	authenticationService authentication.AuthenticationService,
) *Server {
	return &Server{
		router:                router,
		profileService:        profileService,
		authenticationService: authenticationService,
	}
}

func (s *Server) Run() error {
	r := s.Routes()
	err := r.Run()

	if err != nil {
		log.Printf("Server - there was an error calling Run on router: %v", err)
		return err
	}

	return nil
}
