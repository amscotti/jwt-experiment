package todo

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

type Todo struct {
	Id          int
	UserId      float64
	Description string
	Completed   bool
	CreatedOn   time.Time
}

// Simple in memory database for quick experimentation
var todos = []Todo{
	{Id: 0, UserId: 123, Description: "Pick up milk", Completed: false, CreatedOn: time.Now()},
	{Id: 1, UserId: 123, Description: "Learn Go", Completed: false, CreatedOn: time.Now()},
	{Id: 2, UserId: 123, Description: "Write new blog posting", Completed: false, CreatedOn: time.Now()},
	{Id: 3, UserId: 321, Description: "Someone elses todo item", Completed: false, CreatedOn: time.Now()},
	{Id: 4, UserId: 4321, Description: "Someone elses todo item", Completed: false, CreatedOn: time.Now()},
}

func getUserId(c *fiber.Ctx) float64 {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	return userId
}

func getTodos(c *fiber.Ctx) error {
	userId := getUserId(c)

	var usersTodos []Todo
	for _, t := range todos {
		if t.UserId == userId {
			usersTodos = append(usersTodos, t)
		}
	}

	return c.JSON(usersTodos)
}

func addTodo(c *fiber.Ctx) error {
	userId := getUserId(c)

	todo := new(Todo)
	if err := c.BodyParser(todo); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	lastTodo := todos[len(todos)-1]

	todo.Id = lastTodo.Id + 1
	todo.UserId = userId
	todo.Completed = false
	todo.CreatedOn = time.Now()

	todos = append(todos, *todo)

	return c.JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	userId := getUserId(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	todo := new(Todo)
	if err := c.BodyParser(todo); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	successful := false

	for i, t := range todos {
		if t.UserId == userId && t.Id == id {
			t.Completed = todo.Completed
			t.Description = todo.Description
			todos[i] = t

			successful = true
			break
		}
	}

	return c.JSON(fiber.Map{"successful": successful})
}

func removeTodo(c *fiber.Ctx) error {
	userId := getUserId(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	successful := false

	for i, t := range todos {
		if t.UserId == userId && t.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			successful = true
			break
		}
	}

	return c.JSON(fiber.Map{"successful": successful})
}

func GracefulShutdown() error {
	return nil
}

func RegisterHandlers(app *fiber.App, signingKey string) {
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(signingKey),
	}))

	app.Get("/todo", getTodos)
	app.Post("/todo", addTodo)
	app.Put("/todo/:id", updateTodo)
	app.Delete("/todo/:id", removeTodo)
}
