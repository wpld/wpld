package models

import (
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

type Network struct {
	External map[string]string `yaml:external,omitempty`
}

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
	Networks map[string]Network `yaml:"networks"`
	Volumes map[string]interface{} `yaml:"volumes,omitempty"`
}

func (c Compose) Save(fs afero.Fs) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = afero.WriteFile(fs, "docker-compose.yaml", data, 0644)
	if err != nil {
		return err
	}

	return nil
}
