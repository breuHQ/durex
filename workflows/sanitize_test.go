package workflows_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"go.breu.io/durex/workflows"
)

type (
	SanitizeTestSuite struct {
		suite.Suite
	}
)

func (suite *SanitizeTestSuite) TestSanitize() {
	cases := []struct {
		input  string
		output string
	}{
		{"Hello World", "hello_world"},
		{"Hello-World", "hello-world"},
		{"hello-world", "hello-world"},
		{"  hello-world   ", "hello-world"},
		{"hello world!", "hello_world"},
		{"hello.world", "hello_world"},
		{"-----hello---world----", "hello-world"},
		{"hello--world---", "hello-world"},
		{"  hello--world   ", "hello-world"},
		{"!!!", "_"},
		{"----", "-"},
		{"", ""},
		{"   ", "_"},
		{" ---hello--- ", "hello"},
		{"  hello--world--  ", "hello-world"},
		{"-hello-", "hello"},
		{"hello---", "hello"},
		{"____hello____", "hello"},
		{"_hello_", "hello"},
		{"hello---world---", "hello-world"},
		{"-hello-world-", "hello-world"},
		{"hello-world-", "hello-world"},
		{"-hello-world", "hello-world"},
		{"hello-world_", "hello-world"},
		{"_-hello-world_", "hello-world"},
		{"_hello", "hello"},
		{"hello_", "hello"},
	}

	for _, testCase := range cases {
		suite.Run(testCase.input, func() {
			actual := workflows.Sanitize(testCase.input)
			suite.Equal(testCase.output, actual, "Sanitize(%q) should be %q, but got %q", testCase.input, testCase.output, actual)
		})
	}
}

func TestSanitizeSuite(t *testing.T) {
	suite.Run(t, new(SanitizeTestSuite))
}
