package docker

type ContainerSchema struct {
	ID              string                   `json:"Id"`
	Names           []string                 `json:"Names"`
	Image           string                   `json:"Image"`
	ImageID         string                   `json:"ImageID"`
	Command         string                   `json:"Command"`
	Created         int64                    `json:"Created"`
	Ports           []PortSchema             `json:"Ports"`
	Labels          map[string]string        `json:"Labels"`
	State           ContainerState           `json:"State"`
	Status          string                   `json:"Status"`
	HostConfig      ContainerHostConfig      `json:"HostConfig"`
	NetworkSettings ContainerNetworkSettings `json:"NetworkSettings"`
	Mounts          []ContainerMount         `json:"Mounts"`
	SizeRw          int64                    `json:"SizeRw"`
	SizeRootFs      int64                    `json:"SizeRootFs"`
}

type ContainerHostConfig struct {
	NetworkMode string `json:"NetworkMode"`
}

type ContainerNetworkSettings struct {
	Networks map[string]ContainerSingleNetworkSettings `json:"Networks"`
}
type ContainerSingleNetworkSettings struct {
	NetworkMode         string `json:"NetworkID"`
	EndpointID          string `json:"EndpointID"`
	Gateway             string `json:"Gateway"`
	IPAddress           string `json:"IPAddress"`
	IPPrefixLen         int    `json:"IPPrefixLen"`
	IPv6Gateway         string `json:"IPv6Gateway"`
	GlobalIPv6Address   string `json:"GlobalIPv6Address"`
	GlobalIPv6PrefixLen int    `json:"GlobalIPv6PrefixLen"`
	MacAddress          string `json:"MacAddress"`
}

type PortSchema struct {
	IP          string `json:"IP"`
	PrivatePort int    `json:"PrivatePort"`
	PublicPort  int    `json:"PublicPort"`
	Type        string `json:"Type"`
}

type ContainerMount struct {
	Name        string `json:"Name"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Driver      string `json:"Driver"`
	Mode        string `json:"Mode"`
	RW          bool   `json:"RW"`
	Propagation string `json:"Propagation"`
}

type ContainerState string

const (
	StateCreated    ContainerState = "created"
	StateRunning    ContainerState = "running"
	StateRestarting ContainerState = "restarting"
	StateExited     ContainerState = "exited"
	StatePaused     ContainerState = "paused"
	StateDead       ContainerState = "dead"
)
