package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/erknas/todo-list/internal/lib"
	"github.com/erknas/todo-list/internal/storage"
	"github.com/erknas/todo-list/internal/types"
	"github.com/gofiber/fiber/v3"
)

type TaskUpdater interface {
	UpdateTask(context.Context, int, types.TaskRequest) error
}

func UpdateTask(c fiber.Ctx, taskUpdater TaskUpdater) error {
	var req types.TaskRequest

	if err := c.Bind().Body(&req); err != nil {
		slog.Error("failed to decode request body", lib.Err(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"statusCode": http.StatusBadRequest,
			"error":      "invalid request data",
		})
	}

	if errors := req.ValidateUpdateTaskRequest(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}

	id, err := lib.ParseID(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"statusCode": http.StatusBadRequest,
			"error":      err.Error(),
		})
	}

	if err := taskUpdater.UpdateTask(c.Context(), id, req); err != nil {
		if errors.Is(err, storage.ErrNotFound) || errors.Is(err, storage.ErrNoUpdate) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"statusCode": http.StatusBadRequest,
				"msg":        err.Error(),
			})
		}

		slog.Error("update task failed", lib.Err(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"statusCode": http.StatusInternalServerError,
			"error":      "internal server error",
		})
	}

	resp := lib.TaskResponse("task successfully updated", id)

	return c.JSON(resp)
}
