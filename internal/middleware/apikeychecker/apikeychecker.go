package apikeychecker

import (
	"strings"

	"git.itsigo.dev/istigo/pocketmovie/internal/database"
	"github.com/gofiber/fiber/v3"
)

type Config struct {
	DB database.Queries
}

func configDefault(config ...Config) Config {
	cfg := config[0]
	return cfg
}

func New(config Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		setting, _ := config.DB.ListSetting(c, 2)
		referer := string(c.Request().Header.Referer())
		if setting.Value == "" && !strings.Contains(referer, "apikey") {
			return c.Redirect().Status(fiber.StatusMovedPermanently).To("/apikey")
		}
		return c.Next()
	}
}
