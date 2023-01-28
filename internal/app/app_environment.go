package app

import (
	"os/exec"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type EnvironmentInfo struct {
	Arch      string `json:"arch"`
	BuildType string `json:"buildType"`
	Platform  string `json:"platform"`

	HasMongoExport bool `json:"hasMongoExport"`
	HasMongoDump   bool `json:"hasMongoDump"`
}

var env EnvironmentInfo
var envKnown = false

func (a *App) Environment() EnvironmentInfo {
	if !envKnown {
		wailsEnv := runtime.Environment(a.ctx)
		env.Arch = wailsEnv.Arch
		env.BuildType = wailsEnv.BuildType
		env.Platform = wailsEnv.Platform

		_, err := exec.LookPath("mongodump")
		env.HasMongoDump = err == nil

		_, err = exec.LookPath("mongoexport")
		env.HasMongoExport = err == nil

		envKnown = true
	}
	return env
}
