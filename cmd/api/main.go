package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"git.itsigo.dev/istigo/pocketmovie/internal/apis"
	"git.itsigo.dev/istigo/pocketmovie/internal/database"
	"git.itsigo.dev/istigo/pocketmovie/internal/server"
	"github.com/gofiber/fiber/v3"
)

func initial() {
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		os.Mkdir("./data/", 0755)
		fmt.Println("Data created")
		os.Mkdir("./data/img/", 0755)
		fmt.Println("Img created")
	}
}

func main() {
	var port int

	flag.IntVar(&port, "p", 3000, "Provide a port number")

	flag.Parse()

	initial()

	db := database.Init()
	app := server.New(db)

	ctx := context.Background()

	// api key to request manager
	token, _ := db.ListSetting(ctx, 2)
	apis.Init(token.Value)

	app.RegisterFiberRoutes()

	err := app.Listen(fmt.Sprint(":", port), fiber.ListenConfig{
		DisableStartupMessage: false,
	})
	if err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
