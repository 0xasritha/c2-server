package main

import (
	"fmt"
	"gorm.io/gorm"
)

type Listener interface {
	listen()
	// register()
}

func NewListener(db *gorm.DB, commProtocol CommProtocol) (Listener, error) {
	switch commProtocol {
	case HTTPS:
		return &HTTPSListener{db: db}, nil
	case DNS:
		return &DNSListener{db: db}, nil
	}
	return nil, fmt.Errorf("invalid communication protocol: %s", commProtocol)
}
