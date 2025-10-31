package handlers

import (
	"net/http"
	"signal/main/internal/models"
	"text/template"
)

type IncidentHistoryPageData struct {
	HasError     bool
	NumberOfLogs int
	Logs         []models.ErrorLog
}

func IncidentHistoryPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("template/incidentHistory/index.html")
	results, err := models.GetRecentErrorLogs(5)

	tmpl.Execute(w, IncidentHistoryPageData{
		NumberOfLogs: len(results),
		Logs:         results,
		HasError:     err != nil,
	})
}
