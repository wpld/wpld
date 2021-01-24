package utils

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func ParseBytes(val string) (int64, error) {
	re := regexp.MustCompile(`(?i)^(\d+)(K|M|G)(i?)B$`)
	matches := re.FindSubmatch([]byte(val))
	if matches == nil {
		return 0, errors.New("wrong format of the bytes value")
	}

	var base float64 = 1000
	if string(matches[3]) == "i" {
		base = 1024
	}

	switch strings.ToLower(string(matches[2])) {
	case "m":
		base = math.Pow(base, 2)
	case "g":
		base = math.Pow(base, 3)
	}

	result, err := strconv.Atoi(string(matches[1]))
	if err != nil {
		return 0, err
	}

	return int64(float64(result) * base), nil
}
