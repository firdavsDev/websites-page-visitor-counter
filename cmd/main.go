package main

import (
	"context"
	"log"
	"visitor-counter/config"
	"visitor-counter/internal/handlers"
	"visitor-counter/internal/middleware"
	"visitor-counter/internal/storage"

	"github.com/gofiber/fiber/v2/middleware/cors"

	_ "visitor-counter/docs"

	"github.com/gofiber/swagger"

	"github.com/gofiber/fiber/v2"
)

// @title Visitor Counter API
// @version 1.0
// @description This is a simple API to track website visitors.
// @host localhost:8080
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Token for authentication
// @security ApiKeyAuth
// @tag.name Visitor
// @tag.description API for tracking visitors
func main() {
	cfg := config.LoadConfig()

	store, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Conn.Close(context.Background())

	app := fiber.New(
		fiber.Config{
			StrictRouting: true,
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				if e, ok := err.(*fiber.Error); ok {
					return c.Status(e.Code).JSON(fiber.Map{
						"error": e.Message,
					})
				}
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal Server Error",
				})
			},
			DisableStartupMessage: true,
			// Increase header buffer size (default is 4096 bytes)
			ReadBufferSize: 8192, // or higher if needed
		},
	)

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080, http://127.0.0.1:8080",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		// MaxAge:           86400,

	}))

	// Middleware
	// app.Use(middleware.Authenticate(store))

	// Handlers
	h := &handlers.Handler{Store: store}
	// Apply authentication middleware only to protected routes
	api := app.Group("/api", middleware.Authenticate(store))
	api.Get("/track", h.TrackVisitor)

	// Swagger
	// app.Get("/swagger/*", swagger.HandlerDefault) // default
	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://localhost:8080/swagger/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
	}))

	log.Fatal(app.Listen(":" + cfg.ServerPort))
}
