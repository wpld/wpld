package specs

import (
	"wpld/internal/entities"
)

func NewMemcachedSpec() entities.Specification {
	return entities.Specification{
		Name:  "Memcached",
		Image: "memcached:alpine",
	}
}
