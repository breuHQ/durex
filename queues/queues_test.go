package queues_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/durex/queues"
	"go.breu.io/durex/workflows"
)

type (
	QueueTestSuite struct {
		suite.Suite
		testsuite.WorkflowTestSuite

		env   *testsuite.TestWorkflowEnvironment
		queue queues.Queue
	}
)

func (s *QueueTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
	c := func() client.Client {
		return &MockClient{env: s.env}
	}

	s.queue = queues.New(
		queues.WithName("test"),
		queues.WithClient(c),
	)
}

func (s *QueueTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *QueueTestSuite) TestName() {
	expected := "test"

	s.Equal(expected, s.queue.Name().String())
}

func (s *QueueTestSuite) TestPrefix() {
	expected := "io.breu.test"

	s.Equal(expected, s.queue.Prefix())
}

func (s *QueueTestSuite) TestWorkflowID() {
	id := uuid.New()
	opts, _ := workflows.NewOptions(
		workflows.WithBlock("suite"),
		workflows.WithBlockID(id.String()),
	)
	expected := "io.breu.test.suite." + id.String()
	actual := s.queue.WorkflowID(opts)

	s.Equal(expected, actual)
}

func (s *QueueTestSuite) TestExecuteWorkflow() {
	ctx := context.Background()
	id := uuid.New()
	opts, _ := workflows.NewOptions(
		workflows.WithBlock("test"),
		workflows.WithBlockID(id.String()),
	)

	fn := func(ctx workflow.Context, in string) (string, error) {
		return in + " world", nil
	}

	// Execute the workflow
	_, err := s.queue.ExecuteWorkflow(ctx, opts, fn, "hello")

	s.NoError(err)

	expected := "hello world"
	result := ""

	_ = s.env.GetWorkflowResult(&result)

	s.Equal(expected, result)
}

func (s *QueueTestSuite) TestExecuteChildWorkflow() {
	ctx := context.Background()
	parent_id := uuid.New()
	child_id := uuid.New()

	opts, _ := workflows.NewOptions(
		workflows.WithBlock("parent"),
		workflows.WithBlockID(parent_id.String()),
	)

	child := func(ctx workflow.Context, payload string) (string, error) { return payload + "child", nil }
	parent := func(ctx workflow.Context, payload string) (string, error) {
		opts, _ := workflows.NewOptions(
			workflows.WithParent(ctx),
			workflows.WithBlock(child_id.String()),
		)

		result := ""
		f, _ := s.queue.ExecuteChildWorkflow(ctx, opts, child, payload)

		_ = f.Get(ctx, &result)

		return result, nil
	}

	s.env.RegisterWorkflow(child)
	_, err := s.queue.ExecuteWorkflow(ctx, opts, parent, "parent/")

	s.NoError(err)

	expected := "parent/child"
	result := ""

	_ = s.env.GetWorkflowResult(&result)

	s.Equal(expected, result)
}

func (s *QueueTestSuite) TestSignalWorkflow() {
	ctx := context.Background()
	id := uuid.New()
	name := queues.Signal("signal")
	opts, _ := workflows.NewOptions(
		workflows.WithBlock("signal"),
		workflows.WithBlockID(id.String()),
	)

	fn := func(ctx workflow.Context, payload string) (string, error) {
		selector := workflow.NewSelector(ctx)
		signal := workflow.GetSignalChannel(ctx, name.String())
		result := ""

		selector.AddReceive(signal, func(channel workflow.ReceiveChannel, more bool) {
			channel.Receive(ctx, &result)

			result = payload + result
		})

		selector.Select(ctx)

		return result, nil
	}

	s.env.RegisterDelayedCallback(func() {
		_ = s.queue.SignalWorkflow(ctx, opts, name, "world")
	}, 30*time.Second)

	_, err := s.queue.ExecuteWorkflow(ctx, opts, fn, "hello ")

	s.NoError(err)

	expected := "hello world"
	result := ""

	_ = s.env.GetWorkflowResult(&result)

	s.Equal(expected, result)
}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}
