package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/erknas/todo-list/internal/lib"
	"github.com/erknas/todo-list/internal/types"
	"github.com/gofiber/fiber/v3"
)

type TasksGetter interface {
	GetTasks(context.Context) ([]types.Task, error)
}

func GetTasks(c fiber.Ctx, tasksGetter TasksGetter) error {
	tasks, err := tasksGetter.GetTasks(c.Context())
	if err != nil {
		slog.Error("get tasks failed", lib.Err(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"statusCode": http.StatusInternalServerError,
			"error":      "internal server error",
		})
	}

	resp := types.Tasks{
		Tasks: tasks,
	}

	return c.JSON(resp)
}
