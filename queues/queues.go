// Copyright (c) 2023 Breu Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package queues

import (
	"context"
	"fmt"
	"sync"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"

	wrk "go.breu.io/durex/workflows"
)

type (
	// Name is the name of the queue.
	Name string

	// Queue defines the queue interface.
	Queue interface {
		// Name gets the name of the queue as string.
		Name() Name

		String() string

		// Prefix gets the prefix of the queue as string.
		Prefix() string

		// WorkflowID sanitzes the workflow ID given the workflows.Options.
		WorkflowID(opts wrk.Options) string

		// ExecuteWorkflow executes a workflow given the context, workflows.Options, workflow function or function name, and
		// optional payload. Lets say, we have a queue called "default", we can either pass in the workflow function or the
		// function name.
		//
		//  q := queues.New(queues.WithName("default"), queues.WithClient(client))
		//  q.ExecuteWorkflow(
		//    ctx,
		//    workflows.NewOptions(
		//      workflows.WithBlock("healthz"),
		//      workflows.WithBlockID(uuid.New().String()),
		//    ),
		//    WorkflowFn, // or "WorkflowFunctionName"
		//    payload...,    // optional.
		//  )
		ExecuteWorkflow(ctx context.Context, opts wrk.Options, fn any, payload ...any) (WorkflowRun, error)

		// ExecuteChildWorkflow executes a child workflow given the parent workflow context, workflows.Options, workflow
		// function or function name and optional payload. It must be executed from within a workflow.
		//
		//  future, err := q.ExecuteChildWorkflow(
		//    ctx,
		//    workflows.NewOptions(
		//      workflows.WithParent(ctx), // This is important. It tells the queue that this is a child workflow.
		//      workflows.WithBlock("healthz"),
		//      workflows.WithBlockID(uuid.New().String()),
		//    ),
		//    WorkflowFn,    // or "WorkflowFunctionName"
		//    payload...,    // optional.
		//  )
		ExecuteChildWorkflow(ctx workflow.Context, opts wrk.Options, fn any, payload ...any) (ChildWorkflowFuture, error)

		// SignalWorkflow signals a workflow given the workflow ID, signal name and optional payload.
		//
		//  if err := q.SignalWorkflow(
		//    ctx,
		//    workflows.NewOptions(
		//      workflows.WithParent(ctx), // This is important. It tells the queue that this is a child workflow.
		//      workflows.WithBlock("healthz"),
		//      workflows.WithBlockID(uuid.New().String()),
		//    ),
		//    "signal-name",
		//    payload,    // or nil
		//  ); err != nil {
		//    // handle error
		//  }
		SignalWorkflow(ctx context.Context, opts wrk.Options, signal Signal, payload any) error

		// SignalWithStartWorkflow signals a workflow given the workflow ID, signal name and optional payload.
		//
		//  run, err := q.SignalWithStartWorkflow(
		//    ctx,
		//    workflows.NewOptions(
		//      workflows.WithParent(ctx), // This is important. It tells the queue that this is a child workflow.
		//      workflows.WithBlock("healthz"),
		//      workflows.WithBlockID(uuid.New().String()),
		//    ),
		//    "signal-name",
		//    arg,    // or nil
		//    WorkflowFn, // or "WorkflowFunctionName"
		//    payload..., // optional.
		//  )
		SignalWithStartWorkflow(ctx context.Context, opts wrk.Options, signal Signal, args any, fn any, payload ...any) (WorkflowRun, error)

		// SignalExternalWorkflow signals a workflow given the workflow ID, signal name and optional payload.
		//
		//  future, err := q.SignalExternalWorkflow(
		//    ctx,
		//    workflows.NewOptions(
		//      workflows.WithParent(ctx), // This is important. It tells the queue that this is a child workflow.
		//      workflows.WithBlock("healthz"),
		//      workflows.WithBlockID(uuid.New().String()),
		//    ),
		//    "signal-name",
		//    payload,    // or nil
		//  )
		SignalExternalWorkflow(ctx workflow.Context, opts wrk.Options, signal Signal, args any) (WorkflowFuture, error)

		// QueryWorkflow queries a workflow given the workflow opts, query name and optional arguments.
		//
		//  result, err := q.QueryWorkflow(
		//    ctx,
		//    workflows.NewOptions(
		//      workflows.WithWorkflowID("my-workflow-id"),
		//    ),
		//    Signal("query-name"),
		//    arg1, arg2, // Optional arguments passed to the query function.
		//  )
		//
		//  if err != nil {
		//    // handle error
		//  }
		//  // Decode the result.
		QueryWorkflow(ctx context.Context, opts wrk.Options, query Query, args ...any) (converter.EncodedValue, error)

		// CreateWorker configures the worker for the queue.
		//
		// This function configures the worker responsible for executing registered workflows and activities.  It uses a
		// builder pattern, with helper functions prefixed by queues.WithWorkerOption{Option}, where {Option} corresponds to
		// a field in Temporal's worker.Option. This allows configuring worker behavior at runtime, such as setting maximum
		// concurrent tasks, enabling sessions etc.  The worker is a singleton, meaning only one worker can be created per
		// queue. Call this function *before* registering workflows and activities (queues.RegisterWorkflow and
		// queues.RegisterActivity) to ensure correct association.
		//
		//	q := queues.New(queues.WithName("my-queue"), queues.WithClient(client))
		//	q.CreateWorker(
		//		queues.WithWorkerOptionEnableSessionWorker(true), // Enable session worker.
		//	)
		CreateWorker(opts ...WorkerOption)

		// Start starts the worker against the queue. CreateWorker must be called before calling this function.
		Start() error

		// Shutdown shuts down the worker against the queue.
		Shutdown(context.Context) error

		// RegisterWorkflow registers a workflow against the queue. It is a wrapper around the worker.RegisterWorkflow.
		RegisterWorkflow(any)

		// RegisterActivity registers an activity against the queue. It is wrapper around the worker.RegisterActivity.
		RegisterActivity(any)
	}

	// QueueOption is the option for a queue.
	QueueOption func(Queue)

	// Queues is a map of queues.
	Queues map[Name]Queue

	// queue defines the basic queue.
	queue struct {
		name                Name   // The name of the temporal queue.
		prefix              string // The prefix for the Workflow ID.
		workflowMaxAttempts int32  // The maximum number of attempts for a workflow.

		client client.Client // The temporal client.

		once   sync.Once     // The sync.Once for the queue.
		worker worker.Worker // The temporal worker.
	}
)

func (q Name) String() string {
	return string(q)
}

func (q *queue) String() string {
	return q.name.String()
}

func (q *queue) Name() Name {
	return q.name
}

func (q *queue) Prefix() string {
	return q.prefix
}

func (q *queue) WorkflowID(opts wrk.Options) string {
	prefix := ""
	if opts.IsChild() {
		prefix = opts.ParentWorkflowID()
	} else {
		prefix = q.Prefix()
	}

	return fmt.Sprintf("%s.%s", prefix, opts.IDSuffix())
}

func (q *queue) ExecuteWorkflow(ctx context.Context, opts wrk.Options, fn any, payload ...any) (WorkflowRun, error) {
	if opts.IsChild() {
		return nil, ErrChildWorkflowExecutionAttempt
	}

	if q.client == nil {
		return nil, ErrClientNil
	}

	return q.client.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:          q.WorkflowID(opts),
			TaskQueue:   q.Name().String(),
			RetryPolicy: q.RetryPolicy(opts),
		},
		fn,
		payload...,
	)
}

func (q *queue) ExecuteChildWorkflow(ctx workflow.Context, opts wrk.Options, fn any, payload ...any) (ChildWorkflowFuture, error) {
	if !opts.IsChild() {
		return nil, wrk.ErrParentNil
	}

	copts := workflow.ChildWorkflowOptions{
		WorkflowID:  q.WorkflowID(opts),
		RetryPolicy: q.RetryPolicy(opts),
	}

	ctx = workflow.WithChildOptions(ctx, copts)

	return workflow.ExecuteChildWorkflow(ctx, fn, payload...), nil
}

// SignalWorkflow signals a workflow given the workflow ID, signal name and optional payload.
func (q *queue) SignalWorkflow(ctx context.Context, opts wrk.Options, signal Signal, args any) error {
	if q.client == nil {
		return ErrClientNil
	}

	if opts.IsChild() {
		return ErrExternalWorkflowSignalAttempt
	}

	return q.client.SignalWorkflow(ctx, q.WorkflowID(opts), "", signal.String(), args)
}

func (q *queue) SignalWithStartWorkflow(ctx context.Context, opts wrk.Options, signal Signal, args any, fn any, payload ...any) (WorkflowRun, error) {
	if q.client == nil {
		return nil, ErrClientNil
	}

	if opts.IsChild() {
		return nil, wrk.ErrParentNil
	}

	return q.client.SignalWithStartWorkflow(
		ctx,
		q.WorkflowID(opts),
		signal.String(),
		args,
		client.StartWorkflowOptions{
			ID:          q.WorkflowID(opts),
			TaskQueue:   q.Name().String(),
			RetryPolicy: q.RetryPolicy(opts),
		},
		fn,
		payload...,
	)
}

func (q *queue) SignalExternalWorkflow(ctx workflow.Context, opts wrk.Options, signal Signal, args any) (WorkflowFuture, error) {
	if !opts.IsChild() {
		return nil, wrk.ErrParentNil
	}

	return workflow.SignalExternalWorkflow(ctx, q.WorkflowID(opts), "", signal.String(), args), nil
}

func (q *queue) QueryWorkflow(
	ctx context.Context, opts wrk.Options, query Query, args ...any,
) (converter.EncodedValue, error) {
	if q.client == nil {
		return nil, ErrClientNil
	}

	return q.client.QueryWorkflow(ctx, q.WorkflowID(opts), "", query.String(), args...)
}

func (q *queue) RetryPolicy(opts wrk.Options) *temporal.RetryPolicy {
	attempts := opts.MaxAttempts()
	if attempts < wrk.RetryForever &&
		q.workflowMaxAttempts < wrk.RetryForever &&
		q.workflowMaxAttempts > attempts {
		attempts = q.workflowMaxAttempts
	}

	return &temporal.RetryPolicy{MaximumAttempts: attempts, NonRetryableErrorTypes: opts.IgnoredErrors()}
}

func (q *queue) CreateWorker(opts ...WorkerOption) {
	q.once.Do(func() {
		opts := NewWorkerOptions(opts...)

		q.worker = worker.New(q.client, q.Name().String(), opts)
	})
}

func (q *queue) Start() error {
	if q.worker == nil {
		return ErrWorkerNil
	}

	return q.worker.Start()
}

func (q *queue) Shutdown(ctx context.Context) error {
	if q.worker == nil {
		return ErrWorkerNil
	}

	q.worker.Stop()

	return nil
}

func (q *queue) RegisterWorkflow(fn any) {
	q.worker.RegisterWorkflow(fn)
}

func (q *queue) RegisterActivity(fn any) {
	q.worker.RegisterActivity(fn)
}

// WithName sets the queue name and the prefix for the workflow ID.
func WithName(name string) QueueOption {
	return func(q Queue) {
		q.(*queue).name = Name(name)
		q.(*queue).prefix = DefaultPrefix() + name
	}
}

// WithWorkflowMaxAttempts sets the maximum number of attempts for all the workflows in the queue.
// The default value is 0 i.e. RetryForever.
func WithWorkflowMaxAttempts(attempts int32) QueueOption {
	return func(q Queue) {
		q.(*queue).workflowMaxAttempts = attempts
	}
}

// WithClient sets the client for the queue.
func WithClient(c client.Client) QueueOption {
	return func(q Queue) {
		q.(*queue).client = c
	}
}

// New creates a new queue with the given opts.
// For a queue named "default", we will defined it as follows:
//
//	var DefaultQueue = queue.New(
//	  queue.WithName("default"),
//	  queue.WithClient(client),
//	  queue.WithMaxWorkflowAttempts(1),
//	)
func New(opts ...QueueOption) Queue {
	q := &queue{workflowMaxAttempts: wrk.RetryForever}
	for _, opt := range opts {
		opt(q)
	}

	return q
}
