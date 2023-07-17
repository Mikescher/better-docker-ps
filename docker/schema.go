package docker

import (
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"net"
)

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

func (s ContainerSchema) PortsSorted() []PortSchema {
	ports := langext.ArrCopy(s.Ports)

	langext.SortSliceStable(ports, func(p1, p2 PortSchema) bool {
		if p1.PublicPort != p2.PublicPort {
			return p1.PublicPort < p2.PublicPort
		}
		if p1.PrivatePort != p2.PrivatePort {
			return p1.PrivatePort < p2.PrivatePort
		}
		return false
	})

	return ports
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

func (s PortSchema) IsLoopback() bool {
	ip := net.ParseIP(s.IP)
	return ip != nil && ip.IsLoopback()
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

func (ct ContainerState) Num() int {
	switch ct {
	case StateCreated:
		return 0
	case StateRunning:
		return 1
	case StateRestarting:
		return 2
	case StateExited:
		return 3
	case StatePaused:
		return 4
	case StateDead:
		return 5
	default:
		return 999
	}
}
