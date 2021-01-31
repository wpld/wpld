package global

import (
	"database/sql"
	"errors"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
	"wpld/config"
	"wpld/models"
	"wpld/utils"
)

const (
	MYSQL_IMAGE_NAME     = "mysql:5"
	MYSQL_CONTAINER_NAME = "wpld_global_mysql"
)

func RunMySQL(factory models.DockerFactory, pull bool) error {
	if pull {
		img := factory.Image(MYSQL_IMAGE_NAME)
		if err := img.Pull(); err != nil {
			return err
		}
	}

	resources := container.Resources{}
	port := nat.PortBinding{
		HostIP:   "127.0.0.1",
		HostPort: "3306",
	}

	memory := viper.GetString(config.MYSQL_MEMORY)
	if mem, err := utils.ParseBytes(memory); err == nil {
		resources.Memory = mem
	}

	reservation := viper.GetString(config.MYSQL_RESERVATION)
	if reserve, err := utils.ParseBytes(reservation); err == nil {
		resources.MemoryReservation = reserve
	}

	if viper.IsSet(config.MYSQL_PORT) {
		port.HostPort = viper.GetString(config.MYSQL_PORT)
	}

	containerConfig := &container.Config{
		Image: MYSQL_IMAGE_NAME,
		//User: strconv.Itoa(os.Getuid()),
		Env: []string{
			"MYSQL_ROOT_PASSWORD=password",
		},
	}

	host := &container.HostConfig{
		NetworkMode: NETWORK_NAME,
		IpcMode:     "shareable",
		PortBindings: nat.PortMap{
			"3306/tcp": []nat.PortBinding{port},
		},
		Resources: resources,
	}

	mysql := factory.Container(MYSQL_CONTAINER_NAME)
	if err := mysql.Create(containerConfig, host); err != nil {
		return err
	}

	return mysql.Start()
}

func StopMySQL(factory models.DockerFactory, rm bool) error {
	mysql := factory.Container(MYSQL_CONTAINER_NAME)

	if rm {
		return mysql.Remove()
	}

	return mysql.Stop()
}

func WaitForMySQL() error {
	for i := 0; i < 12; i++ {
		db, err := sql.Open("mysql", "root:password@/information_schema")
		if err != nil {
			logrus.Error(err)
		} else {
			if pingErr := db.Ping(); pingErr == nil {
				_ = db.Close()
				return nil
			}
		}

		time.Sleep(5 * time.Second)
	}

	return errors.New("can't connect to MySQL")
}
