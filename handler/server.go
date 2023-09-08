package handler

import (
	"github.com/AthanatiusC/SawitPro/repository"
)

type Server struct {
	Repository repository.RepositoryInterface
	JWTSecret  string
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	Secret     string
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		JWTSecret:  opts.Secret,
	}
}
