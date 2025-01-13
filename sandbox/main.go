package main

import (
	"errors"
	"net"
	"regexp"
)

func findCIDRs(input string) ([]*net.IPNet, error) {
	// Regular expression for CIDR IPv4 or IPv6 matches
	cidrRegex := `((?:\d{1,3}\.){3}\d{1,3}\/\d{1,2})|(?:[a-fA-F0-9:]+\/\d{1,3})`

	compiledRegex, err := regexp.Compile(cidrRegex)
	if err != nil {
		return nil, errors.New("failed to compile regex")
	}

	// Find all matches in the input
	matches := compiledRegex.FindAllString(input, -1)

	parsedCIDRs := []*net.IPNet{}

	for _, match := range matches {
		_, subnet, err := net.ParseCIDR(match)
		if err != nil {
			// Skip invalid CIDRs
			continue
		}
		parsedCIDRs = append(parsedCIDRs, subnet)
	}

	return parsedCIDRs, nil
}
