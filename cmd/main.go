package main

import (
	"github/Danila331/YAP/internal/models"
	"github/Danila331/YAP/internal/server"
	"github/Danila331/YAP/internal/store"
	"log"
)

func main() {

	conn, err := store.ConnectToDatabase()
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer conn.Close()

	// Создаем оркестратора
	orchestrator := models.NewOrchestrator(conn)
	orchestrator.Start()
	server.StartServer()
	if err != nil {
		log.Println(err)
	}
}
