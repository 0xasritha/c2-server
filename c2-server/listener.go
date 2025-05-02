package main

import (
	"fmt"
	"gorm.io/gorm"
)

type Listener interface {
	listen()
}

func NewListener(db *gorm.DB, commProtocol CommProtocol) (Listener, error) {
	switch commProtocol {
	case HTTPS:
		return &HTTPSListener{db: db}, nil
	case DNS:
		return &DNSListener{}, nil
	}
	return nil, fmt.Errorf("invalid communication protocol: %s", commProtocol)
}
