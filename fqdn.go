package main

import "strings"

type FQDN struct {
	TLD        string
	ZoneName   string
	Subdomains []string
}

func ParseFQDN(s string) FQDN {
	parts := strings.Split(s, ".")

	tld := parts[len(parts)-1]
	zoneName := parts[len(parts)-2]

	var subdomains []string

	if len(parts) > 2 {
		for i := len(parts) - 3; i >= 0; i-- {
			subdomains = append(subdomains, parts[i])
		}
	}

	return FQDN{
		TLD:        tld,
		ZoneName:   zoneName,
		Subdomains: subdomains,
	}
}
