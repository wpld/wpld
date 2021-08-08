package entities

import (
	"fmt"
)

type Specification struct {
	Name         string            `yaml:"name"`
	Image        string            `yaml:"image,omitempty"`
	WorkingDir   string            `yaml:"working_dir,omitempty"`
	Entrypoint   []string          `yaml:"entrypoint,omitempty"`
	ExposedPorts []string          `yaml:"exposed_ports,omitempty"`
	CapAdd       []string          `yaml:"cap_add,omitempty"`
	CapDrop      []string          `yaml:"cap_drop,omitempty"`
	Cmd          []string          `yaml:"cmd,omitempty"`
	Volumes      []string          `yaml:"volumes,omitempty"`
	VolumesFrom  []string          `yaml:"volumes_from,omitempty"`
	Ports        []string          `yaml:"ports,omitempty"`
	Env          map[string]string `yaml:"env,omitempty"`
	DependsOn    []string          `yaml:"depends_on,omitempty"`
	Domains      []string          `yaml:"domains,omitempty"`
	IPAddress    string            `yaml:"ipv4_address,omitempty"`
}

func (s Specification) GetEnvs() []string {
	i := 0
	envs := make([]string, len(s.Env))

	for key, value := range s.Env {
		envs[i] = fmt.Sprintf("%s=%s", key, value)
		i++
	}

	return envs
}
