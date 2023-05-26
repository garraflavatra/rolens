package app

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ncruces/zenity"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExportContents string
type ExportFormat string

const (
	ExportContentsAll            ExportContents = "all"
	ExportContentsQuery          ExportContents = "query"
	ExportContentsQueryLimitSkip ExportContents = "querylimitskip"

	ExportFormatJsonArray   ExportFormat = "jsonarray"
	ExportFormatJsonNewline ExportFormat = "jsonnewline"
	ExportFormatCsv         ExportFormat = "csv"
)

type ExportSettings struct {
	Contents  ExportContents `json:"contents"`
	Format    ExportFormat   `json:"format"`
	ViewKey   string         `json:"viewKey"`
	QueryJson string         `json:"query"`
	Limit     int64          `json:"limit"`
	Skip      int64          `json:"skip"`
	OutFile   string         `json:"outfile"`
}

func (a *App) PerformFindExport(hostKey, dbKey, collKey, settingsJson string) bool {
	var settings ExportSettings
	if err := json.Unmarshal([]byte(settingsJson), &settings); err != nil {
		runtime.LogWarning(a.ctx, "Could not parse export settings:")
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Error(err.Error(), zenity.Title("Couldn't parse export settings!"), zenity.ErrorIcon)
		return false
	}

	if _, err := os.Stat(settings.OutFile); err == nil {
		zenity.Error(fmt.Sprintf("File %s already exists, export aborted.", settings.OutFile), zenity.ErrorIcon)
		return false
	}

	views, err := a.Views()
	if err != nil {
		return false
	}

	view, found := views[settings.ViewKey]
	if !found {
		zenity.Error(fmt.Sprintf("View %s is not known", settings.ViewKey), zenity.ErrorIcon)
		return false
	}

	var query bson.M
	if settings.Contents != ExportContentsAll {
		if err = bson.UnmarshalExtJSON([]byte(settings.QueryJson), true, &query); err != nil {
			runtime.LogInfo(a.ctx, "Invalid find query (exporting):")
			runtime.LogInfo(a.ctx, err.Error())
			zenity.Error(err.Error(), zenity.Title("Invalid query"), zenity.ErrorIcon)
			return false
		}
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return false
	}
	defer close()

	projection := bson.M{}

	if settings.ViewKey != "list" {
		for _, col := range view.Columns {
			projection[col.Key] = ""
		}
	}

	cur, err := client.Database(dbKey).Collection(collKey).Find(ctx, query, &options.FindOptions{
		Skip:       &settings.Skip,
		Limit:      &settings.Limit,
		Projection: projection,
	})
	if err != nil {
		runtime.LogInfo(a.ctx, "Unable to get cursor while exporting:")
		runtime.LogInfo(a.ctx, err.Error())
		zenity.Error(err.Error(), zenity.Title("Unable to get cursor"), zenity.ErrorIcon)
		return false
	}

	file, err := os.OpenFile(settings.OutFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		zenity.Error(fmt.Sprintf(err.Error(), zenity.Title("Error while opening file"), settings.OutFile), zenity.ErrorIcon)
		return false
	}
	defer file.Close()

	var csvWriter *csv.Writer
	var csvColumnKeys []any

	switch settings.Format {
	case ExportFormatJsonArray:
		file.WriteString("[\n")
	case ExportFormatCsv:
		csvWriter = csv.NewWriter(file)
	}

	index := -1
	for cur.Next(ctx) {
		index++
		var item map[any]interface{}

		if err := bson.Unmarshal(cur.Current, &item); err != nil {
			runtime.LogInfo(a.ctx, fmt.Sprintf("Unable to unmarshal item %d while exporting", index))
			runtime.LogInfo(a.ctx, err.Error())
			zenity.Error(err.Error(), zenity.Title("Unable to unmarshal item %d"), zenity.ErrorIcon)
			continue
		}

		switch settings.Format {
		case ExportFormatCsv:
			csvItem := make([]string, 0)

			switch settings.ViewKey {
			case "list":
				if csvColumnKeys == nil {
					csvColumnKeys = make([]any, 0)
					for k := range item {
						csvColumnKeys = append(csvColumnKeys, k)
					}
				}

				for _, k := range csvColumnKeys {
					csvItem = append(csvItem, item[k].(string))
				}

			default:
				for _, v := range item {
					csvItem = append(csvItem, v.(string))
				}
			}

			if err := csvWriter.Write(csvItem); err != nil {
				runtime.LogInfo(a.ctx, fmt.Sprintf("Unable to write item %d to CSV while exporting", index))
				runtime.LogInfo(a.ctx, err.Error())
				zenity.Error(err.Error(), zenity.Title("Unable to write item %d to CSV"), zenity.ErrorIcon)
			}

		case ExportFormatJsonArray, ExportFormatJsonNewline:
			itemJson, err := json.Marshal(item)
			if err != nil {
				runtime.LogInfo(a.ctx, fmt.Sprintf("Unable to marshal item %d to JSON while exporting", index))
				runtime.LogInfo(a.ctx, err.Error())
				zenity.Error(err.Error(), zenity.Title("Unable to marshal item %d to JSON"), zenity.ErrorIcon)
			}

			file.Write(itemJson)
			if settings.Format == ExportFormatJsonArray {
				file.WriteString(",")
			}
			file.WriteString("\n")
		}
	}

	if settings.Format == ExportFormatJsonArray {
		file.WriteString("]\n")
	}

	return true
}
