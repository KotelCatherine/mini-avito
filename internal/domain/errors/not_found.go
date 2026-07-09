package errors

import "fmt"

type NotFoundError struct {
	EntityType string
	ID         string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id '%s' not found", e.EntityType, e.ID)
}

func NewNotFoundError(entityType, id string) *NotFoundError {
	return &NotFoundError{
		EntityType: entityType,
		ID:         id,
	}
}

func (e *NotFoundError) Unwrap() error {
	return ErrNotFound
}

func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}
