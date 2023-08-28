# go-temporal-tools: utilities for temporal.io go sdk

![GitHub release (with filter)](https://img.shields.io/github/v/release/breuHQ/go-temporal-tools)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/breuHQ/go-temporal-tools)
[![License](https://img.shields.io/github/license/breuHQ/go-temporal-tools)](./LICENSE)
![GitHub contributors](https://img.shields.io/github/contributors/breuHQ/go-temporal-tools)

## 🚀 Install

```sh
go get go.breu.io/temporal-tools
```

**Compatibility**: go >= 1.21

- ⚠️ Work in Progress

## Why?

After working with temporal.io across multiple projects, we have standarized a set of best practices across our projects.
For example, to create a unique & identifiable workflow id from the UI, we have found that following the [block, element, modifier](https://getbem.com/introduction/) method, a technique for writing maintainable CSS, makes it very readable and maintainable.

We also found that tying the workflow to a queue makes it very easy. For us, `Queue` is where it all begins.

## Getting Started

We start by creating a `Queue` first and then call `ExecuteWorkflow` or `ExecuteChildWorkflow` methods on the `Queue` interface.

### temporal.io setup

Given the `temporal.go` file below, we are going to setup our queues as global varibales like

```go
package temporal

import (
  "fmt"
  "log/slog"
  "sync"
  "time"

  "github.com/avast/retry-go/v4"
  "github.com/ilyakaznacheev/cleanenv"
  "go.breu.io/slog-utils/calldepth"
  "go.temporal.io/sdk/client"
)

type (
  Temporal interface {
    Client() client.Client
    ConnectionString() string
  }

  config struct {
    Host   string `env:"TEMPORAL_HOST" env-default:"localhost"`
    Port   int    `env:"TEMPORAL_PORT" env-default:"7233"`
    client client.Client
    logger *slog.Logger
  }

  Option func(*config)
)

var (
  once sync.Once
)

const (
  MaxAttempts = 10
)

func (c *config) ConnectionString() string {
  return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *config) Client() client.Client {
  if c.client == nil {
  once.Do(func() {
    slog.Info(
      "initializing temporal connection",
      slog.String("connection_string", c.ConnectionString()),
    )

    options := client.Options{
      HostPort: c.ConnectionString(),
      Logger: c.logger,
    }

    retryTemporal := func() error {
      clt, err := client.Dial(options)
      if err != nil {
        return err
      }

      c.client = clt

      slog.Info("temporal: connected")

      return nil
    }

    if err := retry.Do(
      retryTemporal,
      retry.Attempts(MaxAttempts),
      retry.Delay(1*time.Second),
      retry.OnRetry(func(n uint, err error) {
        slog.Info(
          "temporal: failed to connect. retrying connection ...",
          slog.Any("attempt", n+1),
          slog.Any("max attempts", MaxAttempts),
          slog.String("error", err.Error()),
        )
      }),
      ); err != nil {
        panic(fmt.Errorf("failed to connect to temporal: %w", err))
      }
    })
  }

  return c.client
}

func WithHost(host string) Option {
  return func(c *config) {
  c.Host = host
  }
}

func WithPort(port int) Option {
  return func(c *config) {
  c.Port = port
  }
}

func WithLogger(logger *slog.Logger) Option {
  return func(c *config) {
  c.logger = logger
  }
}

func FromEnvironment() Option {
  return func(c *config) {
  if err := cleanenv.ReadEnv(c); err != nil {
  panic(fmt.Errorf("failed to read environment: %w", err))
  }
  }
}

func New(opts ...Option) Temporal {
  c := &config{}

  for _, opt := range opts {
  opt(c)
  }

  return c
}
```

### `Temporal` singleton

```go
package shared

import (
  "os"
  "sync"

  "shared/temporal"

  "go.breu.io/slog-utils/calldepth"
)

var (
  tmprl     temporal.Temporal // Global temporal instance.
  tmprlOnce sync.Once         // Global temporal initialization state.
)

// Temporal returns the global temporal instance.
func Temporal() temporal.Temporal {
  tmprlOnce.Do(func() {
    tmprl = temporal.New(
      temporal.FromEnvironment(),
      temporal.WithLogger(
        calldepth.NewAdapter(
          calldepth.NewLogger(slog.NewJsonLogger()),
          calldepth.WithCallDepth(5), // 5 for activities, 6 for workflows.
         ).WithGroup("temporal")
      ),
    )
  })

  return tmprl
}
```

### `Queue` setup

```go
package shared

import (
  "sync"

  "go.breu.io/temporal-tools/queues"
)

var (
  coreq queues.Queue
  coreqone sync.once
)

// The default value right now is "io.breu".
// Change this as per your requirements. A good practice is to use java's package notation.
func init() {
  queues.SetDefaultPrefix("com.company.pkg")
}

func CoreQueue(
  once.Do(func() {
    queues.New(
      WithName("core"),
      // We have already setup temporal as a signleton in the shared package. We leverage the instantiatilized client.
      WithClient(
        Temporal().Client(),
      ),
    )
  })
)
```

## Usage for Workflows

We make the assumptions that a workflow is tied to a queue. So by calling the Queue.ExecuteWorkflow() to start a workflow.

```go
package main

import (
  "context"
  "shared"

  "github.com/google/uuid"
  "go.temporal.io/sdk/workflow"

  "go.breu.io/temporal-tools/workflows"
)

func main() {
  worker := shared.CoreQueue().CreateWorker()
  defer worker.Stop()

  worker.RegisterWorkflow(Workflow)
  worker.RegisterWorkflow(ChildWorkflow)

  if err != worker.Start(); err != nil {
    // handler error
  }


  opts, err := workflows.NewOptions(
    workflows.WithBlock("block"),
    workflows.WithBlockID(uuid.New()), // d5e012df-5b9e-41cf-9ed5-3439eeafd8e4
  )
  if err != nil {
    // handle error
  }

  exe, err := shared.CoreQueue().ExecuteWorkflow(
    context.Background(),
    // The workflow id created in this case would be
    //  "com.company.pkg.block.d5e012df-5b9e-41cf-9ed5-3439eeafd8e4"
    opts,
    Workflow, // or workflow function name
    payload,
  )
  if err != nil {
    // handle error
  }
}

func Workflow(ctx context.Context, payload any) error {
  childopts, err := workflow.NewOptions(
    WithParent(ctx),
    WithBlock("child"),
    WithBlockID(uuid.New()), // bb42d90a-d3f3-4022-bc09-aeed6e9db659
  )
  if err != nil {
    return err
  }

  future := shared.CoreQueue().ExecuteChildWorkflow(
    ctx,
    // By passing WithParent(ctx), the helper method automagically picks up the right workflow to create the child id i.e.
    // "com.company.pkg.block.d5e012df-5b9e-41cf-9ed5-3439eeafd8e4.child.bb42d90a-d3f3-4022-bc09-aeed6e9db659"
    opts,
    ChildWorkflow,
    payload,
  )

  // Do Something

  return nil
}

func ChildWorkflow(ctx workflow.Context, payload any) {
  // Child workflow
}
```


## 👤 Contributors

![Contributors](https://contrib.rocks/image?repo=breuHQ/go-temporal-tools)


## 📝 License

Copyright © 2023 [Breu Inc.](https://github.com/breuHQ)

This project is [MIT](./LICENSE) licensed.
