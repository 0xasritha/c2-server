package main

import (
	"github.com/shirou/gopsutil/host"
)

type SystemInfo struct {
	Hostname             string `json:"hostname"`              // System hostname
	HostID               string `json:"host_id"`               // Unique machine ID (e.g., hashed MAC or kernel host ID)
	KernelVersion        string `json:"kernel_version"`        // Kernel version string (e.g., 6.1.0-17-amd64)
	KernelArch           string `json:"kernel_arch"`           // Architecture (e.g., amd64, arm64)
	VirtualizationSystem string `json:"virtualization_system"` // Virtualization type (e.g., kvm, qemu, vmware)
	VirtualizationRole   string `json:"virtualization_role"`   // Role: "guest" or "host"
}

func GetSystemInfo() SystemInfo {
	info, err := host.Info()
	if err != nil {
		panic(err) // TODO: standardize error handling
	}

	return SystemInfo{
		Hostname:             info.Hostname,
		HostID:               info.HostID,
		KernelVersion:        info.KernelVersion,
		KernelArch:           info.KernelArch,
		VirtualizationSystem: info.VirtualizationSystem,
		VirtualizationRole:   info.VirtualizationRole,
	}
}
