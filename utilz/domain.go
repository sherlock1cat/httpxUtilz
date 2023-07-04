package utilz

import (
	"errors"
	"log"
	"strings"
)

func GetSubDomain(url string) (domain string, err error) {
	domain = url
	// Extract the domain portion.
	if strings.Contains(url, "://") {
		domainParts := strings.Split(url, "://")
		if len(domainParts) != 2 {
			log.Println("GetSubDomain: Invalid domain name")
			err := errors.New("GetSubDomain: Invalid domain name")
			return "", err
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

	return domain, nil
}
