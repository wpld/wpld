package specs

import (
	"fmt"

	"wpld/internal/entities"
)

func NewDatabaseSpec(db, volume string) entities.Specification {
	return entities.Specification{
		Name:  "MySQL",
		Image: "mariadb:latest",
		Volumes: []string{
			fmt.Sprintf("%s:/var/lib/mysql", volume),
		},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "password",
			"MYSQL_DATABASE":      db,
			"MYSQL_USER":          "wordpress",
			"MYSQL_PASSWORD":      "password",
		},
	}
}
