package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/erknas/todo-list/internal/lib"
	"github.com/erknas/todo-list/internal/storage"
	"github.com/gofiber/fiber/v3"
)

type TaskDeleter interface {
	DeleteTask(context.Context, int) error
}

func DeleteTask(c fiber.Ctx, taskDeleter TaskDeleter) error {
	id, err := lib.ParseID(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"statusCode": http.StatusBadRequest,
			"error":      err.Error(),
		})
	}

	if err := taskDeleter.DeleteTask(c.Context(), id); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"statusCode": http.StatusBadRequest,
				"error":      err.Error(),
			})
		}

		slog.Error("delete task failed", lib.Err(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"statusCode": http.StatusInternalServerError,
			"error":      "internal server error",
		})
	}

	resp := lib.TaskResponse("task successfully deleted", id)

	return c.JSON(resp)
}
