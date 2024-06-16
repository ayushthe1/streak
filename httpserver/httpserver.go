package httpserver

import (
	"log"

	"github.com/ayushthe1/streak/handler"
	"github.com/ayushthe1/streak/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func setupRoutes(app *fiber.App) {
	public := app.Group("/api")
	public.Post("/register", handler.SignupHandler)
	public.Post("/login", handler.LoginHandler)
	public.Post("/logout", handler.LogoutHandler)
	public.Post("/verify-contact", handler.VerifyContactHandler)

	protected := app.Group("/api", middleware.IsAuthenticate)
	protected.Get("/chat-history", handler.ChatHistoryHandler)
	protected.Get("/contact-list", handler.ContactHandler)
	protected.Post("file-upload", handler.FileUploadHandler)

}

func StartHttpServer() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:4000 ",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	port := "3000"
	setupRoutes(app)
	log.Println("Starting HTTP Server on port", port)
	app.Listen(":" + port)
}
