package queues

import (
	"errors"
)

var (
	// ErrClientNil is returned when the temporal client is nil.
	ErrClientNil = errors.New("client is nil")
	// ErrChildWorkflowExecutionAttempt is returned when attempting to execute a child workflow without the parent.
	ErrChildWorkflowExecutionAttempt = errors.New("attempting to execute child workflow directly. use ExecuteWorkflow instead")
)
