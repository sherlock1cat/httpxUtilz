package utilz

import (
	"log"
	"strings"
)

func GetSubDomain(url string) (domain string) {
	domain = url
	// Extract the domain portion.
	if strings.Contains(url, "://") {
		domainParts := strings.Split(url, "://")
		if len(domainParts) != 2 {
			log.Fatal("GetSubDomain: Invalid domain name")
			return
		}
		domain = domainParts[1]
	}

	// If the domain portion contains a port, further extraction is required.
	if strings.Contains(domain, ":") {
		domainParts := strings.Split(domain, ":")
		domain = domainParts[0]
	}

	// If the domain portion contains a slash, further extraction is required.
	if strings.Contains(domain, "/") {
		domainParts := strings.Split(domain, "/")
		domain = domainParts[0]
	}

	return
}
