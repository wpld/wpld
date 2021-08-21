package specs

import (
	"fmt"

	"wpld/internal/entities"
)

const (
	PROJECT_TYPE_PLUGIN     = "A plugin"
	PROJECT_TYPE_THEME      = "A theme"
	PROJECT_TYPE_WP_CONTENT = "Whole wp-content folder"
)

func NewWordPressSpec(slug, db, volume, php, projectType string) entities.Specification {
	volumes := []string{
		fmt.Sprintf("%s:/var/www/html", volume),
	}

	switch projectType {
	case PROJECT_TYPE_PLUGIN:
		volumes = append(volumes, fmt.Sprintf("./:/var/www/html/wp-content/plugins/%s", slug))
	case PROJECT_TYPE_THEME:
		volumes = append(volumes, fmt.Sprintf("./:/var/www/html/wp-content/themes/%s", slug))
	case PROJECT_TYPE_WP_CONTENT:
		volumes = append(volumes, "./:/var/www/html/wp-content/")
	}

	return entities.Specification{
		Name:    "WordPress",
		Image:   fmt.Sprintf("wordpress:5-php%s-fpm-alpine", php),
		Volumes: volumes,
		Env: map[string]string{
			"WORDPRESS_DB_HOST":      "db",
			"WORDPRESS_DB_USER":      "wordpress",
			"WORDPRESS_DB_PASSWORD":  "password",
			"WORDPRESS_DB_NAME":      db,
			"WORDPRESS_DEBUG":        "on",
			"WORDPRESS_CONFIG_EXTRA": "define( \"WP_DEBUG_DISPLAY\", false );\ndefine( \"WP_DEBUG_LOG\", true );",
		},
	}
}
