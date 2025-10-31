package models

import (
	"context"
	"encoding/json"
	"fmt"
	"signal/main/internal/utils"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ErrorLog struct {
	CreatedDate  string `bson:"created_date"`
	Content      string `bson:"content"`
	StackTrace   string `bson:"stack_trace"`
	ErrorMessage string `bson:"error_message"`
	Error        string `bson:"error"`
	StatusCode   int    `bson:"status_code,omitempty"`
}

func CreateErrorLog(data string, StackTrace string, msg string, err error, statusCode int) ErrorLog {
	log := ErrorLog{
		CreatedDate:  time.Now().Format(time.RFC3339),
		Content:      data,
		StackTrace:   StackTrace,
		ErrorMessage: msg,
		Error:        err.Error(),
		StatusCode:   statusCode,
	}

	return log
}

func SaveErrorLog(data []byte, StackTrace []byte, msg string, err error, statusCode int) error {
	log := CreateErrorLog(string(data), string(StackTrace), msg, err, statusCode)
	if msg == utils.GetSignalError() {
		UpdateConfig("failed_to_send")
	}

	return SaveErrorLogs([]ErrorLog{log})
}

func SaveErrorLogs(logs []ErrorLog) error {
	if collections == nil || collections["error_log"] == nil {
		return fmt.Errorf("database not connected. Please call ConnectDB() before querying")
	}

	_, err := collections["error_log"].InsertMany(context.TODO(), logs)
	if err != nil {
		return fmt.Errorf("failed to save error logs: %v", err)
	}
	return nil
}

func GetRecentErrorLogs(limit int) ([]ErrorLog, error) {
	return GetErrorLogs(limit, bson.D{}, bson.D{{Key: "created_date", Value: -1}})
}

func GetErrorLogs(numOfLogs int, filter bson.D, sorting bson.D) ([]ErrorLog, error) {
	var results []ErrorLog

	if collections == nil || collections["error_log"] == nil {
		return results, fmt.Errorf("database not connected. Please call ConnectDB() before querying")
	}

	cursor, err := collections["error_log"].Find(
		context.TODO(),
		filter,
		options.Find().SetLimit(int64(numOfLogs)),
		options.Find().SetSort(sorting),
	)

	if err != nil {
		return results, fmt.Errorf("unexpected error occurred while fetching error logs: %v", err)
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return results, fmt.Errorf("failed to retrieve error logs: %v", err)
	}

	var logs []ErrorLog
	for _, result := range results {
		createdDate, _ := time.Parse(time.RFC3339, result.CreatedDate)
		result.CreatedDate = createdDate.Format("2006-01-02 15:04")

		_, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			return results, fmt.Errorf("failed to marshal error log: %v", err)
		}
		logs = append(logs, result)
	}

	return logs, nil
}
