package descarga_correos

import (
	"log"
	"time"

	"gorm.io/gorm"

	"sigefae/env"
	"sigefae/internal/correo"
	"sigefae/procesos_aplicacion/crear_correos"
	"sigefae/procesos_aplicacion/descarga_correos/graph"
	"sigefae/procesos_aplicacion/descarga_correos/sync"
)

func Start(cfg *env.Env, database *gorm.DB) {
	// Definimos el intervalo en 5 minutos
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	auth := graph.NewAuth(cfg.GraphClientID, cfg.GraphClientSecret, cfg.GraphTenantID)
	client := graph.NewClient(auth, cfg.GraphUserEmail)
	syncService := sync.New()
	correoService := correo.New(database)

	log.Println("[Descarga Correos] Servicio iniciado en segundo plano (5 min interval)")

	// Ejecutar la primera sincronización inmediatamente
	if err := processMails(client, syncService, correoService); err != nil {
		log.Printf("[Descarga Correos] Error sincronizando correos: %v\n", err)
	}

	for {
		<-ticker.C
		log.Println("[Descarga Correos] Iniciando sincronización de correos...")
		if err := processMails(client, syncService, correoService); err != nil {
			log.Printf("[Descarga Correos] Error sincronizando correos: %v\n", err)
		} else {
			log.Println("[Descarga Correos] Sincronización finalizada correctamente.")
		}
	}
}

func processMails(client *graph.Client, syncService *sync.Service, correoService *correo.Service) error {
	lastSync := syncService.LastSync()

	messages, err := client.ListMessages(lastSync)
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return nil
	}

	var latestSync time.Time = lastSync

	for _, msg := range messages {
		if msg.ReceivedDateTime.After(latestSync) {
			latestSync = msg.ReceivedDateTime
		}

		if err := client.DownloadAttachments(msg); err != nil {
			log.Printf("Error descargando adjuntos para el correo %s: %v\n", msg.ID, err)
			continue
		}

		if err := crear_correos.ProcesarYCrear(msg, correoService); err != nil {
			log.Printf("Error registrando el correo %s en la base de datos: %v\n", msg.ID, err)
			continue
		}
	}

	if latestSync.After(lastSync) {
		if err := syncService.Update(latestSync); err != nil {
			return err
		}
	}

	return nil
}
