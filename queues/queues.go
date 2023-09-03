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

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/temporal-tools/workflows"
)

type (
	// Name is the name of the queue.
	Name string

	// Queue defines the queue interface.
	Queue interface {
		// Name gets the name of the queue as string.
		Name() Name

		// Prefix gets the prefix of the queue as string.
		Prefix() string

		// WorkflowID sanitzes the workflow ID given the workflows.Options.
		WorkflowID(options workflows.Options) string

		// ExecuteWorkflow executes a workflow given the context, workflows.Options, workflow function or function name, and
		// optional payload.
		// Lets say, we have a queue called "default", we can either pass in the workflow function or the function name.
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
		ExecuteWorkflow(ctx context.Context, options workflows.Options, fn any, payload ...any) (client.WorkflowRun, error)

		// ExecuteChildWorkflow executes a child workflow given the parent workflow context, workflows.Options,
		// workflow function or function name and optional payload. It must be executed from within a workflow.
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
		ExecuteChildWorkflow(ctx workflow.Context, options workflows.Options, fn any, payload ...any) (ChildWorkflowFuture, error)

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
		SignalWorkflow(ctx context.Context, options workflows.Options, signalName string, payload any) error

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
		SignalExternalWorkflow(ctx workflow.Context, options workflows.Options, signal string, payload any) (workflow.Future, error)

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
		SignalWithStartWorkflow(
			ctx context.Context, options workflows.Options, signal string, arg any, fn any, payload ...any,
		) (client.WorkflowRun, error)

		// CreateWorker creates a worker against the queue.
		CreateWorker() worker.Worker
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
	}

	ChildWorkflowFuture workflow.ChildWorkflowFuture
)

func (q Name) String() string {
	return string(q)
}

func (q *queue) Name() Name {
	return q.name
}

func (q *queue) Prefix() string {
	return q.prefix
}

func (q *queue) WorkflowID(options workflows.Options) string {
	pfix := ""
	if options.IsChild() {
		pfix, _ = options.ParentWorkflowID()
	} else {
		pfix = q.Prefix()
	}

	return fmt.Sprintf("%s.%s", pfix, options.IDSuffix())
}

func (q *queue) ExecuteWorkflow(ctx context.Context, opts workflows.Options, fn any, payload ...any) (client.WorkflowRun, error) {
	attempts := opts.MaxAttempts()
	if attempts != workflows.RetryForever && q.workflowMaxAttempts != workflows.RetryForever && q.workflowMaxAttempts > attempts {
		attempts = q.workflowMaxAttempts
	}

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
			RetryPolicy: &temporal.RetryPolicy{MaximumAttempts: attempts},
		},
		fn,
		payload...,
	)
}

func (q *queue) ExecuteChildWorkflow(ctx workflow.Context, opts workflows.Options, fn any, payload ...any) (ChildWorkflowFuture, error) {
	attempts := opts.MaxAttempts()

	if !opts.IsChild() {
		return nil, workflows.ErrParentNil
	}

	if attempts != workflows.RetryForever && q.workflowMaxAttempts != workflows.RetryForever && q.workflowMaxAttempts > attempts {
		attempts = q.workflowMaxAttempts
	}

	copts := workflow.ChildWorkflowOptions{
		WorkflowID:  q.WorkflowID(opts),
		RetryPolicy: &temporal.RetryPolicy{MaximumAttempts: attempts},
	}

	ctx = workflow.WithChildOptions(ctx, copts)

	return workflow.ExecuteChildWorkflow(ctx, fn, payload...), nil
}

// SignalWorkflow signals a workflow given the workflow ID, signal name and optional payload.
func (q *queue) SignalWorkflow(ctx context.Context, opts workflows.Options, signalName string, arg any) error {
	if q.client == nil {
		return ErrClientNil
	}

	if opts.IsChild() {
		return workflows.ErrParentNil
	}

	return q.client.SignalWorkflow(ctx, q.WorkflowID(opts), "", signalName, arg)
}

func (q *queue) SignalExternalWorkflow(ctx workflow.Context, opts workflows.Options, signalName string, arg any) (workflow.Future, error) {
	if !opts.IsChild() {
		return nil, ErrChildWorkflowExecutionAttempt
	}

	return workflow.SignalExternalWorkflow(ctx, q.WorkflowID(opts), "", signalName, arg), nil
}

func (q *queue) SignalWithStartWorkflow(
	ctx context.Context, opts workflows.Options, sig string, arg any, fn any, payload ...any,
) (client.WorkflowRun, error) {
	if q.client == nil {
		return nil, ErrClientNil
	}

	if opts.IsChild() {
		return nil, workflows.ErrParentNil
	}

	attempts := opts.MaxAttempts()
	if attempts != workflows.RetryForever && q.workflowMaxAttempts != workflows.RetryForever && q.workflowMaxAttempts > attempts {
		attempts = q.workflowMaxAttempts
	}

	return q.client.SignalWithStartWorkflow(
		ctx,
		q.WorkflowID(opts),
		sig,
		arg,
		client.StartWorkflowOptions{
			ID:          q.WorkflowID(opts),
			TaskQueue:   q.Name().String(),
			RetryPolicy: &temporal.RetryPolicy{MaximumAttempts: attempts},
		},
		fn,
		payload...,
	)
}

func (q *queue) CreateWorker() worker.Worker {
	options := worker.Options{OnFatalError: func(err error) { panic(err) }}
	return worker.New(q.client, q.Name().String(), options)
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

// New creates a new queue with the given options.
// For a queue named "default", we will defined it as follows:
//
//	var DefaultQueue = queue.New(
//	  queue.WithName("default"),
//	  queue.WithClient(client),
//	  queue.WithMaxWorkflowAttempts(1),
//	)
func New(opts ...QueueOption) Queue {
	q := &queue{workflowMaxAttempts: workflows.RetryForever}
	for _, opt := range opts {
		opt(q)
	}

	return q
}
