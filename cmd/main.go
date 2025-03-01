package main

import (
	"context"
	"log"

	"github.com/erknas/todo-list/internal/config"
	"github.com/erknas/todo-list/internal/handlers"
	"github.com/erknas/todo-list/internal/storage"
	"github.com/erknas/todo-list/migrations/migrator"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

// @title TODO-list API
// @version 1.0
// @description This is a simple API for managing tasks.
// @host localhost:3000
// @BasePath /tasks
func main() {
	var (
		ctx = context.Background()
		cfg = config.Load()
	)

	if err := migrator.New(cfg); err != nil {
		log.Fatalf("failed to migrate: %s", err)
	}

	storage, err := storage.NewPostgresPool(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %s", err)
	}
	defer storage.Close()

	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/tasks")

	api.Post("/", func(c fiber.Ctx) error {
		return handlers.CreateTask(c, storage)
	})

	api.Get("/", func(c fiber.Ctx) error {
		return handlers.GetTasks(c, storage)
	})

	api.Put("/:id", func(c fiber.Ctx) error {
		return handlers.UpdateTask(c, storage)
	})

	api.Delete("/:id", func(c fiber.Ctx) error {
		return handlers.DeleteTask(c, storage)
	})

	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Page not found")
	})

	log.Fatal(app.Listen(":" + cfg.Port))
}
