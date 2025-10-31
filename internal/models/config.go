package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Config struct {
	Id           bson.ObjectID `bson:"_id,omitempty"`
	StartTime    string        `bson:"start_time"`
	MessageSent  int           `bson:"message_sent"`
	FailedToSend int           `bson:"failed_to_send"`
	Direction    string        `bson:"direction"`
}

func GetConfig() (Config, error) {
	var result Config

	if collections == nil || collections["config"] == nil {
		return result, fmt.Errorf("database not connected. Please call ConnectDB() before querying")
	}

	filter := bson.D{{}}
	sorting := bson.D{{Key: "start_time", Value: -1}}
	err := collections["config"].FindOne(context.TODO(), filter, options.FindOne().SetSort(sorting)).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func CreateConfig() {
	if collections == nil || collections["config"] == nil {
		return
	}

	doc := Config{
		StartTime:    time.Now().Format(time.RFC3339),
		MessageSent:  0,
		FailedToSend: 0,
		Direction:    "",
	}

	collections["config"].InsertOne(context.TODO(), doc)
}

func UpdateConfig(field string) {
	config, err := GetConfig()
	if err != nil {
		return
	}

	filter := bson.D{{Key: "_id", Value: config.Id}}
	var update bson.D
	switch field {
	case "message_sent":
		update = bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: config.MessageSent + 1}, {Key: "direction", Value: "up"}}}}
	case "failed_to_send":
		update = bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: config.FailedToSend + 1}, {Key: "direction", Value: "down"}}}}
	case "start_time":
		update = bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: time.Now().Format(time.RFC3339)}}}}
	default:
		return
	}

	collections["config"].UpdateOne(context.TODO(), filter, update)
}
