package entities

import (
	"errors"
	"fmt"

	mapset "github.com/deckarep/golang-set"
)

type Project struct {
	ID       string                   `yaml:"id"`
	Name     string                   `yaml:"name"`
	Domains  []string                 `yaml:"domains"`
	Volumes  []string                 `yaml:"volumes"`
	Services map[string]Specification `yaml:"services"`
}

func (p Project) GetNetworkName() string {
	return p.ID
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

				service := Service{
					ID:      p.GetContainerIDForService(id),
					Alias:   id,
					Network: p.GetNetworkName(),
					Project: p.ID,
					Spec:    p.Services[id],
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
