package compose

type Service struct {
	Name string `yaml:"name,omitempty"`
	Image string `yaml:"image,omitempty"`
	Build Build `yaml:"build,omitempty"`
	Volumes []string `yaml:"volumes,omitempty"`
	Env map[string]string `yaml:"env,omitempty"`
}
