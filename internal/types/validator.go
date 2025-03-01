package types

func (req NewTaskRequest) ValidateCreateTaskRequest() map[string]string {
	errors := make(map[string]string)

	if req.Status != "new" && req.Status != "in_progress" && req.Status != "done" {
		errors["status"] = "wrong status"
	}

	return errors
}

func (req NewTaskRequest) ValidateUpdateTaskRequest() map[string]string {
	errors := make(map[string]string)

	if len(req.Status) == 0 {
		return nil
	}

	if req.Status != "new" && req.Status != "in_progress" && req.Status != "done" {
		errors["status"] = "wrong status"
	}

	return errors
}
