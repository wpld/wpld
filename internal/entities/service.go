package entities

type Service struct {
	ID           string
	Network      string
	Aliases      []string
	Project      string
	Spec         Specification
	Domains      []string
	AttachStdout bool
	AttachStdin  bool
	AttachStderr bool
	Tty          bool
}
