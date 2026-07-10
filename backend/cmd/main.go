package main

import (
	"fmt"
	"log"

	env "sigefae/internal/.env"
	"sigefae/internal/graph"
)

func main() {
	// Cargar variables de entorno
	cfg, err := env.Load(".env")
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	fmt.Println("Environment variables loaded successfully.")

	// Crear autenticación
	auth := graph.NewAuth(
		cfg.GraphClientID,
		cfg.GraphClientSecret,
		cfg.GraphTenantID,
	)

	// Crear cliente de Graph
	client := graph.NewClient(
		auth,
		cfg.GraphUserEmail,
	)

	// Obtener correos
	messages, err := client.ListMessages()
	if err != nil {
		log.Fatal(err)
	}

	for _, msg := range messages {
		fmt.Println("--------------------------------")
		fmt.Println("Asunto:", msg.Subject)

		if !msg.HasAttachments {
			fmt.Println("Sin adjuntos")
			continue
		}

		if err := client.DownloadAttachments(msg); err != nil {
			log.Println("Error descargando adjuntos:", err)
		}
	}
}