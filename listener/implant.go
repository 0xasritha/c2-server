package main

import (
	"gorm.io/gorm"
	"time"
)

type CommProtocol int

const (
	HTTPS CommProtocol = iota
	DNS
)

func (c CommProtocol) String() string {
	return [...]string{"HTTPS", "DNS"}[c-1]
}

// TODO: remove unneeded fields
type Implant struct {
	gorm.Model
	Hostname       string    // Optional: hostname of the system it's running on
	IP             string    // Last known IP address
	LastSeen       time.Time // Timestamp of last beacon or activity
	OS             string    // Operating system name/version
	Architecture   string    // e.g., amd64, arm64
	CommProtocol             // http, https, dns
	Version        string    // Implant version, for upgrade tracking
	Active         bool      // Implant status
	PublicKey      []byte    // If using asymmetric encryption
	NonceCounter   uint64    // To prevent replay attacks
	SessionToken   string    // JWT or session ID
	TaskQueue      []Task    // Pending tasks for this implant
	CompletedTasks []Task    // History or recent completed tasks
	//	Labels         []string  // Optional: e.g., ["phishing", "test", "persistent"]
	// UserAgent string // Optional: helps with fingerprinting
}
