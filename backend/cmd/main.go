package main

import (
	"fmt"
	"log"

	env "sigefae/internal/.env"
	"sigefae/internal/graph"
	"sigefae/internal/sync"
)

func main() {
	cfg, err := env.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Environment variables loaded successfully.")

	auth := graph.NewAuth(
		cfg.GraphClientID,
		cfg.GraphClientSecret,
		cfg.GraphTenantID,
	)

	client := graph.NewClient(
		auth,
		cfg.GraphUserEmail,
	)

	syncService := sync.New()

	lastSync := syncService.LastSync()

	messages, err := client.ListMessages(lastSync)
	if err != nil {
		log.Fatal(err)
	}

	var newestSync = lastSync

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