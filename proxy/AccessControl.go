package proxy

import (
	"net"
	"net/http"
	"strings"

	"github.com/JojiiOfficial/ReverseProxy/models"
	log "github.com/sirupsen/logrus"
)

// Return true if request is not denied by some rules or IP has access
func isRequestAllowed(req *http.Request, location *models.RouteLocation) bool {
	// if location ist denied by 'all' and no exception is specified, return 'not allowed'
	if location.HasDenyRoule && len(location.Allow) == 0 {
		return false
	} else if location.HasDenyRoule {
		// To IP
		sourceHost := req.RemoteAddr
		// Use a custom header if specified
		if len(location.SrcIPHeader) > 0 {
			sourceHost = req.Header.Get(location.SrcIPHeader)
			// If header is empty, warn and return 'not allowed'
			if len(sourceHost) == 0 {
				log.Warn("SrcIPHeader is empty!")
				return false
			}
		}

		// Remove port and parse to net.IP
		sIP := strings.Split(sourceHost, ":")[0]
		ip := net.ParseIP(sIP)

		// Loop allowed IPs
		for _, allowedIPr := range location.Allow {
			// if specified IP is a cidr parse it otherwise compare it
			if strings.Contains(allowedIPr, "/") && ip.To4() != nil {
				_, cidr, err := net.ParseCIDR(allowedIPr)
				if err != nil {
					log.Error(err)
					return false
				}

				// return "allowed" if requesting IP is in cidr range
				if cidr.Contains(ip.To4()) {
					return true
				}
			} else {
				// return "allowed" if ip matches with specified ips
				if allowedIPr == sIP {
					return true
				}
			}
		}

		// Otherwise return 'deny all'
		return false
	}

	// By default any request is allowed
	return true
}
