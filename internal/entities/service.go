package entities

type Service struct {
	ID           string
	Network      string
	Aliases      []string
	Project      string
	Spec         Specification
	AttachStdout bool
	AttachStdin  bool
	AttachStderr bool
	Tty          bool
}
