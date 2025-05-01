package main

import "gorm.io/gorm"

type Result struct {
	gorm.Model
	Content *string `json:"content"`
	Success bool    `json:"success"`
	TaskID  uint    `json:"task_id"` // foreign key
}
