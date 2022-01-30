package tasks

import (
	"errors"
)

var (
	ProjectNotFoundErr   = errors.New("project not found")
	ServiceNotFoundErr   = errors.New("service not found")
	WpServiceNotFoundErr = errors.New("wp service not found")
	ScriptNotFoundErr    = errors.New("script not found")
)
