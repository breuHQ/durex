package queues

import (
	"errors"
)

var (
	// ErrClientNil is returned when the client is nil.
	ErrClientNil                     = errors.New("client is nil")
	ErrChildWorkflowExecutionAttempt = errors.New("attempting to execute child workflow directly. use ExecuteWorkflow instead")
)
