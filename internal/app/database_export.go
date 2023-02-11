package app

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ExportType string
type FileType string
type ExportInfo struct {
	Type     ExportType `json:"type"`
	FileType FileType   `json:"fileType"`
	OutDir   string     `json:"outdir"`
	Filename string     `json:"filename"`
	HostKey  string     `json:"hostKey"`
	DbKey    string     `json:"dbKey"`
	CollKeys []string   `json:"collKeys"`
}

const (
	ExportTypeExport ExportType = "export"
	ExportTypeDump   ExportType = "dump"

	FileTypeJson FileType = "json"
	FileTypeBson FileType = "bson"
	FileTypeCsv  FileType = "csv"
)

func (a *App) PerformExport(jsonData string) bool {
	var info ExportInfo
	err := json.Unmarshal([]byte(jsonData), &info)
	if err != nil {
		runtime.LogError(a.ctx, "Could not unmarshal export form")
		runtime.LogError(a.ctx, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Could not unmarshal JSON",
			Message: err.Error(),
		})
		return false
	}

	hosts, err := a.Hosts()
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Could not retrieve hosts",
			Message: err.Error(),
		})
		return false
	}
	host := hosts[info.HostKey]

	switch info.Type {
	case ExportTypeExport:
		if !a.Env.HasMongoExport {
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Type:  runtime.ErrorDialog,
				Title: "You need to install mongoexport to perform an export.",
			})
			return false
		}

		args := make([]string, 0)
		args = append(args, fmt.Sprintf(`--uri="%v"`, host.URI))
		args = append(args, fmt.Sprintf(`--type="%v"`, info.FileType))

		if info.DbKey != "" {
			args = append(args, fmt.Sprintf(`--db="%v"`, info.DbKey))

			if info.CollKeys != nil {
				args = append(args, fmt.Sprintf(`--collection="%v"`, info.CollKeys[0]))
			}
		}

		args = append(args, fmt.Sprintf(`--out="%v.%v"`, path.Join(info.OutDir, info.Filename), info.FileType))
		cmd := exec.Command("mongoexport", args...)
		var stdout strings.Builder
		var stderr strings.Builder
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()

		runtime.LogInfo(a.ctx, "Performing export with args: "+strings.Join(args, " "))

		fmt.Println(args)
		fmt.Println(stdout.String())
		fmt.Println(stderr.String())
		fmt.Println(err)

	case ExportTypeDump:
		if !a.Env.HasMongoDump {
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Type:  runtime.ErrorDialog,
				Title: "You need to install mongodump to perform a dump.",
			})
			return false
		}

	default:
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Unrecognised export type",
			Message: string(info.Type),
		})
		return false
	}

	return err == nil
}
