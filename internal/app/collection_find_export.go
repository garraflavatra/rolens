package app

import (
	"archive/zip"
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

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
	ExportFormatExcel     ExportFormat = "excel"
)

var (
	//go:embed collection_find_export_excel/app.xml
	excelAppXml string
	//go:embed collection_find_export_excel/core.xml
	excelCoreXml string
	//go:embed collection_find_export_excel/rels.xml
	excelRelsXml string
	//go:embed collection_find_export_excel/styles.xml
	excelStylesXml string
	//go:embed collection_find_export_excel/theme.xml
	excelThemeXml string
	//go:embed collection_find_export_excel/contenttypes.xml
	excelContentTypesXml string

	alphabet = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
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

func excelColIndex(idx int) string {
	str := make([]rune, 0)

	for idx > 0 {
		rem := idx % 26
		if rem == 0 {
			str = append(str, 'Z')
			idx = (idx / 26) - 1
		} else {

			str = append(str, alphabet[rem-1])
			idx = (idx / 26)
		}
	}

	// Reverse string
	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}

	return string(str)
}

func (a *App) PerformFindExport(hostKey, dbKey, collKey, settingsJson string) bool {
	runtime.LogInfof(a.ctx, "Export started for %s/%s/%s. Settings: %s", hostKey, dbKey, collKey, settingsJson)

	var settings ExportSettings
	if err := json.Unmarshal([]byte(settingsJson), &settings); err != nil {
		runtime.LogWarningf(a.ctx, "Could not parse export settings: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Couldn't parse export settings!",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
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
		runtime.LogDebugf(a.ctx, "Export: unknown view %s", settings.ViewKey)
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Message: fmt.Sprintf("View %s is not known", settings.ViewKey),
			Type:    runtime.ErrorDialog,
		})
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

	case ExportFormatExcel:
		defaultFilename = "export.xlsx"
		fileFilter = runtime.FileFilter{
			DisplayName: "Microsoft Excel Workbook (*.xlsx)",
			Pattern:     "*.xlsx",
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
		runtime.LogWarningf(a.ctx, "Export: error while choosing export destination: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error while choosing export destination",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}
	if settings.OutFile == "" {
		runtime.LogDebug(a.ctx, "Export: no destination specified")
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Message: "Please specify an export destination.",
			Type:    runtime.ErrorDialog,
		})
		return false
	}
	if _, err := os.Stat(settings.OutFile); err == nil {
		runtime.LogDebugf(a.ctx, "Export: destination %s already exists", settings.OutFile)
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Message: fmt.Sprintf("File %s already exists, export aborted.", settings.OutFile),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	var query bson.M
	if settings.Contents != ExportContentsAll {
		if err = bson.UnmarshalExtJSON([]byte(settings.QueryJson), true, &query); err != nil {
			runtime.LogDebugf(a.ctx, "Invalid find query (exporting): %s", settings.QueryJson)
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Title:   "Invalid query",
				Message: err.Error(),
				Type:    runtime.ErrorDialog,
			})
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
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Couldn't get cursor",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	file, err := os.OpenFile(settings.OutFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		runtime.LogDebugf(a.ctx, "Export: unable to open file %s", settings.OutFile)
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error opening file",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}
	defer file.Close()

	var index uint = 0
	var columnKeys []string

	var csvWriter *csv.Writer

	var excelZipWriter *zip.Writer
	var excelSheetWriter io.Writer
	var excelStrings = make([]string, 0)

	switch settings.Format {
	case ExportFormatJsonArray:
		file.WriteString("[")

	case ExportFormatCsv:
		csvWriter = csv.NewWriter(file)

	case ExportFormatExcel:
		excelZipWriter = zip.NewWriter(file)

		files := map[string]string{
			"docProps/app.xml":          excelAppXml,
			"docProps/core.xml":         strings.Replace(excelCoreXml, "{TITLE}", fmt.Sprintf("%s.%s", dbKey, collKey), 1),
			"xl/theme/theme1.xml":       excelThemeXml,
			"xl/rels/workbook.xml.rels": excelRelsXml,
			"xl/styles.xml":             excelStylesXml,
			"[Content-Types].xml":       excelContentTypesXml,
		}

		for fname, body := range files {
			f, err := excelZipWriter.Create(fname)
			if err != nil {
				runtime.LogErrorf(a.ctx, "Export: Excel zip.Create error: %s", err.Error())
				runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
					Title:   "ZIP error!",
					Message: err.Error(),
					Type:    runtime.ErrorDialog,
				})
				return false
			}

			_, err = f.Write([]byte(body))
			if err != nil {
				runtime.LogErrorf(a.ctx, "Export: Excel zip.Create.Write error: %s", err.Error())
				runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
					Title:   "ZIP error!",
					Message: err.Error(),
					Type:    runtime.ErrorDialog,
				})
				return false
			}
		}

		excelZipWriter.Create("_rels/")

		excelSheetWriter, err = excelZipWriter.Create("xl/worksheets/sheet1.xml")
		if err != nil {
			runtime.LogErrorf(a.ctx, "Export: Excel ZIP error creating worksheet: %s", err.Error())
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Title:   "ZIP error!",
				Message: err.Error(),
				Type:    runtime.ErrorDialog,
			})
			return false
		}

		excelSheetWriter.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<worksheet
    xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"
    xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
    xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006" mc:Ignorable="x14ac xr xr2 xr3"
    xmlns:x14ac="http://schemas.microsoft.com/office/spreadsheetml/2009/9/ac"
    xmlns:xr="http://schemas.microsoft.com/office/spreadsheetml/2014/revision"
    xmlns:xr2="http://schemas.microsoft.com/office/spreadsheetml/2015/revision2"
    xmlns:xr3="http://schemas.microsoft.com/office/spreadsheetml/2016/revision3" xr:uid="{17671867-E8D5-A14C-B382-03A6AA54A004}">
    <dimension ref="A1:C5"/>
    <sheetViews>
        <sheetView tabSelected="1" workbookViewId="0">
            <selection activeCell="A1" sqref="A1"/>
        </sheetView>
    </sheetViews>
    <sheetFormatPr baseColWidth="10" defaultRowHeight="16" x14ac:dyDescent="0.2"/>
    <sheetData>`))
	}

	for cur.Next(ctx) {
		if settings.ViewKey == "list" && columnKeys == nil {
			columnKeys = make([]string, 0)
			els, err := cur.Current.Elements()
			if err != nil {
				runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
					Title:   "BSON is invalid",
					Message: err.Error(),
					Type:    runtime.ErrorDialog,
				})
			}

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
					columnKeys = append(columnKeys, el.Key())
				}
			}

			runtime.LogDebugf(a.ctx, "Export column keys: %v", columnKeys)

			switch settings.Format {
			case ExportFormatCsv:
				if err := csvWriter.Write(columnKeys); err != nil {
					runtime.LogInfof(a.ctx, "Unable to write item %d to CSV while exporting: %s", index, err.Error())
					runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
						Title:   fmt.Sprintf("Unable to write item %d to CSV", index),
						Message: err.Error(),
						Type:    runtime.ErrorDialog,
					})
				}

			case ExportFormatExcel:
				excelSheetWriter.Write([]byte(fmt.Sprintf(`<row r="1" spans="1:%d" s="1" customFormat="1" x14ac:dyDescent="0.2">`, len(columnKeys))))
				for idx, key := range columnKeys {
					// excelStringsWriter.Write([]byte(fmt.Sprintf("<si><t>%s</t></si>", key)))
					excelStrings = append(excelStrings, key)
					excelSheetWriter.Write([]byte(fmt.Sprintf(`<c r="%s1" s="1" t="s"><v>%d</v></c>`, excelColIndex(idx+1), len(excelStrings))))
				}
				excelSheetWriter.Write([]byte("</row>"))
			}
		}

		switch settings.Format {
		case ExportFormatCsv:
			csvItem := make([]string, 0)

			switch settings.ViewKey {
			case "list":
				for _, k := range columnKeys {
					r, err := cur.Current.LookupErr(k)
					if err != nil {
						csvItem = append(csvItem, "")
						continue
					}

					var v any
					if err := r.Unmarshal(&v); err != nil {
						runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
							Title:   fmt.Sprintf("Unable to unmarshal field %s", k),
							Message: err.Error(),
							Type:    runtime.ErrorDialog,
						})
						csvItem = append(csvItem, "")
						continue
					}

					csvItem = append(csvItem, fmt.Sprintf("%v", v))
				}

			default:
				// @todo
			}

			if err := csvWriter.Write(csvItem); err != nil {
				runtime.LogInfof(a.ctx, "Export: Unable to write item %d to CSV while exporting: %s", index, err.Error())
				runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
					Title:   fmt.Sprintf("Unable to write item %d to CSV", index),
					Message: err.Error(),
					Type:    runtime.ErrorDialog,
				})
			}

			csvWriter.Flush()

		case ExportFormatJsonArray, ExportFormatNdJson:
			itemJson, err := bson.MarshalExtJSON(cur.Current, true, false)
			if err != nil {
				runtime.LogInfof(a.ctx, "Export: Unable to marshal item %d to JSON while exporting: %s", index, err.Error())
				runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
					Title:   fmt.Sprintf("Unable to write item %d to CSV", index),
					Message: err.Error(),
					Type:    runtime.ErrorDialog,
				})
			}

			if (settings.Format == ExportFormatJsonArray) && (index != 0) {
				file.WriteString(",\n")
			}

			file.Write(itemJson)

			if settings.Format == ExportFormatNdJson {
				file.WriteString("\n")
			}

		case ExportFormatExcel:
			excelRow := make([]string, 0)

			for _, k := range columnKeys {
				r, err := cur.Current.LookupErr(k)
				if err != nil {
					excelRow = append(excelRow, "")
					continue
				}

				var v any
				if err := r.Unmarshal(&v); err != nil {
					runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
						Title:   fmt.Sprintf("Unable to unmarshal field %s", k),
						Message: err.Error(),
						Type:    runtime.ErrorDialog,
					})
					excelRow = append(excelRow, "")
					continue
				}

				excelRow = append(excelRow, fmt.Sprintf("%v", v))
			}

			excelSheetWriter.Write([]byte(fmt.Sprintf(`<row r="%d" spans="1:%d" s="1" x14ac:dyDescent="0.2">`, index+2, len(columnKeys))))
			for idx, val := range excelRow {
				// excelStringsWriter.Write([]byte(fmt.Sprintf("<si><t>%s</t></si>", val)))
				excelStrings = append(excelStrings, val)
				excelSheetWriter.Write([]byte(fmt.Sprintf(`<c r="%s%d" t="s"><v>%d</v></c>`, excelColIndex(idx+1), index+2, len(excelStrings))))
			}
			excelSheetWriter.Write([]byte("</row>"))
		}

		index++
	}

	switch settings.Format {
	case ExportFormatJsonArray:
		file.WriteString("]\n")

	case ExportFormatExcel:
		excelSheetWriter.Write([]byte(`</sheetData><pageMargins left="0.7" right="0.7" top="0.75" bottom="0.75" header="0.3" footer="0.3"/></worksheet>`))

		excelStringsWriter, err := excelZipWriter.Create("xl/sharedStrings.xml")
		if err != nil {
			runtime.LogErrorf(a.ctx, "Export: Excel ZIP error creating shared strings: %s", err.Error())
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Title:   "ZIP error!",
				Message: err.Error(),
				Type:    runtime.ErrorDialog,
			})
			return false
		}

		excelStringsWriter.Write([]byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" count="%d" uniqueCount="%d">%s`, len(excelStrings), len(excelStrings), "\r\n")))
		for _, str := range excelStrings {
			excelStringsWriter.Write([]byte(fmt.Sprintf("<si><t>%s</t></si>\r\n", str)))
		}
		excelStringsWriter.Write([]byte("</sst>"))

		if err := excelZipWriter.Close(); err != nil {
			runtime.LogErrorf(a.ctx, "Export: Excel ZIP error while closing: %s", err.Error())
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Title:   "ZIP error!",
				Message: err.Error(),
				Type:    runtime.ErrorDialog,
			})
		}
	}

	a.ui.Reveal(settings.OutFile)
	runtime.LogInfof(a.ctx, "Export succeeded: %d items", count)
	return true
}
