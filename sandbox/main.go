package ipparser

import (
	"log/slog"
	"net"
	"regexp"
	"strings"
)

func ParseCidrs(input string) ([]*net.IPNet, error) {
	if strings.TrimSpace(input) == "" {
		return nil, nil
	}

	ipv4Regex := `\b(?:\d{1,3}\.){3}\d{1,3}(?:/\d{1,2})?\b`
	ipv6Regex := `\b(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}(?:/\d{1,3})?\b|` +
		`\b(?:[0-9a-fA-F]{1,4}:){1,7}:(?:/\d{1,3})?\b|` +
		`\b:(?::[0-9a-fA-F]{1,4}){1,7}(?:/\d{1,3})?\b|` +
		`\b(?:[0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}(?:/\d{1,3})?\b`

	combined := regexp.MustCompile(ipv4Regex + "|" + ipv6Regex)
	matches := combined.FindAllString(input, -1)

	var result []*net.IPNet
	for _, match := range matches {
		var ipnet *net.IPNet
		var err error

		if !strings.Contains(match, "/") {
			ip := net.ParseIP(match)
			if ip == nil {
				continue
			}
			if ip.To4() != nil {
				match += "/32"
				slog.Info("is ipv4", "address", match)
			} else {
				slog.Info("is ipv6", "address", match)
				match += "/128"
			}
		}

		_, ipnet, err = net.ParseCIDR(match)
		if err != nil {
			return nil, err
		}

		if ipnet != nil {
			result = append(result, ipnet)
		}
	}

	return result, nil
}
