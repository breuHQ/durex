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

package workflows

import (
	"strings"

	"github.com/gobeam/stringy"
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
	}

	// Option sets the specified options.
	Option func(Options) error

	// props defines the interface for creating id properties.
	props map[string]string

	// optionProvider defines the interface for creating workflow options.
	options struct {
		parent workflow.Context

		block     string
		blockID   string
		elm       string
		elmID     string
		mod       string
		modID     string
		props     props    // props is a map of id properties.
		propOrder []string // propOrder is the order in which the properties are added.

		maxattempts int32
	}
)

func (w *options) IsChild() bool {
	return w.parent != nil
}

// IDSuffix sanitizes the suffix and returns it.
func (w *options) IDSuffix() string {
	parts := []string{w.block, w.blockID, w.elm, w.elmID, w.mod, w.modID}
	for _, key := range w.propOrder {
		parts = append(parts, key, w.props[key])
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

func (w *options) ParentWorkflowID() (string, error) {
	if w.parent == nil {
		return "", ErrParentNil
	}

	return workflow.GetInfo(w.parent).WorkflowExecution.ID, nil
}

// MaxAttempts returns the max attempts for the workflow.
func (w *options) MaxAttempts() int32 {
	return w.maxattempts
}

// WithParent sets the parent workflow context.
func WithParent(parent workflow.Context) Option {
	return func(o Options) error {
		o.(*options).parent = parent
		return nil
	}
}

// WithBlock sets the block name.
func WithBlock(block string) Option {
	return func(o Options) error {
		if o.(*options).block != "" {
			return NewDuplicateIDPropError("block", o.(*options).block, block)
		}

		blockstr := stringy.New(block)
		block := blockstr.SnakeCase("?", "", "#", "").ToLower()

		o.(*options).block = block

		return nil
	}
}

// WithBlockID sets the block value.
func WithBlockID(val string) Option {
	return func(o Options) error {
		if o.(*options).blockID != "" {
			return NewDuplicateIDPropError("blockID", o.(*options).blockID, val)
		}

		valstr := stringy.New(val)
		val = valstr.SnakeCase("?", "", "#", "").ToLower()

		o.(*options).blockID = val

		return nil
	}
}

// WithElement sets the element name.
func WithElement(element string) Option {
	return func(o Options) error {
		if o.(*options).elm != "" {
			return NewDuplicateIDPropError("element", o.(*options).elm, element)
		}

		elementstr := stringy.New(element)
		element := elementstr.SnakeCase("?", "", "#", "").ToLower()

		o.(*options).elm = element

		return nil
	}
}

// WithElementID sets the element value.
func WithElementID(val string) Option {
	return func(o Options) error {
		if o.(*options).elmID != "" {
			return NewDuplicateIDPropError("element id", o.(*options).elmID, val)
		}

		valstr := stringy.New(val)
		val = valstr.SnakeCase("?", "", "#", "").ToLower()

		o.(*options).elmID = val

		return nil
	}
}

// WithMod sets the modifier name.
func WithMod(modifier string) Option {
	return func(o Options) error {
		if o.(*options).mod != "" {
			return NewDuplicateIDPropError("modifier", o.(*options).mod, modifier)
		}

		modifierstr := stringy.New(modifier)
		modifier := modifierstr.SnakeCase("?", "", "#", "").ToLower()

		o.(*options).mod = modifier

		return nil
	}
}

// WithModID sets the modifier value.
func WithModID(val string) Option {
	return func(o Options) error {
		if o.(*options).modID != "" {
			return NewDuplicateIDPropError("modifier id", o.(*options).modID, val)
		}

		valstr := stringy.New(val)
		val = valstr.SnakeCase("?", "", "#", "").ToLower()

		o.(*options).modID = val

		return nil
	}
}

// WithProp sets the prop given a key & value.
func WithProp(key, val string) Option {
	return func(o Options) error {
		o.(*options).propOrder = append(o.(*options).propOrder, key)
		o.(*options).props[stringy.New(key).SnakeCase("?", "", "#", "").ToLower()] = stringy.New(val).SnakeCase("?", "", "#", "").ToLower()

		return nil
	}
}

// WithMaxAttempts sets the max attempts for the workflow.
func WithMaxAttempts(attempts int32) Option {
	return func(o Options) error {
		o.(*options).maxattempts = attempts
		return nil
	}
}

// NewOptions sets workflow options required to run a workflow like workflow id, max attempts, etc.
//
// The idempotent workflow ID Sometimes we need to signal the workflow from a completely disconnected
// part of the application. For us, it is important to arrive at the same workflow ID regardless of the conditions.
// We try to follow the block, element, modifier pattern popularized by advocates of mantainable CSS. For more info,
// https://getbem.com.
//
// Example:
// For the block github with installation id 123, the element being the repository with id 456, and the modifier being the
// pull request with id 789, we would call
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
		props:       make(props),
		propOrder:   make([]string, 0),
		maxattempts: RetryForever,
	}

	for _, opt := range opts {
		err = opt(w)
		if err != nil {
			continue
		}
	}

	return w, err
}
