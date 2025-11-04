package server

import (
	"git.itsigo.dev/istigo/pocketmovie/internal/database"
	"github.com/gofiber/fiber/v3"
)

type FiberServer struct {
	*fiber.App

	db *database.Queries
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "PocketMovie",
			AppName:      "PocketMovie",
		}),
		db: database.Init(),
	}
	return server
}
