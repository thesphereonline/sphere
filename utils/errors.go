package utils

import "fmt"

type BlockchainError struct {
	Message string
	Err     error
}

func (e *BlockchainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewError(message string, err error) error {
	return &BlockchainError{
		Message: message,
		Err:     err,
	}
}
