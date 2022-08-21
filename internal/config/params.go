package config

// Serve is the configuration for the net connection.
type Serve struct {
	// /var/run/docker.sock:/docker/:*,-POST,-PUT,-DELETE
	Socket []string `cfg:"socket"`
}

// Application is the configuration for the application.
// Holds flags, parameters.
var Application = struct {
	LogLevel string `cfg:"logLevel" env:"LOG_LEVEL"`
	Router   Router `cfg:"router"`
}{
	LogLevel: "info",
}
