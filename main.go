package main

import (
	"signal/main/internal"
	"signal/main/internal/models"
	"signal/main/internal/services"
)

func main() {
	models.ConnectDB()
	go services.ForwardMessengerMessages()
	internal.StartServer()
	defer models.DisconnectDB()
}
