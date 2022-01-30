package entities

type Script struct {
	Service    string   `yaml:"service"`
	WorkingDir string   `yaml:"working_dir,omitempty"`
	Command    []string `yaml:"command"`
}
