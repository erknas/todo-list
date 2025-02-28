package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/erknas/todo-list/internal/lib"
	"github.com/erknas/todo-list/internal/types"
	"github.com/gofiber/fiber/v3"
)

type TaskCreator interface {
	CreateTask(context.Context, types.TaskRequest) (int, error)
}

func CreateTask(c fiber.Ctx, taskCreator TaskCreator) error {
	var req types.TaskRequest

	if err := c.Bind().Body(&req); err != nil {
		slog.Error("failed to decode request body", lib.Err(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"statusCode": http.StatusBadRequest,
			"error":      "invalid request data",
		})
	}

	if errors := req.ValidateTaskRequest(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}

	id, err := taskCreator.CreateTask(c.Context(), req)
	if err != nil {
		slog.Error("create task failed", lib.Err(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"statusCode": http.StatusInternalServerError,
			"error":      "internal server error",
		})
	}

	resp := lib.TaskResponse("task successfully created", id)

	return c.JSON(resp)
}
