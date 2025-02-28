package types

import "github.com/gofiber/fiber/v3"

func (t TaskRequest) ValidateCreateTaskRequest() fiber.Map {
	errors := make(fiber.Map)

	if t.Status != "new" && t.Status != "in_progress" && t.Status != "done" {
		errors["status"] = "wrong status"
	}

	return errors
}

func (t TaskRequest) ValidateUpdateTaskRequest() fiber.Map {
	errors := make(fiber.Map)

	if len(t.Status) == 0 {
		return nil
	}

	if t.Status != "new" && t.Status != "in_progress" && t.Status != "done" {
		errors["status"] = "wrong status"
	}

	return errors
}
