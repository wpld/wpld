package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
)

func NetworkDebugInfo(net types.NetworkResource) map[string]interface{} {
	return map[string]interface{}{
		"id":      net.ID,
		"scope":   net.Scope,
		"driver":  net.Driver,
		"ipam":    net.IPAM,
		"options": net.Options,
	}
}

func VolumeDebugInfo(vol volume.Volume) map[string]interface{} {
	return map[string]interface{}{
		"scope":      vol.Scope,
		"driver":     vol.Driver,
		"mountpoint": vol.Mountpoint,
		"status":     vol.Status,
		"options":    vol.Options,
	}
}
