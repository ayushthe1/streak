package httpserver

import (
	"log"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/ayushthe1/streak/chatbot"
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
	public.Post("/wv", handler.GetQueryDataFromWeaviate)
	public.Post("/hook", chatbot.WebhookHandler)
	// public.Post("/chatbot", chatbot.ChatbotHandler)

	protected := app.Group("/api", middleware.IsAuthenticate)
	protected.Get("/chat-history", handler.ChatHistoryHandler)
	protected.Get("/contact-list", handler.ContactHandler)
	protected.Get("/activities", handler.ActivityHandler)
	protected.Post("upload", handler.FileUploadHandler)

}

func StartHttpServer() {
	app := fiber.New()

	// Setup Prometheus for metrics collection
	prometheus := fiberprometheus.New("streak-service")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://frontend:4000, https://streak.ayushsharma.co.in, http://host.docker.internal:4000 ",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	port := "3000"
	setupRoutes(app)
	log.Println("Starting HTTP Server on port", port)
	app.Listen(":" + port)
}
