package docker

type Option func(*Docker)

func WithPersistContainers() Option {
	return func(d *Docker) {
		d.persistContainers = true
	}
}
