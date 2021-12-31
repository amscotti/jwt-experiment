package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func login(signingKey string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.FormValue("user")
		pass := c.FormValue("pass")

		if user != "john" || pass != "doe" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Create the Claims
		claims := jwt.MapClaims{
			"name":    "John Doe",
			"user_id": 123,
			"role_id": 1,
			"exp":     time.Now().Add(time.Hour * 72).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte(signingKey))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": t})
	}
}

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
	app.Post("/login", login(signingKey))

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(signingKey),
	}))

	app.Get("/is-vaild", isVaild)
}
