package tasks

import (
	"errors"
)

var (
	ProjectNotFoundErr = errors.New("project not found")
	ServiceNotFoundErr = errors.New("service not found")
)
