package config

// Application is the configuration for the application.
// Holds flags, parameters.
var Application = struct {
	Hosts []string
	// /var/run/docker.sock:/docker/:*,-POST,-PUT,-DELETE
	Sockets []string
}{
	Hosts: []string{"0.0.0.0:8080"},
}
