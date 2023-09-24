package workflows

import (
	"regexp"

	"github.com/gobeam/stringy"
)

var (
	uuidrxstr = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$" // taken from github.com/go-playground/validator
)

var (
	uuidrx = regexp.MustCompile(uuidrxstr)
)

// format sanitizes a string and returns it.
func format(s string) string {
	if uuidrx.MatchString(s) {
		return s
	}

	return stringy.New(s).SnakeCase("?", "", "#", "").ToLower()
}
