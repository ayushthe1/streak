package routes

import (
	"github.com/ayushthe1/streak/controller"
	"github.com/ayushthe1/streak/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controller.Signup)
	app.Post("/api/login", controller.Login)

	// Authenticated Routes
	app.Use(middleware.IsAuthenticate)

}
