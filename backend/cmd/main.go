package main

import (
	"fmt"
	"log"

	"sigefae/internal/db"
	"sigefae/internal/env"
	"sigefae/internal/graph"
	"sigefae/internal/sync"
)

func main() {

	// ========================================
	// Configuración
	// ========================================
	cfg, err := env.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// ========================================
	// Base de datos
	// ========================================
	if _, err := db.Connect(cfg); err != nil {
		log.Fatal(err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Environment loaded successfully.")
	fmt.Println("Database connected successfully.")

	// ========================================
	// Cliente Microsoft Graph
	// ========================================
	auth := graph.NewAuth(
		cfg.GraphClientID,
		cfg.GraphClientSecret,
		cfg.GraphTenantID,
	)

	client := graph.NewClient(
		auth,
		cfg.GraphUserEmail,
	)

	// ========================================
	// Sincronización
	// ========================================
	syncService := sync.New()

	lastSync := syncService.LastSync()

	messages, err := client.ListMessages(lastSync)
	if err != nil {
		log.Fatal(err)
	}

	newestSync := lastSync

	for _, msg := range messages {

		fmt.Println("--------------------------------")
		fmt.Println("Asunto:", msg.Subject)

		if !msg.HasAttachments {

			fmt.Println("Sin adjuntos")

			if msg.ReceivedDateTime.After(newestSync) {
				newestSync = msg.ReceivedDateTime
			}

			continue
		}

		if err := client.DownloadAttachments(msg); err != nil {
			log.Println(err)
			continue
		}

		if msg.ReceivedDateTime.After(newestSync) {
			newestSync = msg.ReceivedDateTime
		}
	}

	if !newestSync.IsZero() {
		if err := syncService.Update(newestSync); err != nil {
			log.Println(err)
		}
	}
}