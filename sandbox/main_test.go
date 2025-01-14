package ipparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCidrs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		wantErr  bool
	}{
		{
			name:     "empty input",
			input:    "",
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "single IPv4 CIDR",
			input:    "The network is 192.168.1.0/24",
			expected: []string{"192.168.1.0/24"},
			wantErr:  false,
		},
		{
			name:     "single IPv6 CIDR",
			input:    "IPv6 network: 2001:db8::/32",
			expected: []string{"2001:db8::/32"},
			wantErr:  false,
		},
		{
			name:     "multiple mixed CIDRs",
			input:    "Networks: 10.0.0.0/8 and 2001:db8::/32 and 172.16.0.0/12",
			expected: []string{"10.0.0.0/8", "2001:db8::/32", "172.16.0.0/12"},
			wantErr:  false,
		},
		{
			name:     "single IPs should be converted to /32 or /128",
			input:    "IPs: 192.168.1.1 and 2001:db8::1",
			expected: []string{"192.168.1.1/32", "2001:db8::1/128"},
			wantErr:  false,
		},
		{
			name:     "invalid IP or CIDR",
			input:    "Invalid: 256.256.256.256/24",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseCidrs(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, result, len(tt.expected))

			// Convert results to strings for comparison
			resultStrings := make([]string, len(result))
			for i, ipnet := range result {
				resultStrings[i] = ipnet.String()
			}

			assert.ElementsMatch(t, tt.expected, resultStrings)
		})
	}
}
