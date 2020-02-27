package rocketoff

import (
	"fmt"
)

type ErrInvalidType struct {
	expected string
}

func (e ErrInvalidType) Error() string {
	return fmt.Sprintf("invalid data type received, expected %s", e.expected)
}
