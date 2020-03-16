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
	if location.Deny == "all" && len(location.Allow) == 0 {
		return false
	} else if location.Deny == "all" {
		// To IP
		sIP := strings.Split(req.RemoteAddr, ":")[0]
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

		// Otherwise return false -> deny all
		return false
	}

	// By default any request is allowed
	return true
}
