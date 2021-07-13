package entities

type Specification struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image,omitempty"`
	WorkingDir  string            `yaml:"working_dir,omitempty"`
	Entrypoint  []string          `yaml:"entrypoint,omitempty"`
	CapAdd      []string          `yaml:"cap_add,omitempty"`
	CapDrop     []string          `yaml:"cap_drop,omitempty"`
	Cmd         []string          `yaml:"cmd,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	VolumesFrom []string          `yaml:"volumes_from,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Env         map[string]string `yaml:"env,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty"`
}
