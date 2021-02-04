package compose

import "gopkg.in/yaml.v2"

type Compose struct {
	Name     string             `yaml:"name"`
	Hostname string             `yaml:"hostname"`
	Services map[string]Service `yaml:"services"`
}

func (c Compose) Serialize() []byte {
	if data, err := yaml.Marshal(c); err != nil {
		return []byte("")
	} else {
		return data
	}
}
