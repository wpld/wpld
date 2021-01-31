package cases

import (
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"os"
)

func GetProjectConfig() (*viper.Viper, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	config := viper.New()
	config.SetFs(afero.NewOsFs())
	config.SetConfigName("wpld")
	config.SetConfigType("yaml")
	config.AddConfigPath(dir)
	if err = config.ReadInConfig(); err != nil {
		return nil, err
	}

	return config, nil
}
