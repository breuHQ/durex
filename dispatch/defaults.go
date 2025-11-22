// Package dispatch provides helper functions for configuring workflow.Context objects with default activity options.
package dispatch

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/durex/queues"
)

// WithDefaultActivityContext returns a workflow.Context with the default activity options applied.
// The default options include a StartToCloseTimeout of 60 seconds.
//
// Example:
//
//	ctx = shared.WithDefaultActivityContext(ctx)
func WithDefaultActivityContext(ctx workflow.Context) workflow.Context {
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
	})
}

// WithIgnoredErrorsContext returns a workflow.Context with activity options configured with a
// StartToCloseTimeout of 60 seconds and a RetryPolicy that allows a single attempt and ignores
// specified error types.
//
// Example:
//
//	ignored := []string{"CustomErrorType"}
//	ctx = shared.WithIgnoredErrorsContext(ctx, ignored...)
func WithIgnoredErrorsContext(ctx workflow.Context, args ...string) workflow.Context {
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:        1,
			NonRetryableErrorTypes: args,
		},
	})
}

// WithMarathonContext returns a workflow.Context with activity options configured for long-running activities.
// It sets the StartToCloseTimeout to 60 minutes and the HeartbeatTimeout to 30 seconds.
//
// Example:
//
//	ctx = shared.WithMarathonContext(ctx)
func WithMarathonContext(ctx workflow.Context) workflow.Context {
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 24 * time.Hour,
		HeartbeatTimeout:    30 * time.Second,
	})
}

// WithCustomQueueContext returns a workflow.Context with activity options configured with a
// StartToCloseTimeout of 60 seconds and a dedicated task queue. This allows scheduling activities
// on a different queue than the one the workflow is running on.
//
// Example:
//
//	ctx = shared.WithCustomQueueContext(ctx, queues.MyTaskQueue)
func WithCustomQueueContext(ctx workflow.Context, q queues.Queue) workflow.Context {
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		TaskQueue:           q.String(),
	})
}
