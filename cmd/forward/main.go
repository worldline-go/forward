package main

import (
	"github.com/rakunlabs/into"
	"github.com/rakunlabs/logi"
	"github.com/worldline-go/forward/cmd/forward/args"
	"github.com/worldline-go/forward/internal/info"
)

var (
	version = "v0.0.0"
	commit  = "-"
	date    = "-"
)

func main() {
	// change information
	info.AppInfo.Version = version
	info.AppInfo.BuildCommit = commit
	info.AppInfo.BuildDate = date

	into.Init(
		args.Execute,
		into.WithLogger(logi.InitializeLog(logi.WithCaller(false))),
		into.WithStartFn(nil),
		into.WithStopFn(nil),
	)
}
