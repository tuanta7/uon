package system

import (
	"io"
	"strings"
	"testing"
)

func TestIPCommandRejectsInvalidOptions(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"missing interface", []string{"--dhcp"}, "--interface is required"},
		{"blank interface", []string{"--interface", " ", "--dhcp"}, "--interface is required"},
		{"missing mode", []string{"--interface", "eth0"}, "provide --static or --dhcp"},
		{"disabled DHCP", []string{"--interface", "eth0", "--dhcp=false"}, "provide --static or --dhcp"},
		{"conflicting modes", []string{"--interface", "eth0", "--static", "192.168.1.10/24", "--dhcp"}, "none of the others can be"},
		{"DHCP with gateway", []string{"--interface", "eth0", "--dhcp", "--gateway", "192.168.1.1/24"}, "none of the others can be"},
		{"gateway without static", []string{"--interface", "eth0", "--gateway", "192.168.1.1/24"}, "provide --static or --dhcp"},
		{"positional argument", []string{"eth0"}, "unknown command"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newIPCommand()
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)
			cmd.SetArgs(tt.args)
			if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Execute() error = %v, want error containing %q", err, tt.want)
			}
		})
	}
}
