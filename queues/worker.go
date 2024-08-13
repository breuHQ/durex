package queues

import (
	"context"
	"time"

	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/worker"
)

type (
	// WorkerOption is a function that configures a worker.Options struct.
	WorkerOption func(*worker.Options)
)

// WithWorkerOptionMaxConcurrentActivityExecutionSize sets the maximum concurrent activity executions this worker can have.
// The zero value of this uses the default value (1000).
func WithWorkerOptionMaxConcurrentActivityExecutionSize(size int) WorkerOption {
	return func(o *worker.Options) {
		o.MaxConcurrentActivityExecutionSize = size
	}
}

// WithWorkerOptionWorkerActivitiesPerSecond sets the rate limiting on number of activities that can be executed per second per worker.
// This can be used to limit resources used by the worker. The zero value of this uses the default value (100,000).
func WithWorkerOptionWorkerActivitiesPerSecond(rate float64) WorkerOption {
	return func(o *worker.Options) {
		o.WorkerActivitiesPerSecond = rate
	}
}

// WithWorkerOptionMaxConcurrentLocalActivityExecutionSize sets the maximum concurrent local activity executions this worker can have.
// The zero value of this uses the default value (1000).
func WithWorkerOptionMaxConcurrentLocalActivityExecutionSize(size int) WorkerOption {
	return func(o *worker.Options) {
		o.MaxConcurrentLocalActivityExecutionSize = size
	}
}

// WithWorkerOptionWorkerLocalActivitiesPerSecond sets the rate limiting on number of local
// activities that can be executed per second per worker.
//
// This can be used to limit resources used by the worker. The zero value of this uses the default value (100,000).
func WithWorkerOptionWorkerLocalActivitiesPerSecond(rate float64) WorkerOption {
	return func(o *worker.Options) {
		o.WorkerLocalActivitiesPerSecond = rate
	}
}

// WithWorkerOptionTaskQueueActivitiesPerSecond sets the rate limiting on number of
// activities that can be executed per second for the entire task queue.
//
// This is managed by the server. The zero value of this uses the default value (100,000).
func WithWorkerOptionTaskQueueActivitiesPerSecond(rate float64) WorkerOption {
	return func(o *worker.Options) {
		o.TaskQueueActivitiesPerSecond = rate
	}
}

// WithWorkerOptionMaxConcurrentActivityTaskPollers sets the maximum number of goroutines that will concurrently poll the
// temporal-server to retrieve activity tasks. The default value is 2.
func WithWorkerOptionMaxConcurrentActivityTaskPollers(count int) WorkerOption {
	return func(o *worker.Options) {
		o.MaxConcurrentActivityTaskPollers = count
	}
}

// WithWorkerOptionMaxConcurrentWorkflowTaskExecutionSize sets the maximum concurrent workflow task executions this worker can have.
// The zero value of this uses the default value (1000). This value cannot be 1.
func WithWorkerOptionMaxConcurrentWorkflowTaskExecutionSize(size int) WorkerOption {
	return func(o *worker.Options) {
		o.MaxConcurrentWorkflowTaskExecutionSize = size
	}
}

// WithWorkerOptionMaxConcurrentWorkflowTaskPollers sets the maximum number of goroutines that will concurrently poll the
// temporal-server to retrieve workflow tasks. The default value is 2. This value cannot be 1.
func WithWorkerOptionMaxConcurrentWorkflowTaskPollers(count int) WorkerOption {
	return func(o *worker.Options) {
		o.MaxConcurrentWorkflowTaskPollers = count
	}
}

// WithWorkerOptionEnableLoggingInReplay enables logging in replay mode. This is only useful for debugging purposes.
// The default value is false.
func WithWorkerOptionEnableLoggingInReplay(enable bool) WorkerOption {
	return func(o *worker.Options) {
		o.EnableLoggingInReplay = enable
	}
}

// WithWorkerOptionStickyScheduleToStartTimeout sets the sticky schedule to start timeout.
// The default value is 5 seconds.
func WithWorkerOptionStickyScheduleToStartTimeout(timeout time.Duration) WorkerOption {
	return func(o *worker.Options) {
		o.StickyScheduleToStartTimeout = timeout
	}
}

// WithWorkerOptionBackgroundActivityContext sets the root context for all activities.
func WithWorkerOptionBackgroundActivityContext(ctx context.Context) WorkerOption {
	return func(o *worker.Options) {
		o.BackgroundActivityContext = ctx
	}
}

// WithWorkerOptionWorkflowPanicPolicy sets how the workflow worker deals with non-deterministic history events and panics.
// The default value is BlockWorkflow.
func WithWorkerOptionWorkflowPanicPolicy(policy worker.WorkflowPanicPolicy) WorkerOption {
	return func(o *worker.Options) {
		o.WorkflowPanicPolicy = policy
	}
}

// WithWorkerOptionWorkerStopTimeout sets the worker graceful stop timeout.
// The default value is 0 seconds.
func WithWorkerOptionWorkerStopTimeout(timeout time.Duration) WorkerOption {
	return func(o *worker.Options) {
		o.WorkerStopTimeout = timeout
	}
}

// WithWorkerOptionEnableSessionWorker enables running session workers for activities within a session.
// The default value is false.
func WithWorkerOptionEnableSessionWorker(enable bool) WorkerOption {
	return func(o *worker.Options) {
		o.EnableSessionWorker = enable
	}
}

// WithWorkerOptionMaxConcurrentSessionExecutionSize sets the maximum number of concurrently running sessions the resource supports.
// The default value is 1000.
func WithWorkerOptionMaxConcurrentSessionExecutionSize(size int) WorkerOption {
	return func(o *worker.Options) {
		o.MaxConcurrentSessionExecutionSize = size
	}
}

// WithWorkerOptionDisableWorkflowWorker disables the workflow worker for this worker.
// The default value is false.
func WithWorkerOptionDisableWorkflowWorker(disable bool) WorkerOption {
	return func(o *worker.Options) {
		o.DisableWorkflowWorker = disable
	}
}

// WithWorkerOptionLocalActivityWorkerOnly sets the worker to only handle workflow tasks and local activities.
// The default value is false.
func WithWorkerOptionLocalActivityWorkerOnly(localOnly bool) WorkerOption {
	return func(o *worker.Options) {
		o.LocalActivityWorkerOnly = localOnly
	}
}

// WithWorkerOptionIdentity sets the identity for the worker, overwriting the client-level Identity value.
func WithWorkerOptionIdentity(identity string) WorkerOption {
	return func(o *worker.Options) {
		o.Identity = identity
	}
}

// WithWorkerOptionDeadlockDetectionTimeout sets the maximum amount of time that a workflow task will be allowed to run.
// The default value is 1 second.
func WithWorkerOptionDeadlockDetectionTimeout(timeout time.Duration) WorkerOption {
	return func(o *worker.Options) {
		o.DeadlockDetectionTimeout = timeout
	}
}

// WithWorkerOptionMaxHeartbeatThrottleInterval sets the maximum amount of time between sending each pending heartbeat to the server.
// The default value is 60 seconds.
func WithWorkerOptionMaxHeartbeatThrottleInterval(interval time.Duration) WorkerOption {
	return func(o *worker.Options) {
		o.MaxHeartbeatThrottleInterval = interval
	}
}

// WithWorkerOptionDefaultHeartbeatThrottleInterval sets the default amount of time between sending each pending heartbeat to the server.
// The default value is 30 seconds.
func WithWorkerOptionDefaultHeartbeatThrottleInterval(interval time.Duration) WorkerOption {
	return func(o *worker.Options) {
		o.DefaultHeartbeatThrottleInterval = interval
	}
}

// WithWorkerOptionInterceptors sets the interceptors to apply to the worker.
func WithWorkerOptionInterceptors(interceptors []interceptor.WorkerInterceptor) WorkerOption {
	return func(o *worker.Options) {
		o.Interceptors = interceptors
	}
}

// WithWorkerOptionOnFatalError sets the callback invoked on fatal error.
func WithWorkerOptionOnFatalError(fn func(error)) WorkerOption {
	return func(o *worker.Options) {
		o.OnFatalError = fn
	}
}

// WithWorkerOptionDisableEagerActivities disables eager activities.
// The default value is false.
func WithWorkerOptionDisableEagerActivities(disable bool) WorkerOption {
	return func(o *worker.Options) {
		o.DisableEagerActivities = disable
	}
}

// WithWorkerOptionMaxConcurrentEagerActivityExecutionSize sets the maximum number of eager activities that can be running.
// The default value of 0 means unlimited.
func WithWorkerOptionMaxConcurrentEagerActivityExecutionSize(size int) WorkerOption {
	return func(o *worker.Options) {
		o.MaxConcurrentEagerActivityExecutionSize = size
	}
}

// WithWorkerOptionDisableRegistrationAliasing disables allowing workflow and activity functions registered with custom names
// from being called with their function references.
// The default value is false.
func WithWorkerOptionDisableRegistrationAliasing(disable bool) WorkerOption {
	return func(o *worker.Options) {
		o.DisableRegistrationAliasing = disable
	}
}

// WithWorkerOptionBuildID assigns a BuildID to this worker.
// NOTE: Experimental.
func WithWorkerOptionBuildID(buildID string) WorkerOption {
	return func(o *worker.Options) {
		o.BuildID = buildID
	}
}

// WithWorkerOptionUseBuildIDForVersioning opts this worker into the Worker Versioning feature.
// NOTE: Experimental.
func WithWorkerOptionUseBuildIDForVersioning(use bool) WorkerOption {
	return func(o *worker.Options) {
		o.UseBuildIDForVersioning = use
	}
}

// NewWorkerOptions creates a new worker.Options struct with the given options applied.
//
//	Example usage:
//	workerOptions := NewWorkerOptions(
//		WithWorkerOptionOnFatalError(func(err error) {
//			// Handle fatal error
//		}),
//		WithWorkerOptionDisableEagerActivities(true),
//		WithWorkerOptionMaxConcurrentEagerActivityExecutionSize(10),
//		WithWorkerOptionDisableRegistrationAliasing(false),
//		WithWorkerOptionBuildID("my-build-id"),
//		WithWorkerOptionUseBuildIDForVersioning(true),
//	)
func NewWorkerOptions(opts ...WorkerOption) worker.Options {
	options := worker.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return options
}
