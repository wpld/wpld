package utils

type Service struct {
	Image string `yaml:"image,omitempty"`
	Expose []int `yaml:"expose,omitempty"`
	Volumes []string `yaml:"volumes,omitempty"`
	DependsOn []string `yaml:"depends_on,omitempty"`
	Networks []string `yaml:"networks,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
}

type Compose struct {
	Version string `yaml:"version"`
	Services map[string]Service `yaml:"services"`
	Networks map[string]interface{} `yaml:"networks"`
	Volumes map[string]interface{} `yaml:"volumes,omitempty"`
}
