package listener

import (
	"fmt"
	"gorm.io/gorm"
)

// TODO: move to types?
type CommProtocol int

const (
	HTTPS CommProtocol = iota
	DNS
)

func (c CommProtocol) String() string {
	names := []string{"HTTPS", "DNS"}
	if int(c) < len(names) {
		return names[c]
	}
	return "UNKNOWN"
}

type Listener interface {
	Listen()
	// register()
}

func New(db *gorm.DB, commProtocol CommProtocol) (Listener, error) {
	switch commProtocol {
	case HTTPS:
		return &HTTPSListener{db: db}, nil
	case DNS:
		return &DNSListener{db: db}, nil
	}
	return nil, fmt.Errorf("invalid communication protocol: %s", commProtocol)
}
