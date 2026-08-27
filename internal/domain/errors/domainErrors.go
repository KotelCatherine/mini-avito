package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrAlradyExists = errors.New("alredy exists")
	ErrInvalidInput = errors.New("ivalid input")

	ErrAdNotFound      = fmt.Errorf("ad: %w", ErrNotFound)
	ErrAdInvalidPrice  = fmt.Errorf("ad price: %w", ErrInvalidInput)
	ErrInvalidTitle    = fmt.Errorf("ad title: %w", ErrInvalidInput)
	ErrInvalidCategory = fmt.Errorf("ad category: %w", ErrInvalidInput)
)
