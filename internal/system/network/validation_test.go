package network

import "testing"

func TestValidateAddress(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "IPv4 CIDR preserves host", input: "192.168.1.10/24", want: "192.168.1.10/24"},
		{name: "IPv6 CIDR preserves host", input: "2001:db8::1/64", want: "2001:db8::1/64"},
		{name: "whitespace", input: " 192.168.1.10/24 \n", want: "192.168.1.10/24"},
		{name: "bare IPv4", input: "192.168.1.10", wantErr: true},
		{name: "bare IPv6", input: "2001:db8::1", wantErr: true},
		{name: "empty", input: " ", wantErr: true},
		{name: "missing prefix", input: "192.168.1.10/", wantErr: true},
		{name: "invalid IPv4 prefix", input: "192.168.1.10/33", wantErr: true},
		{name: "invalid IPv6 prefix", input: "2001:db8::1/129", wantErr: true},
		{name: "invalid address", input: "999.168.1.10/24", wantErr: true},
		{name: "dotted mask", input: "192.168.1.10/255.255.255.0", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateAddress(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateAddress(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("validateAddress(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
