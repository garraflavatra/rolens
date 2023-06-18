package app

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ncruces/zenity"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExportContents string
type ExportFormat string

const (
	ExportContentsAll            ExportContents = "all"
	ExportContentsQuery          ExportContents = "query"
	ExportContentsQueryLimitSkip ExportContents = "querylimitskip"

	ExportFormatJsonArray ExportFormat = "jsonarray"
	ExportFormatNdJson    ExportFormat = "ndjson"
	ExportFormatCsv       ExportFormat = "csv"
)

type ExportSettings struct {
	Contents  ExportContents `json:"contents"`
	Format    ExportFormat   `json:"format"`
	ViewKey   string         `json:"viewKey"`
	QueryJson string         `json:"query"`
	Limit     uint           `json:"limit"`
	Skip      uint           `json:"skip"`
	OutFile   string         `json:"outfile"`
}

func getptr[T any](v T) *T {
	return &v
}

func (a *App) PerformFindExport(hostKey, dbKey, collKey, settingsJson string) bool {
	runtime.LogInfof(a.ctx, "Export started for %s/%s/%s. Settings: %s", hostKey, dbKey, collKey, settingsJson)

	var settings ExportSettings
	if err := json.Unmarshal([]byte(settingsJson), &settings); err != nil {
		runtime.LogWarningf(a.ctx, "Could not parse export settings: %s", err.Error())
		zenity.Error(err.Error(), zenity.Title("Couldn't parse export settings!"), zenity.ErrorIcon)
		return false
	}

	switch settings.Contents {
	case ExportContentsAll:
		settings.QueryJson = "{}"
		settings.Limit = 0
		settings.Skip = 0
	case ExportContentsQuery:
		settings.Limit = 0
		settings.Skip = 0
	case ExportContentsQueryLimitSkip:
	}

	views, err := a.Views()
	if err != nil {
		runtime.LogWarningf(a.ctx, "Export: error while retrieving view: %s", err.Error())
		return false
	}

	view, found := views[settings.ViewKey]
	if !found {
		zenity.Error(fmt.Sprintf("View %s is not known", settings.ViewKey), zenity.ErrorIcon)
		runtime.LogDebugf(a.ctx, "Export: unknown view %s", settings.ViewKey)
		return false
	}

	var fileFilter runtime.FileFilter
	defaultFilename := ""

	switch settings.Format {
	case ExportFormatCsv:
		defaultFilename = "export.csv"
		fileFilter = runtime.FileFilter{
			DisplayName: "CSV files (*.csv)",
			Pattern:     "*.csv",
		}
	case ExportFormatJsonArray:
		defaultFilename = "export.json"
		fileFilter = runtime.FileFilter{
			DisplayName: "JSON files (*.json)",
			Pattern:     "*.json",
		}
	case ExportFormatNdJson:
		defaultFilename = "export.ndjson"
		fileFilter = runtime.FileFilter{
			DisplayName: "Newline delimited JSON files (*.ndjson)",
			Pattern:     "*.ndjson",
		}
	}

	settings.OutFile, err = runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:                "Choose export destination",
		DefaultDirectory:     a.Env.DownloadDirectory,
		CanCreateDirectories: true,
		DefaultFilename:      defaultFilename,
		Filters:              []runtime.FileFilter{fileFilter},
	})
	if err != nil {
		zenity.Error("An error occured while choosing the export destination", zenity.ErrorIcon)
		runtime.LogWarningf(a.ctx, "Export: error while choosing export destination: %s", err.Error())
		return false
	}
	if settings.OutFile == "" {
		zenity.Error("You must specify an export destination.", zenity.ErrorIcon)
		runtime.LogDebug(a.ctx, "Export: no destination specified")
		return false
	}

	if _, err := os.Stat(settings.OutFile); err == nil {
		zenity.Error(fmt.Sprintf("File %s already exists, export aborted.", settings.OutFile), zenity.ErrorIcon)
		runtime.LogDebugf(a.ctx, "Export: destination %s already exists", settings.OutFile)
		return false
	}

	var query bson.M
	if settings.Contents != ExportContentsAll {
		if err = bson.UnmarshalExtJSON([]byte(settings.QueryJson), true, &query); err != nil {
			runtime.LogDebugf(a.ctx, "Invalid find query (exporting): %s", settings.QueryJson)
			zenity.Error(err.Error(), zenity.Title("Invalid query"), zenity.ErrorIcon)
			return false
		}
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return false
	}
	defer close()

	pgr, _ := zenity.Progress(zenity.Title("Performing export…"))

	projection := bson.M{}
	if settings.ViewKey != "list" {
		for _, col := range view.Columns {
			projection[col.Key] = ""
		}
	}

	var count uint = 0
	if settings.Limit != 0 {
		count = uint(settings.Limit)
	} else {
		c, _ := client.Database(dbKey).Collection(collKey).CountDocuments(ctx, query, &options.CountOptions{
			Skip: getptr(int64(settings.Skip)),
		})
		count = uint(c)
	}

	cur, err := client.Database(dbKey).Collection(collKey).Find(ctx, query, &options.FindOptions{
		Skip:       getptr(int64(settings.Skip)),
		Limit:      getptr(int64(settings.Limit)),
		Projection: projection,
	})
	if err != nil {
		runtime.LogInfof(a.ctx, "Export: unable to get cursor while exporting: %s", err.Error())
		zenity.Error(err.Error(), zenity.Title("Unable to get cursor"), zenity.ErrorIcon)
		return false
	}

	file, err := os.OpenFile(settings.OutFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		zenity.Error(fmt.Sprintf(err.Error(), zenity.Title("Error while opening file"), settings.OutFile), zenity.ErrorIcon)
		runtime.LogDebugf(a.ctx, "Export: unable to open file %s", settings.OutFile)
		return false
	}
	defer file.Close()

	var csvWriter *csv.Writer
	var csvColumnKeys []string
	var index uint = 0

	switch settings.Format {
	case ExportFormatJsonArray:
		file.WriteString("[")

	case ExportFormatCsv:
		csvWriter = csv.NewWriter(file)
	}

	for cur.Next(ctx) {
		switch settings.Format {
		case ExportFormatCsv:
			els, err := cur.Current.Elements()
			if err != nil {
				zenity.Error(err.Error(), zenity.Title("BSON invalid"), zenity.ErrorIcon)
			}

			csvItem := make([]string, 0)

			switch settings.ViewKey {
			case "list":
				if csvColumnKeys == nil {
					csvColumnKeys = make([]string, 0)

					for _, el := range els {
						if el.Key() == "" {
							continue
						}

						switch el.Value().Type {
						case bsontype.Boolean,
							bsontype.Decimal128,
							bsontype.Double,
							bsontype.Int32,
							bsontype.Int64,
							bsontype.Null,
							bsontype.ObjectID,
							bsontype.Regex,
							bsontype.String,
							bsontype.Symbol,
							bsontype.Timestamp,
							bsontype.Undefined:
							csvColumnKeys = append(csvColumnKeys, el.Key())
						}
					}

					runtime.LogDebugf(a.ctx, "Export csvColumnKeys: %v", csvColumnKeys)

					if err := csvWriter.Write(csvColumnKeys); err != nil {
						runtime.LogInfof(a.ctx, "Unable to write item %d to CSV while exporting: %s", index, err.Error())
						zenity.Error(err.Error(), zenity.Title("Unable to write item %d to CSV"), zenity.ErrorIcon)
					}
				}

				for _, k := range csvColumnKeys {
					r, err := cur.Current.LookupErr(k)
					if err != nil {
						csvItem = append(csvItem, "")
						continue
					}

					var v any
					if err := r.Unmarshal(&v); err != nil {
						zenity.Error(err.Error(), zenity.Title(fmt.Sprintf("Unable to unmarshal field %s", k)), zenity.ErrorIcon)
						csvItem = append(csvItem, "")
						continue
					}

					csvItem = append(csvItem, fmt.Sprintf("%v", v))
				}

			default:
				// @todo
			}

			if err := csvWriter.Write(csvItem); err != nil {
				runtime.LogInfof(a.ctx, "Unable to write item %d to CSV while exporting: %s", index, err.Error())
				zenity.Error(err.Error(), zenity.Title("Unable to write item %d to CSV"), zenity.ErrorIcon)
			}

			csvWriter.Flush()

		case ExportFormatJsonArray, ExportFormatNdJson:
			itemJson, err := bson.MarshalExtJSON(cur.Current, true, false)
			if err != nil {
				runtime.LogInfof(a.ctx, "Unable to marshal item %d to JSON while exporting: %s", index, err.Error())
				zenity.Error(err.Error(), zenity.Title("Unable to marshal item %d to JSON"), zenity.ErrorIcon)
			}

			if (settings.Format == ExportFormatJsonArray) && (index != 0) {
				file.WriteString(",\n")
			}

			file.Write(itemJson)

			if settings.Format == ExportFormatNdJson {
				file.WriteString("\n")
			}
		}

		index++

		if count != 0 && pgr != nil {
			p := (float32(index) / float32(count)) * 100.0
			pgr.Value(int(p))
		}
	}

	if settings.Format == ExportFormatJsonArray {
		file.WriteString("]\n")
	}

	if pgr != nil {
		pgr.Complete()
		pgr.Close()
	}

	a.ui.Reveal(settings.OutFile)
	runtime.LogInfo(a.ctx, "Export succeeded")
	return true
}
