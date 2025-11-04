package main

import (
	"fmt"
	"os"
	"strconv"

	"git.itsigo.dev/istigo/pocketmovie/internal/server"
	"github.com/gofiber/fiber/v3"
	_ "github.com/joho/godotenv/autoload"
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
	initial()

	app := server.New()

	app.RegisterFiberRoutes()

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := app.Listen(fmt.Sprintf(":%d", port), fiber.ListenConfig{
		DisableStartupMessage: true,
	})
	if err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
