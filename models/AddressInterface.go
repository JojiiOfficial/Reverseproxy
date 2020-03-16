package models

import "net/http"

// AddressInterface an interface for address
type AddressInterface struct {
	Address  string
	Task     InterfaceTask
	TaskData TaskData
}

// TaskData data for interface Task
type TaskData struct {
	Redirect RedirectData
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

// RedirectData data for interface Task to redirect
type RedirectData struct {
	Location string
	Body     string
	HTTPCode int
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
