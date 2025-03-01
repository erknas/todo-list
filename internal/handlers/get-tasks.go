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

// GetTasks handles the retrieval of tasks.
// @Summary Get all tasks
// @Description Retrieve a list of all tasks.
// @Tags tasks
// @Produce json
// @Success 200 {object} types.Tasks
// @Failure 500 {object} lib.APIError
// @Router /tasks [get]
func GetTasks(c fiber.Ctx, tasksGetter TasksGetter) error {
	tasks, err := tasksGetter.GetTasks(c.Context())
	if err != nil {
		slog.Error("get tasks failed", lib.Err(err))
		return c.Status(http.StatusInternalServerError).JSON(lib.InternalServerError())
	}

	resp := types.Tasks{
		Tasks: tasks,
	}

	return c.JSON(resp)
}
