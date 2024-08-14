package queues_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/temporal-tools/queues"
	"go.breu.io/temporal-tools/workflows"
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
	c := &MockClient{env: s.env}

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
	_, _ = s.queue.ExecuteWorkflow(ctx, opts, fn, "hello")

	expected := "hello world"
	result := ""

	_ = s.env.GetWorkflowResult(&result)

	s.Equal(expected, result)
}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}
