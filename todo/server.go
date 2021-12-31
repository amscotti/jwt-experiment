package todo

import (
	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func isVaild(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)

	return c.JSON(fiber.Map{"vaild": true, "name": name})
}

func GracefulShutdown() error {
	return nil
}

func RegisterHandlers(app *fiber.App, signingKey string) {
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(signingKey),
	}))

	app.Get("/is-vaild", isVaild)
}
