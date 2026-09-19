package netplan

type Config struct {
	Network Network `yaml:"network"`
}

type Network struct {
	Version   int                  `yaml:"version"`
	Ethernets map[string]Interface `yaml:"ethernets"`
}

type Interface struct {
	DHCP4     bool     `yaml:"dhcp4"` // true = DHCP
	DHCP6     bool     `yaml:"dhcp6,omitempty"`
	Addresses []string `yaml:"addresses,omitempty"`
	Routes    []Route  `yaml:"routes,omitempty"`
}

type Route struct {
	To  string `yaml:"to"`
	Via string `yaml:"via"`
}
