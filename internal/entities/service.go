package entities

type Service struct {
	ID      string
	Network string
	Alias   string
	Project string
	Spec    Specification
	Domains []string
}
