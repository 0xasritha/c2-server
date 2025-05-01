package main

import (
	"gorm.io/gorm"
)

type CommProtocol int

const (
	HTTPS CommProtocol = iota
	DNS
)

func (c CommProtocol) String() string {
	return [...]string{"HTTPS", "DNS"}[c-1]
}

type Metadata struct {
	Hostname             string   `json:"hostname"`              // System hostname
	HostID               string   `json:"host_id"`               // Unique machine ID (e.g., hashed MAC or kernel host ID)
	KernelVersion        string   `json:"kernel_version"`        // Kernel version string (e.g., 6.1.0-17-amd64)
	KernelArch           string   `json:"kernel_arch"`           // Architecture (e.g., amd64, arm64)
	VirtualizationSystem string   `json:"virtualization_system"` // Virtualization type (e.g., kvm, qemu, vmware)
	VirtualizationRole   string   `json:"virtualization_role"`   // Role: "guest" or "host"
	CommProtocol         string   `json:"comm_protocol"`         // Communication protocol (e.g., http, https, dns) // TODO: should this be type string or CommProtocol
	IPs                  []string `json:"ips,omitempty"`         // (Optional) Last known IP addresses // TODO: implement
}

// TODO: remove unneeded fields
type Implant struct {
	gorm.Model // TODO: i want ID to be a UUID
	Metadata   // Server metadata ?? #TODO: what other fields should be in metadata?
	// 	LastSeen time.Time // Timestamp of last beacon or activity
	// PublicKey      []byte    // If using asymmetric encryption
	// NonceCounter   uint64    // To prevent replay attacks
	// SessionToken   string    // JWT or session ID
	PendingTasks   []Task // Pending tasks for this implant
	CompletedTasks []Task // History or recent completed tasks
	//	Labels         []string  // Optional: e.g., ["phishing", "test", "persistent"]
	// UserAgent string // Optional: helps with fingerprinting
}
