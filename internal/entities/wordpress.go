package entities

type WordPress struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Email        string `yaml:"email"`
	Multisite    bool   `yaml:"multisite"`
	Subdirectory bool   `yaml:"subdirectory"`
}
