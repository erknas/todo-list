package storage

import "fmt"

type TaskNotFoundError struct {
	ID int
}

func (e TaskNotFoundError) Error() string {
	return fmt.Sprintf("task does not exist with ID: %d", e.ID)
}
