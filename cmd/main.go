package main

import (
	"os"

	"github.com/AthanatiusC/SawitPro/generated"
	"github.com/AthanatiusC/SawitPro/handler"
	"github.com/AthanatiusC/SawitPro/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	secret := os.Getenv("SECRET")

	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})

	opts := handler.NewServerOptions{
		Repository: repo,
		Secret:     secret,
	}

	return handler.NewServer(opts)
}
