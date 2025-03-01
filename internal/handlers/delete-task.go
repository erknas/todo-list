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

type TaskDeleter interface {
	DeleteTask(context.Context, int) error
}

// DeleteTask handles the deletion of a task by ID.
// @Summary Delete a task
// @Description Delete a task by its ID.
// @Tags tasks
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} types.TaskResp
// @Failure 400 {object} lib.APIError
// @Failure 500 {object} lib.APIError
// @Router /tasks/:id [delete]
func DeleteTask(c fiber.Ctx, taskDeleter TaskDeleter) error {
	id, err := lib.ParseID(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(lib.InvalidID())
	}

	if err := taskDeleter.DeleteTask(c.Context(), id); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.Status(http.StatusBadRequest).JSON(lib.TaskNotFound(id))
		}

		slog.Error("delete task failed", lib.Err(err))
		return c.Status(http.StatusInternalServerError).JSON(lib.InternalServerError())
	}

	resp := types.TaskResp{
		TaskID: id,
		Msg:    "task successfully deleted",
	}

	return c.JSON(resp)
}
