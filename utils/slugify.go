package utils

import (
	"regexp"
	"strings"
)

var slugifyRegexp = regexp.MustCompile(`\W+`)

func Slugify(val string) string {
	result := string(slugifyRegexp.ReplaceAll([]byte(val), []byte("-")))
	result = strings.ToLower(result)
	result = strings.Trim(result, "-")

	return result
}
