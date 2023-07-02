package app

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ExecuteShellScriptResult struct {
	Output           string `json:"output"`
	Stderr           string `json:"stderr"`
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
	fname := path.Join(dirname, fmt.Sprintf("%s.mongosh.js", id.String()))

	if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
		runtime.LogWarningf(a.ctx, "Shell: failed to mkdir %s", err.Error())
		result.ErrorTitle = "Could not create temporary directory"
		result.ErrorDescription = err.Error()
		return
	}

	scriptHeader := fmt.Sprintf("// Namespace: %s.%s\n", dbKey, collKey)

	if dbKey != "" {
		url, err := url.Parse(host.URI)
		if err != nil {
			runtime.LogWarningf(a.ctx, "Shell: failed to parse host URI %s: %s", host.URI, err.Error())
			result.ErrorTitle = "Could parse host URI"
			result.ErrorDescription = err.Error()
			return
		}

		url.Path = "/" + dbKey
		scriptHeader = scriptHeader + fmt.Sprintf("db = connect('%s');\n", url.String())
	}

	if collKey != "" {
		scriptHeader = scriptHeader + fmt.Sprintf("coll = db.getCollection('%s');\n", collKey)
	}

	scriptHeader = scriptHeader + "\n// Start of user script\n"
	script = scriptHeader + script

	if err := os.WriteFile(fname, []byte(script), os.ModePerm); err != nil {
		runtime.LogWarningf(a.ctx, "Shell: failed to write script to %s", err.Error())
		result.ErrorTitle = "Could not create temporary script file"
		result.ErrorDescription = err.Error()
		return
	}

	var outbuf, errbuf strings.Builder
	cmd := exec.Command("mongosh", "--file", fname, host.URI)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()

	if exiterr, ok := err.(*exec.ExitError); ok {
		result.Status = exiterr.ExitCode()
	} else if err != nil {
		runtime.LogWarningf(a.ctx, "Shell: failed to execute: mongosh --file %s: %s", fname, err.Error())
		result.ErrorTitle = "mongosh failure"
		result.ErrorDescription = err.Error()
		return
	} else {
		result.Status = 0
	}

	os.Remove(fname)
	result.Output = outbuf.String()
	result.Stderr = errbuf.String()
	return
}
