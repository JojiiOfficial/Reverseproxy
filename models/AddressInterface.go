package models

// AddressInterface an interface for address
type AddressInterface struct {
	Address string
	Task    InterfaceTask
}

// InterfaceTask task for an Address interface
type InterfaceTask string

// ...
const (
	HTTPRedirectTask InterfaceTask = "httpredirect"
	ProxyTask        InterfaceTask = "proxy"
)

// GetTask gets task from AddressInterface. If not set, return Default task
func (aif *AddressInterface) GetTask() InterfaceTask {
	// Return default task if not set
	if string(aif.Task) == "" {
		return ProxyTask
	}

	return aif.Task
}
