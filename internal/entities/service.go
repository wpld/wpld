package entities

type Service struct {
	ID      string
	Network string
	Aliases []string
	Project string
	Spec    Specification
	Domains []string
}
