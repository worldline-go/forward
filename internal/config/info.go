package config

import "time"

// Info holds the information of build.
type Info struct {
	StartDate   time.Time `json:"startDate"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	BuildCommit string    `json:"buildCommit"`
	BuildDate   string    `json:"buildDate"`
}

// AppInfo is hold info about build date, version, commit.
var AppInfo = Info{
	StartDate: time.Now(),
	Name:      "forward",
	Version:   "v0.0.0",
}

// Prefix for vault and consul configuration.
type Prefix struct {
	Vault  string `cfg:"vault"`
	Consul string `cfg:"consul"`
}

// Load for config loader.
type Load struct {
	Prefix  Prefix `cfg:"prefix"`
	AppName string `cfg:"appName" env:"APP_NAME"`
}

// LoadConfig is the configuration for the config loader.
var LoadConfig = Load{
	AppName: AppInfo.Name,
}
