package main

import (
	"net"
	"testing"
)

func TestFindCIDRs(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCIDRs []string
		expectError   bool
	}{
		{
			name:          "Single IPv4 CIDR",
			input:         "The IP range 192.168.1.0/24 is used for internal networking.",
			expectedCIDRs: []string{"192.168.1.0/24"},
			expectError:   false,
		},
		{
			name:          "Multiple mixed CIDRs",
			input:         "Valid CIDRs: 10.0.0.0/8, 172.16.0.0/12, and fe80::/10",
			expectedCIDRs: []string{"10.0.0.0/8", "172.16.0.0/12", "fe80::/10"},
			expectError:   false,
		},
		{
			name:          "Non-CIDR IP addresses ignored",
			input:         "Here are some IPs: 1.1.1.1 and 2001:db8::ff00:42:8329 but no CIDRs.",
			expectedCIDRs: []string{},
			expectError:   false,
		},
		{
			name:          "Embedded CIDR in random text",
			input:         "Random text includes stuff like 127.0.0.1 and 192.168.100.0/24 interspersed.",
			expectedCIDRs: []string{"192.168.100.0/24"},
			expectError:   false,
		},
		{
			name:          "Invalid CIDRs ignored",
			input:         "Some invalid CIDRs: 300.300.0.0/24, and fe80:::/130",
			expectedCIDRs: []string{},
			expectError:   false,
		},
		{
			name:          "Empty input",
			input:         "",
			expectedCIDRs: []string{},
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundCIDRs, err := findCIDRs(tt.input)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Validate the CIDRs match
			expected := make([]*net.IPNet, len(tt.expectedCIDRs))
			for i, cidr := range tt.expectedCIDRs {
				_, subnet, err := net.ParseCIDR(cidr)
				if err != nil {
					t.Fatalf("unexpected error parsing expected CIDR: %v", err)
				}
				expected[i] = subnet
			}

			if len(foundCIDRs) != len(expected) {
				t.Errorf("expected %v CIDRs, but got %v", len(expected), len(foundCIDRs))
			}

			for _, found := range foundCIDRs {
				matched := false
				for _, exp := range expected {
					if found.String() == exp.String() {
						matched = true
						break
					}
				}
				if !matched {
					t.Errorf("unexpected CIDR found: %v", found)
				}
			}
		})
	}
}
