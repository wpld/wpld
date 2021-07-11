package entities

type Project struct {
	Name     string             `yaml:"name"`
	Domains  []string           `yaml:"domains"`
	Volumes  []string           `yaml:"volumes"`
	Services map[string]Service `yaml:"services"`
}
