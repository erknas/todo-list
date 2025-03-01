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
	UpdateTask(context.Context, int, types.NewTaskRequest) error
}

// UpdateTask handles the update of a task by ID.
// @Summary Update a task
// @Description Update a task by its ID with the provided details.
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body types.NewTaskRequest true "Task data"
// @Success 200 {object} types.TaskResp
// @Failure 400 {object} lib.APIError
// @Failure 500 {object} lib.APIError
// @Router /tasks/:id [put]
func UpdateTask(c fiber.Ctx, taskUpdater TaskUpdater) error {
	var req types.NewTaskRequest

	if err := c.Bind().Body(&req); err != nil {
		slog.Error("failed to decode request body", lib.Err(err))
		return c.Status(http.StatusBadRequest).JSON(lib.InvalidJSON())
	}

	if errors := req.ValidateUpdateTaskRequest(); len(errors) > 0 {
		return c.Status(http.StatusUnprocessableEntity).JSON(lib.InvalidRequestData(errors))
	}

	id, err := lib.ParseID(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(lib.InvalidID())
	}

	if err := taskUpdater.UpdateTask(c.Context(), id, req); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.Status(http.StatusBadRequest).JSON(lib.TaskNotFound(id))
		}
		if errors.Is(err, storage.ErrNoUpdate) {
			return c.Status(http.StatusBadRequest).JSON(lib.NothigToUpdate())
		}

		slog.Error("update task failed", lib.Err(err))
		return c.Status(http.StatusInternalServerError).JSON(lib.InternalServerError())
	}

	resp := types.TaskResp{
		TaskID: id,
		Msg:    "task successfully updated",
	}

	return c.JSON(resp)
}
