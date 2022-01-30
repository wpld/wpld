package entities

import (
	"errors"
	"fmt"

	mapset "github.com/deckarep/golang-set"
)

type Project struct {
	ID       string                   `yaml:"id"`
	Name     string                   `yaml:"name"`
	WP       *WordPress               `yaml:"wordpress,omitempty"`
	Volumes  []string                 `yaml:"volumes"`
	Services map[string]Specification `yaml:"services"`
	Scripts  map[string]Script        `yaml:"scripts"`
}

func (p Project) GetNetwork() Network {
	return Network{
		Name: p.ID,
	}
}

func (p Project) GetContainerIDForService(serviceID string) string {
	return fmt.Sprintf("%s__%s", p.ID, serviceID)
}

func (p Project) GetServices() ([]Service, error) {
	index := 0
	servicesLen := len(p.Services)

	services := make([]Service, servicesLen)
	deps := make(map[string]mapset.Set, servicesLen)

	for id, spec := range p.Services {
		deps[id] = mapset.NewSet()
		for _, dep := range spec.DependsOn {
			deps[id].Add(dep)
		}
	}

	for len(deps) > 0 {
		ready := mapset.NewSet()

		for id, dep := range deps {
			if dep.Cardinality() == 0 {
				ready.Add(id)
				delete(deps, id)

				containerID := p.GetContainerIDForService(id)
				service := Service{
					ID:      containerID,
					Network: p.GetNetwork(),
					Project: p.Name,
					Spec:    p.Services[id],
					Aliases: []string{
						id,
						containerID,
					},
				}

				for i, from := range service.Spec.VolumesFrom {
					if _, ok := p.Services[from]; ok {
						service.Spec.VolumesFrom[i] = p.GetContainerIDForService(from)
					}
				}

				services[index] = service
				index++
			}
		}

		if ready.Cardinality() == 0 {
			return nil, errors.New("circular dependencies found")
		} else {
			for id, dep := range deps {
				deps[id] = dep.Difference(ready)
			}
		}
	}

	return services, nil
}
