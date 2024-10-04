package workflows

import (
	"regexp"

	"github.com/gobeam/stringy"
)

var (
	uuid = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$" // taken from github.com/go-playground/validator
)

var (
	re_uuid = regexp.MustCompile(uuid)
)

// sanitize sanitizes a string and returns it.
func sanitize(s string) string {
	if re_uuid.MatchString(s) {
		return s
	}

	return stringy.New(s).SnakeCase("?", "", "#", "").ToLower()
}
