package main

import (
	"fmt"
	"log"

	env "sigefae/internal/.env"
	"sigefae/internal/graph"
	"sigefae/internal/sync"
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

	// Crear cliente de Microsoft Graph
	client := graph.NewClient(
		auth,
		cfg.GraphUserEmail,
	)

	// Servicio de sincronización
	syncService := sync.New()

	// Obtener la fecha del último correo sincronizado
	lastSync := syncService.LastSync()

	// Consultar únicamente los correos nuevos
	messages, err := client.ListMessages(lastSync)
	if err != nil {
		log.Fatal(err)
	}

	// Procesar correos
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

	// Actualizar la fecha del último correo procesado
	if len(messages) > 0 {
		if err := syncService.Update(
			messages[len(messages)-1].ReceivedDateTime,
		); err != nil {
			log.Println(err)
		}
	}
}