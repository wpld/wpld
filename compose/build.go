package compose

type Build struct {
	Name       string            `yaml:"name"`
	Dockerfile string            `yaml:"dockerfile,omitempty"`
	Context    string            `yaml:"context,omitempty"`
	Args       map[string]string `yaml:"args,omitempty"`
}
