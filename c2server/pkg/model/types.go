package model

import "gorm.io/gorm"

type TaskType int

const (
	Configure TaskType = iota
	Ping
	ShellCommand
	FileTransfer
)

func (t TaskType) String() string {
	names := []string{"Configure", "Ping", "ShellCommand", "FileTransfer"}
	if int(t) < len(names) {
		return names[t]
	}
	return "UNKNOWN"
}

type SystemInfo struct {
	Hostname             string   `json:"hostname"`
	HostID               string   `json:"host_id"`
	KernelVersion        string   `json:"kernel_version"`
	KernelArch           string   `json:"kernel_arch"`
	VirtualizationSystem string   `json:"virtualization_system"`
	VirtualizationRole   string   `json:"virtualization_role"`
	IPs                  []string `json:"ips,omitempty" gorm:"type:text[]"`
}

// TODO: check over the new 1-1 association w/ Task + Result -> check over implant logic as well
type Implant struct {
	gorm.Model
	SystemInfo

	Dwell        string  `json:"dwell" gorm:"default:'5m'"`
	Jitter       float64 `json:"jitter" gorm:"default:0"`
	CommProtocol string  `json:"comm_protocol" gorm:"default:HTTPS"` // enum from listener.CommProtocol

	Tasks []Task `gorm:"foreignKey:ImplantID;constraint:OnDelete:CASCADE;"`
}

type Task struct {
	gorm.Model
	TaskType  `json:"type"`
	Options   map[string]interface{} `json:"options" gorm:"serializer:json;type:jsonb"`
	Completed bool                   `json:"completed"`
	Result    `json:"result" gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE;"`
	ImplantID uint `json:"implant_id" gorm:"index"`
}

type Result struct {
	gorm.Model
	Content string `json:"content"`
	Success bool   `json:"success"`
	TaskID  uint   `json:"task_id" gorm:"index"`
}
