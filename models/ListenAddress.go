package models

import (
	"net/http"
	"strings"
)

// ListenAddress describes the address to listen on. For each Address a specific httpServer is started.
// It must be specified in a config file before it can be used

// ListenAddress config for ports
type ListenAddress struct {
	Address  string
	SSL      bool
	Task     InterfaceTask
	TaskData TaskData
}

// TaskData data for interface Task
type TaskData struct {
	Redirect RedirectData
}

// RedirectData data for interface Task to redirect
type RedirectData struct {
	Body     string
	HTTPCode int
}

// InterfaceTask task for an Address interface
type InterfaceTask string

// ...
const (
	HTTPRedirectTask InterfaceTask = "httpredirect"
	ProxyTask        InterfaceTask = "proxy"
)

// GetTask gets task from AddressInterface. If not set, return Default task
func (address *ListenAddress) GetTask() InterfaceTask {
	// Return default task if not set
	if string(address.Task) == "" {
		return ProxyTask
	}

	return address.Task
}

// GetBody returns body. If empty return default body
func (redirectData RedirectData) GetBody() string {
	if len(redirectData.Body) == 0 {
		return "Moved permanently"
	}
	return redirectData.Body
}

// GetHTTPCode returns body. If empty return default body
func (redirectData RedirectData) GetHTTPCode() int {
	if len(redirectData.Body) == 0 {
		return http.StatusMovedPermanently
	}
	return redirectData.HTTPCode
}

// GetAddress returns address of a listenAddress
func (address ListenAddress) GetAddress() string {
	return address.Address
}

// GetPort returns port of address
func (address ListenAddress) GetPort() string {
	return strings.Split(address.Address, ":")[1]
}
