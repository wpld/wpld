package specs

import (
	"fmt"

	"wpld/internal/entities"
)

func NewWordPressSpec(db, volume, php string) entities.Specification {
	return entities.Specification{
		Name:  "WordPress",
		Image: fmt.Sprintf("wordpress:5-php%s-fpm-alpine", php),
		Volumes: []string{
			fmt.Sprintf("%s:/var/www/html", volume),
		},
		Env: map[string]string{
			"WORDPRESS_DB_HOST":     "db",
			"WORDPRESS_DB_USER":     "wordpress",
			"WORDPRESS_DB_PASSWORD": "password",
			"WORDPRESS_DB_NAME":     db,
		},
	}
}
