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
	CreateTask(context.Context, types.NewTaskRequest) (types.NewTaskResponse, error)
}

// CreateTask handles the creation of a new task.
// @Summary Create a task
// @Description Create a new task with the provided details.
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body types.NewTaskRequest true "Task data"
// @Success 200 {object} types.NewTaskResponse
// @Failure 400 {object} lib.APIError
// @Failure 500 {object} lib.APIError
// @Router /tasks [post]
func CreateTask(c fiber.Ctx, taskCreator TaskCreator) error {
	var req types.NewTaskRequest

	if err := c.Bind().Body(&req); err != nil {
		slog.Error("failed to decode request body", lib.Err(err))
		return c.Status(http.StatusBadRequest).JSON(lib.InvalidJSON())
	}

	if errors := req.ValidateCreateTaskRequest(); len(errors) > 0 {
		return c.Status(http.StatusUnprocessableEntity).JSON(lib.InvalidRequestData(errors))
	}

	resp, err := taskCreator.CreateTask(c.Context(), req)
	if err != nil {
		slog.Error("create task failed", lib.Err(err))
		return c.Status(http.StatusInternalServerError).JSON(lib.InternalServerError())
	}

	return c.JSON(resp)
}
