package app

import (
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SavedQuery struct {
	Query
	Name    string `json:"name"`
	Remarks string `json:"remarks"`
	HostKey string `json:"hostKey"`
	DbKey   string `json:"dbKey"`
	CollKey string `json:"collKey"`
}

func updateQueryFile(a *App, newData map[string]SavedQuery) error {
	filePath := path.Join(a.Env.DataDirectory, "queries.json")
	jsonData, err := json.MarshalIndent(newData, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, jsonData, 0644)
	return err
}

func (a *App) SavedQueries() map[string]SavedQuery {
	filePath := path.Join(a.Env.DataDirectory, "queries.json")
	jsonData, err := ioutil.ReadFile(filePath)

	if err != nil {
		// It's ok if the file cannot be opened, for example if it is not accessible.
		// Therefore no error is returned.
		runtime.LogInfof(a.ctx, "Could not open queries.json: %s", err.Error())
		return make(map[string]SavedQuery, 0)
	}

	if len(jsonData) == 0 {
		return make(map[string]SavedQuery, 0)
	} else {
		var queries map[string]SavedQuery
		err = json.Unmarshal(jsonData, &queries)

		if err != nil {
			runtime.LogInfof(a.ctx, "queries.json file contains malformatted JSON data: %s", err.Error())
			return nil
		}

		return queries
	}
}

func (a *App) SaveQuery(jsonData string) string {
	var query SavedQuery
	err := json.Unmarshal([]byte(jsonData), &query)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Add query: malformed form: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Malformed JSON",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return ""
	}

	queries := a.SavedQueries()
	queries[query.Name] = query
	err = updateQueryFile(a, queries)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Could not update query list",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return ""
	}

	return query.Name
}

func (a *App) RemoveQuery(queryName string) {
	queries := a.SavedQueries()
	delete(queries, queryName)
	if err := updateQueryFile(a, queries); err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Could not update query list",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
	}
}

func (a *App) UpdateQueries(jsonData string) bool {
	var queries map[string]SavedQuery
	err := json.Unmarshal([]byte(jsonData), &queries)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Update queries: malformed form: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Malformed JSON",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	err = updateQueryFile(a, queries)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Could not save queries",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	return true
}
