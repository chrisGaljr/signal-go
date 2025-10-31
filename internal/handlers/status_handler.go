package handlers

import (
	"fmt"
	"net/http"
	"signal/main/internal/models"
	"signal/main/internal/services"
	"signal/main/internal/utils"
	"text/template"
	"time"
)

type StatusPageData struct {
	Status        string
	IsOk          bool
	Error         error
	LastChecked   string
	Bars          []utils.Bar
	Direction     string
	SignalSuccess int
	Uptime        string
}

func StatusPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("template/status/index.html")

	status, err := services.CheckChromedpHealth()
	messageSuccessRate, direction, uptime := getStatistics()

	tmpl.Execute(w, StatusPageData{
		Status:        status,
		IsOk:          err == nil,
		Error:         err,
		LastChecked:   time.Now().Format("2006-01-02 15:04"),
		Bars:          utils.GetStatusBarData(messageSuccessRate),
		Direction:     direction,
		SignalSuccess: messageSuccessRate,
		Uptime:        uptime,
	})
}

func getStatistics() (int, string, string) {
	config, err := models.GetConfig()
	if err != nil {
		return 0, "unknown", "unknown"
	}

	messageSuccessRate := float64(100)
	if (config.MessageSent + config.FailedToSend) > 0 {
		messageSuccessRate = (float64(config.MessageSent) / float64(config.MessageSent+config.FailedToSend)) * 100
	}

	startTime, _ := time.Parse(time.RFC3339, config.StartTime)
	diff := time.Since(startTime)
	oneDayInHours, oneYearInHours := 24, 8760
	years := int(diff.Hours()) / oneYearInHours
	days := int(diff.Hours()) / oneDayInHours
	hours := int(diff.Hours()) % oneDayInHours

	if days == 0 && hours == 0 {
		return int(messageSuccessRate), config.Direction, "less than an hour"
	}

	var uptime string
	if days == 0 {
		uptime = fmt.Sprintf("%d hour(s)", hours)
	} else if years > 0 {
		uptime = fmt.Sprintf("%d year(s)", years)
	} else {
		uptime = fmt.Sprintf("%d day(s), %d hour(s)", days, hours)
	}

	return int(messageSuccessRate), config.Direction, uptime
}
