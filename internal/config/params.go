package config

// Serve is the configuration for the net connection.
type Serve struct {
	// /var/run/docker.sock:/docker/:*,-POST,-PUT,-DELETE
	Socket []string `json:"socket"`
}

// Application is the configuration for the application.
// Holds flags, parameters.
var Application = struct {
	Host      string `json:"host"`
	LogLevel  string `json:"logLevel"`
	LogFormat string `json:"logFormat"`
	Serve     Serve  `json:"serve"`
}{
	Host:      "0.0.0.0:8080",
	LogLevel:  "info",
	LogFormat: "common",
}
