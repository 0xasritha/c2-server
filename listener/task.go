package main

import "gorm.io/gorm"

type TaskType int

const (
	Configure TaskType = iota
	Ping
	ListThreads
	Execute
)

func (t TaskType) String() string {
	return [...]string{"Configure", "Ping", "ListThreads", "Execute", ""}[t-1]
}

type Task struct {
	gorm.Model
	TaskType string `json:"type"`
	// completed or done?
	completed bool `json:"completed"`

	// TODO: which one to use?
	Options map[string]interface{} `json:"options" gorm:"serializer:json;type:jsonb""`
	// Options datatypes.JSONMap      `gorm:"type:jsonb"`
	Results []Result `json:"results"` //TODO: Results can be null
}
