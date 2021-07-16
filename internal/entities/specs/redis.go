package specs

import (
	"wpld/internal/entities"
)

func NewRedisSpec() entities.Specification {
	return entities.Specification{
		Name:  "Redis",
		Image: "redis:alpine",
	}
}
