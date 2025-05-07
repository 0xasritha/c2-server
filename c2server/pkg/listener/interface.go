package listener

// TODO: do we want to maintain CommProtocol -> is this even being used anywere? if get rid it, then remove in implant as well
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
}
