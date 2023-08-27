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
	"errors"
	"fmt"
)

var (
	// ErrParentNil is returned when the parent workflow context is nil.
	ErrParentNil = errors.New("parent workflow context is nil")
)

type (
	DuplicateIDPropError struct {
		prop  string
		value string
		next  string
	}
)

func (e *DuplicateIDPropError) Error() string {
	return fmt.Sprintf("attempting to override %s existing value %s with %s", e.prop, e.value, e.next)
}

func NewDuplicateIDPropError(prop, value, next string) error {
	return &DuplicateIDPropError{prop, value, next}
}
