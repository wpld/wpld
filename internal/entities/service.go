package entities

type Service struct {
	ID           string
	Network      Network
	Aliases      []string
	Project      string
	Spec         Specification
	AttachStdout bool
	AttachStdin  bool
	AttachStderr bool
	Tty          bool
}
