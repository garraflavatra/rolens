package app

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ExecuteShellScriptResult struct {
	Output           string `json:"output"`
	Status           int    `json:"status"`
	ErrorTitle       string `json:"errorTitle"`
	ErrorDescription string `json:"errorDescription"`
}

func (a *App) ExecuteShellScript(hostKey, dbKey, collKey, script string) (result ExecuteShellScriptResult) {
	if !a.Env.HasMongoShell {
		result.ErrorTitle = "mongosh not found"
		result.ErrorDescription = "The mongosh executable is required to run a shell script. Please see https://www.mongodb.com/docs/mongodb-shell/install/"
		return
	}

	hosts, err := a.Hosts()
	if err != nil {
		result.ErrorTitle = "Could not get hosts"
		result.ErrorDescription = err.Error()
		return
	}

	host, hostFound := hosts[hostKey]
	if !hostFound {
		result.ErrorTitle = "The specified host does not seem to exist"
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Shell: failed to generate a UUID: %s", err.Error())
		result.ErrorTitle = "Could not generate UUID"
		result.ErrorDescription = err.Error()
		return
	}

	dirname := path.Join(a.Env.DataDirectory, "Shell Scripts")
	fname := path.Join(dirname, fmt.Sprintf("%s.js", id.String()))

	if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
		runtime.LogWarningf(a.ctx, "Shell: failed to mkdir %s", err.Error())
		result.ErrorTitle = "Could not create temporary directory"
		result.ErrorDescription = err.Error()
		return
	}

	connstr := host.URI
	if !strings.HasSuffix(connstr, "/") {
		connstr = connstr + "/"
	}

	connstr = connstr + dbKey
	script = fmt.Sprintf("db = connect('%s');\ncoll = db.getCollection('%s');\n\n%s", connstr, collKey, script)

	if err := os.WriteFile(fname, []byte(script), os.ModePerm); err != nil {
		runtime.LogWarningf(a.ctx, "Shell: failed to write script to %s", err.Error())
		result.ErrorTitle = "Could not create temporary script file"
		result.ErrorDescription = err.Error()
		return
	}

	cmd := exec.Command("mongosh", "--file", fname, connstr)
	stdout, err := cmd.Output()

	if exiterr, ok := err.(*exec.ExitError); ok {
		result.Status = exiterr.ExitCode()
	} else if err != nil {
		runtime.LogWarningf(a.ctx, "Shell: failed to execute: mongosh --file %s: %s", fname, err.Error())
		result.ErrorTitle = "Could not execute script"
		result.ErrorDescription = err.Error()
		return
	} else {
		result.Status = 0
	}

	result.Output = string(stdout)
	return
}
