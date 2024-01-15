package config

type PortForward struct {
	Type         string `json:"type" yaml:"type"`
	Port         int    `json:"port" yaml:"port"`
	RemoteDomain string `json:"remote_domain" yaml:"remote_domain"`
	RemotePort   int    `json:"remote_port" yaml:"remote_port"`
	Timing       int    `json:"timing" yaml:"timing"`
}
