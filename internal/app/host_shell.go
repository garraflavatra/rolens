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

type SaveShellScriptResult struct {
	Host             Host   `json:"host"`
	Fname            string `json:"filename"`
	ErrorTitle       string `json:"errorTitle"`
	ErrorDescription string `json:"errorDescription"`
}

func (a *App) ExecuteShellScript(hostKey, dbKey, collKey, script string) (result ExecuteShellScriptResult) {
	if !a.Env.HasMongoShell {
		result.ErrorTitle = "mongosh not found"
		result.ErrorDescription = "The mongosh executable is required to run a shell script. Please see https://www.mongodb.com/docs/mongodb-shell/install/"
		return
	}

	saveRes := a.SaveShellScript(hostKey, dbKey, collKey, script, true)
	if (saveRes.ErrorTitle != "") || (saveRes.ErrorDescription != "") {
		result.ErrorTitle = saveRes.ErrorTitle
		result.ErrorDescription = saveRes.ErrorDescription
		return
	}

	var outbuf, errbuf strings.Builder
	cmd := exec.Command("mongosh", "--file", saveRes.Fname, saveRes.Host.URI)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()

	if exiterr, ok := err.(*exec.ExitError); ok {
		result.Status = exiterr.ExitCode()
	} else if err != nil {
		runtime.LogWarningf(a.ctx, "Shell: failed to execute: mongosh --file %s: %s", saveRes.Fname, err.Error())
		result.ErrorTitle = "mongosh failure"
		result.ErrorDescription = err.Error()
		return
	} else {
		result.Status = 0
	}

	os.Remove(saveRes.Fname)
	result.Output = outbuf.String()
	result.Stderr = errbuf.String()
	return
}

func (a *App) SaveShellScript(hostKey, dbKey, collKey, script string, temp bool) (result SaveShellScriptResult) {
	hosts, err := a.Hosts()
	if err != nil {
		runtime.LogWarningf(a.ctx, "Shell: could not get hosts: %s", err.Error())
		result.ErrorTitle = "Could not get hosts"
		result.ErrorDescription = err.Error()
		return
	}

	host, hostFound := hosts[hostKey]
	if !hostFound {
		runtime.LogWarningf(a.ctx, "Shell: host %s does not exist", host)
		result.ErrorTitle = "The specified host does not seem to exist"
		return
	}

	result.Host = host
	id, err := uuid.NewRandom()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Shell: failed to generate a UUID: %s", err.Error())
		result.ErrorTitle = "Could not generate UUID"
		result.ErrorDescription = err.Error()
		return
	}

	if temp {
		dirname, err := os.MkdirTemp(os.TempDir(), "rolens-script")
		if err != nil {
			runtime.LogErrorf(a.ctx, "Shell: failed to create temporary directory: %s", err.Error())
			result.ErrorTitle = "Could not generate temporary directory for script"
			result.ErrorDescription = err.Error()
			return
		}
		result.Fname = path.Join(dirname, fmt.Sprintf("%s.mongosh.js", id.String()))
	} else {
		result.Fname, err = runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
			DefaultFilename:      "New Script.js",
			DefaultDirectory:     path.Join(a.Env.DataDirectory, "Shell Scripts"),
			Title:                "Save MongoDB Shell Script",
			CanCreateDirectories: true,
			Filters: []runtime.FileFilter{
				{
					DisplayName: "MongoDB Shell Script (*.js)",
					Pattern:     "*.js",
				},
			},
		})

		if err != nil {
			runtime.LogErrorf(a.ctx, "Shell: failed to save script: %s", err.Error())
			result.ErrorTitle = "Could not save shell script"
			result.ErrorDescription = err.Error()
			return
		}
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

	scriptHeader = scriptHeader + "\n"
	script = scriptHeader + strings.TrimLeft(strings.TrimRight(script, " \t\n"), "\n")

	if err := os.WriteFile(result.Fname, []byte(script), 0755); err != nil {
		runtime.LogWarningf(a.ctx, "Shell: failed to write script to %s: %s", result.Fname, err.Error())
		result.ErrorTitle = "Could not create temporary script file"
		result.ErrorDescription = err.Error()
		return
	}

	return
}

func (a *App) OpenShellScript() string {
	dir := path.Join(a.Env.DataDirectory, "Shell Scripts")
	os.MkdirAll(dir, 0755)

	fname, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory:     path.Join(a.Env.DataDirectory, "Shell Scripts"),
		Title:                "Load a MongoDB Shell Script",
		CanCreateDirectories: true,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "MongoDB Shell Script (*.js)",
				Pattern:     "*.js",
			},
		},
	})

	if err != nil {
		runtime.LogWarningf(a.ctx, "Shell: error opening script: %s", err.Error())
		return ""
	}

	script, err := os.ReadFile(fname)

	if err != nil {
		runtime.LogWarningf(a.ctx, "Shell: error reading script %s: %s", fname, err.Error())
		return ""
	}

	return string(script)
}

func (a *App) SaveShellOuput(output string) {
	fname, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename:      "mongosh-output.txt",
		DefaultDirectory:     a.Env.DownloadDirectory,
		Title:                "Save mongosh output",
		CanCreateDirectories: true,
	})

	if err != nil {
		runtime.LogWarningf(a.ctx, "Shell: error exporting output to %s: %s", fname, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title: "Error exporting output",
			Message: err.Error(),
			Type: runtime.ErrorDialog,
		})
	}

	if err := os.WriteFile(fname, []byte(output), 0755); err != nil {
		runtime.LogWarningf(a.ctx, "Shell: error writing shell output to %s: %s", fname, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title: "Error writing shell output",
			Message: err.Error(),
			Type: runtime.ErrorDialog,
		})
	}
}
