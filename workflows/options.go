package workflows

import (
	"strings"

	"go.temporal.io/sdk/workflow"
)

const (
	RetryForever int32 = 0 // Default workflow max Attempts. 0 means forever.
)

type (
	// Options defines the interface for creating workflow options.
	Options interface {
		IsChild() bool                     // IsChild returns true if the workflow id is a child workflow id.
		ParentWorkflowID() (string, error) // ParentWorkflowID returns the parent workflow id.
		IDSuffix() string                  // IDSuffix santizes the suffix of the workflow id and then formats it as a string.
		MaxAttempts() int32                // MaxAttempts returns the max attempts for the workflow.
		IgnoredErrors() []string           // IgnoredErrors returns the list of errors that are ok to ignore.
	}

	// Option sets the specified options.
	Option func(Options) error

	// props defines the interface for creating id properties.
	props map[string]string

	options struct {
		parent_context workflow.Context // The parent workflow context.

		ParentID       string   `json:"parent"`          // The parent workflow ID.
		Block          string   `json:"block"`           // The block name.
		BlockID        string   `json:"block_id"`        // The block identifier.
		Element        string   `json:"element"`         // The element name.
		ElementID      string   `json:"element_id"`      // The element identifier.
		Mod            string   `json:"mod"`             // The modifier name.
		ModID          string   `json:"mod_id"`          // The modifier identifier.
		Props          props    `json:"props"`           // The properties associated with the workflow.
		Order          []string `json:"order"`           // The order of execution for the workflow elements.
		MaximumAttempt int32    `json:"maximum_attempt"` // The maximum number of attempts for the workflow.
		IgnoreErrors   []string `json:"ignore_errors"`   // A list of errors to ignore during the workflow execution.
	}
)

// IsChild returns true if the workflow id is a child workflow id.
func (w *options) IsChild() bool {
	return w.ParentID != ""
}

func (w *options) ID() string {
	id := w.IDSuffix()

	if w.IsChild() {
		return w.ParentID + "." + id
	}

	return id
}

// IDSuffix sanitizes the suffix and returns it.
func (w *options) IDSuffix() string {
	parts := []string{w.Block, w.BlockID, w.Element, w.ElementID, w.Mod, w.ModID}
	for _, key := range w.Order {
		parts = append(parts, key, w.Props[key])
	}

	sanitized := make([]string, 0)

	// removing empty strings and trimming spaces
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			sanitized = append(sanitized, part)
		}
	}

	return strings.Join(sanitized, ".")
}

// ParentWorkflowID returns the parent workflow id.
func (w *options) ParentWorkflowID() (string, error) {
	if w.parent_context == nil {
		return "", ErrParentNil
	}

	return workflow.GetInfo(w.parent_context).WorkflowExecution.ID, nil
}

// MaxAttempts returns the max attempts for the workflow.
func (w *options) MaxAttempts() int32 {
	return w.MaximumAttempt
}

// IgnoredErrors returns the list of errors that are ok to ignore.
func (w *options) IgnoredErrors() []string {
	return w.IgnoreErrors
}

// WithParent sets the parent workflow context.
func WithParent(parent workflow.Context) Option {
	return func(o Options) error {
		o.(*options).ParentID = workflow.GetInfo(parent).WorkflowExecution.ID
		return nil
	}
}

// WithBlock sets the block name.
func WithBlock(val string) Option {
	return func(o Options) error {
		if o.(*options).Block != "" {
			return NewDuplicateIDPropError("block", o.(*options).Block, val)
		}

		o.(*options).Block = sanitize(val)

		return nil
	}
}

// WithBlockID sets the block value.
func WithBlockID(val string) Option {
	return func(o Options) error {
		if o.(*options).BlockID != "" {
			return NewDuplicateIDPropError("blockID", o.(*options).BlockID, val)
		}

		o.(*options).BlockID = sanitize(val)

		return nil
	}
}

// WithElement sets the element name.
func WithElement(val string) Option {
	return func(o Options) error {
		if o.(*options).Element != "" {
			return NewDuplicateIDPropError("element", o.(*options).Element, val)
		}

		o.(*options).Element = sanitize(val)

		return nil
	}
}

// WithElementID sets the element value.
func WithElementID(val string) Option {
	return func(o Options) error {
		if o.(*options).ElementID != "" {
			return NewDuplicateIDPropError("element id", o.(*options).ElementID, val)
		}

		o.(*options).ElementID = sanitize(val)

		return nil
	}
}

// WithMod sets the modifier name.
func WithMod(val string) Option {
	return func(o Options) error {
		if o.(*options).Mod != "" {
			return NewDuplicateIDPropError("modifier", o.(*options).Mod, val)
		}

		o.(*options).Mod = sanitize(val)

		return nil
	}
}

// WithModID sets the modifier value.
func WithModID(val string) Option {
	return func(o Options) error {
		if o.(*options).ModID != "" {
			return NewDuplicateIDPropError("modifier id", o.(*options).ModID, val)
		}

		o.(*options).ModID = sanitize(val)

		return nil
	}
}

// WithProp sets the prop given a key & value.
func WithProp(key, val string) Option {
	return func(o Options) error {
		o.(*options).Order = append(o.(*options).Order, key)
		o.(*options).Props[sanitize(key)] = sanitize(val)

		return nil
	}
}

// WithMaxAttempts sets the max attempts for the workflow.
func WithMaxAttempts(attempts int32) Option {
	return func(o Options) error {
		o.(*options).MaximumAttempt = attempts
		return nil
	}
}

// WithIgnoredError adds an error to the list of errors that are ok to ignore for the workflow.
func WithIgnoredError(err string) Option {
	return func(o Options) error {
		o.(*options).IgnoreErrors = append(o.(*options).IgnoreErrors, err)
		return nil
	}
}

// NewOptions sets workflow options required to run a workflow like workflow id, max attempts, etc.
//
// The idempotent workflow ID Sometimes we need to signal the workflow from a completely disconnected part of the
// application. For us, it is important to arrive at the same workflow ID regardless of the conditions. We try to follow
// the block, element, modifier pattern popularized by advocates of maintainable CSS. For more info,
// https://getbem.com/.
//
// Example: For the block github with installation id 123, the element being the repository with id 456, and the
// modifier being the pull request with id 789, we would call
//
//	opts := NewOptions(
//	  WithBlock("github"),
//	  WithBlockID("123"),
//	  WithElement("repository"),
//	  WithElementID("456"),
//	  WithModifier("repository"),
//	  WithModifierID("789"),
//	  // Sometimes we need to over-ride max attempts. The default is workflows.RetryForever.
//	  WithMaxAttempts(3),
//	  // If we want to ignore an error, we can do so.
//	  WithIgnoredError("SomeError"),
//	)
//
//	id := opts.ID()
//
// NOTE: The design is work in progress and may change in future.
func NewOptions(opts ...Option) (Options, error) {
	var (
		err error
	)

	w := &options{
		Props:          make(props),
		Order:          make([]string, 0),
		MaximumAttempt: RetryForever,
		IgnoreErrors:   make([]string, 0),
	}

	for _, opt := range opts {
		err = opt(w)
		if err != nil {
			continue
		}
	}

	return w, err
}
