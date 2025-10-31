package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"runtime/debug"
	"signal/main/internal/models"
	"signal/main/internal/utils"
)

func CheckSignalIsUp() bool {
	url := os.Getenv("SIGNAL_REST_BASE_URL") + "v1/about"
	response, err := http.Get(url)
	if err != nil {
		models.SaveErrorLog([]byte(url), debug.Stack(), "Error getting Signal API status", err, 0)
		return false
	}
	defer response.Body.Close()

	return response.StatusCode == 200
}

func SendSignalMessage(msg string, recipients []string) {
	signalIsOk := CheckSignalIsUp()
	if !signalIsOk || msg == "" {
		return
	}

	url := os.Getenv("SIGNAL_REST_BASE_URL") + "v2/send"
	bod := PostBody{
		Message:    msg,
		Number:     os.Getenv("BACKEND_NUMBER"),
		Recipients: recipients,
	}
	jsonData, _ := json.Marshal(bod)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		models.SaveErrorLog(jsonData, debug.Stack(), utils.GetSignalError(), err, 0)
		return
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		models.SaveErrorLog(jsonData, debug.Stack(), utils.GetSignalError(), nil, response.StatusCode)
	} else {
		models.UpdateConfig("message_sent")
	}
}

type PostBody struct {
	Message    string   `json:"message"`
	Number     string   `json:"number"`
	Recipients []string `json:"recipients"`
}
