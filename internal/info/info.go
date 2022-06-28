package info

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
