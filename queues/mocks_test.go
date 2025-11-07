// nolint
package queues_test

import (
	"context"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/operatorservice/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/testsuite"
)

type (
	MockClient struct {
		env *testsuite.TestWorkflowEnvironment
	}
)

func NewMockClient(env *testsuite.TestWorkflowEnvironment) *MockClient {
	return &MockClient{env: env}
}

func (m *MockClient) ExecuteWorkflow(
	ctx context.Context, options client.StartWorkflowOptions, workflow any, args ...any,
) (client.WorkflowRun, error) {
	m.env.ExecuteWorkflow(workflow, args...)
	return &mockWorkflowRun{env: m.env}, nil
}

func (m *MockClient) GetWorkflow(ctx context.Context, workflowID string, runID string) client.WorkflowRun {
	return &mockWorkflowRun{env: m.env}
}

func (m *MockClient) SignalWorkflow(ctx context.Context, workflowID string, runID string, signal string, arg any) error {
	m.env.SignalWorkflow(signal, arg)

	return nil
}

func (m *MockClient) SignalWithStartWorkflow(
	ctx context.Context, workflowID string, signalName string, signalArg any, options client.StartWorkflowOptions, workflow any, workflowArgs ...any,
) (client.WorkflowRun, error) {
	return nil, nil
}

func (m *MockClient) NewWithStartWorkflowOperation(options client.StartWorkflowOptions, workflow any, args ...any) client.WithStartWorkflowOperation {
	return nil
}

func (m *MockClient) CancelWorkflow(ctx context.Context, workflowID string, runID string) error {
	return nil
}

func (m *MockClient) TerminateWorkflow(ctx context.Context, workflowID string, runID string, reason string, details ...any) error {
	return nil
}

func (m *MockClient) GetWorkflowHistory(
	ctx context.Context, workflowID string, runID string, isLongPoll bool, filterType enums.HistoryEventFilterType,
) client.HistoryEventIterator {
	return nil
}

func (m *MockClient) CompleteActivity(ctx context.Context, taskToken []byte, result any, err error) error {
	return nil
}

func (m *MockClient) CompleteActivityByID(
	ctx context.Context, namespace, workflowID, runID, activityID string, result any, err error,
) error {
	return nil
}

func (m *MockClient) RecordActivityHeartbeat(ctx context.Context, taskToken []byte, details ...any) error {
	return nil
}

func (m *MockClient) RecordActivityHeartbeatByID(ctx context.Context, namespace, workflowID, runID, activityID string, details ...any) error {
	return nil
}

func (m *MockClient) ListClosedWorkflow(ctx context.Context, request *workflowservice.ListClosedWorkflowExecutionsRequest) (*workflowservice.ListClosedWorkflowExecutionsResponse, error) {
	return &workflowservice.ListClosedWorkflowExecutionsResponse{}, nil
}

func (m *MockClient) ListOpenWorkflow(ctx context.Context, request *workflowservice.ListOpenWorkflowExecutionsRequest) (*workflowservice.ListOpenWorkflowExecutionsResponse, error) {
	return &workflowservice.ListOpenWorkflowExecutionsResponse{}, nil
}

func (m *MockClient) ListWorkflow(ctx context.Context, request *workflowservice.ListWorkflowExecutionsRequest) (*workflowservice.ListWorkflowExecutionsResponse, error) {
	return &workflowservice.ListWorkflowExecutionsResponse{}, nil
}

func (m *MockClient) ListArchivedWorkflow(ctx context.Context, request *workflowservice.ListArchivedWorkflowExecutionsRequest) (*workflowservice.ListArchivedWorkflowExecutionsResponse, error) {
	return &workflowservice.ListArchivedWorkflowExecutionsResponse{}, nil
}

func (m *MockClient) ScanWorkflow(ctx context.Context, request *workflowservice.ScanWorkflowExecutionsRequest) (*workflowservice.ScanWorkflowExecutionsResponse, error) {
	return &workflowservice.ScanWorkflowExecutionsResponse{}, nil
}

func (m *MockClient) CountWorkflow(ctx context.Context, request *workflowservice.CountWorkflowExecutionsRequest) (*workflowservice.CountWorkflowExecutionsResponse, error) {
	return &workflowservice.CountWorkflowExecutionsResponse{}, nil
}

func (m *MockClient) GetSearchAttributes(ctx context.Context) (*workflowservice.GetSearchAttributesResponse, error) {
	return &workflowservice.GetSearchAttributesResponse{}, nil
}

func (m *MockClient) QueryWorkflow(ctx context.Context, workflowID string, runID string, queryType string, args ...any) (converter.EncodedValue, error) {
	return nil, nil
}

func (m *MockClient) QueryWorkflowWithOptions(ctx context.Context, request *client.QueryWorkflowWithOptionsRequest) (*client.QueryWorkflowWithOptionsResponse, error) {
	return &client.QueryWorkflowWithOptionsResponse{}, nil
}

func (m *MockClient) DescribeWorkflowExecution(ctx context.Context, workflowID, runID string) (*workflowservice.DescribeWorkflowExecutionResponse, error) {
	return &workflowservice.DescribeWorkflowExecutionResponse{}, nil
}

func (m *MockClient) DescribeWorkflow(ctx context.Context, workflowID, runID string) (*client.WorkflowExecutionDescription, error) {
	return &client.WorkflowExecutionDescription{}, nil
}

func (m *MockClient) DescribeTaskQueue(ctx context.Context, taskqueue string, taskqueueType enums.TaskQueueType) (*workflowservice.DescribeTaskQueueResponse, error) {
	return &workflowservice.DescribeTaskQueueResponse{}, nil
}

func (m *MockClient) ResetWorkflowExecution(ctx context.Context, request *workflowservice.ResetWorkflowExecutionRequest) (*workflowservice.ResetWorkflowExecutionResponse, error) {
	return &workflowservice.ResetWorkflowExecutionResponse{}, nil
}

func (m *MockClient) UpdateWorkerBuildIdCompatibility(ctx context.Context, options *client.UpdateWorkerBuildIdCompatibilityOptions) error {
	return nil
}

func (m *MockClient) GetWorkerBuildIdCompatibility(ctx context.Context, options *client.GetWorkerBuildIdCompatibilityOptions) (*client.WorkerBuildIDVersionSets, error) {
	return &client.WorkerBuildIDVersionSets{}, nil
}

func (m *MockClient) GetWorkerTaskReachability(ctx context.Context, options *client.GetWorkerTaskReachabilityOptions) (*client.WorkerTaskReachability, error) {
	return &client.WorkerTaskReachability{}, nil
}

func (m *MockClient) CheckHealth(ctx context.Context, request *client.CheckHealthRequest) (*client.CheckHealthResponse, error) {
	return &client.CheckHealthResponse{}, nil
}

func (m *MockClient) UpdateWorkflow(ctx context.Context, options client.UpdateWorkflowOptions) (client.WorkflowUpdateHandle, error) {
	return nil, nil
}

func (m *MockClient) UpdateWorkflowExecutionOptions(ctx context.Context, options client.UpdateWorkflowExecutionOptionsRequest) (client.WorkflowExecutionOptions, error) {
	return client.WorkflowExecutionOptions{}, nil
}

func (m *MockClient) UpdateWithStartWorkflow(ctx context.Context, options client.UpdateWithStartWorkflowOptions) (client.WorkflowUpdateHandle, error) {
	return nil, nil
}

func (m *MockClient) GetWorkflowUpdateHandle(ref client.GetWorkflowUpdateHandleOptions) client.WorkflowUpdateHandle {
	return nil
}

func (m *MockClient) WorkflowService() workflowservice.WorkflowServiceClient {
	return nil
}

func (m *MockClient) OperatorService() operatorservice.OperatorServiceClient {
	return nil
}

func (m *MockClient) ScheduleClient() client.ScheduleClient {
	return nil
}

func (m *MockClient) DeploymentClient() client.DeploymentClient {
	return nil
}

func (m *MockClient) WorkerDeploymentClient() client.WorkerDeploymentClient {
	return nil
}

func (m *MockClient) Close() {}

func (m *MockClient) DescribeTaskQueueEnhanced(ctx context.Context, options client.DescribeTaskQueueEnhancedOptions) (client.TaskQueueDescription, error) {
	return client.TaskQueueDescription{}, nil
}

func (m *MockClient) GetWorkerVersioningRules(ctx context.Context, options client.GetWorkerVersioningOptions) (*client.WorkerVersioningRules, error) {
	return nil, nil
}

func (m *MockClient) UpdateWorkerVersioningRules(ctx context.Context, options client.UpdateWorkerVersioningRulesOptions) (*client.WorkerVersioningRules, error) {
	return nil, nil
}

type mockWorkflowRun struct {
	env *testsuite.TestWorkflowEnvironment
}

func (r *mockWorkflowRun) GetID() string {
	return "mock-workflow-id"
}

func (r *mockWorkflowRun) GetRunID() string {
	return "mock-run-id"
}

func (r *mockWorkflowRun) Get(ctx context.Context, valuePtr any) error {
	return r.env.GetWorkflowResult(valuePtr)
}

func (r *mockWorkflowRun) GetWithOptions(ctx context.Context, valuePtr any, options client.WorkflowRunGetOptions) error {
	return r.env.GetWorkflowResult(valuePtr)
}
