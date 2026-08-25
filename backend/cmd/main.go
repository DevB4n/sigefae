package main

import (
	"fmt"
	"log"

	"sigefae/env"
	"sigefae/internal/api"
	"sigefae/internal/db"
	"sigefae/procesos_aplicacion/descarga_correos"
)

func main() {

	// ========================================
	// Configuración
	// ========================================

	cfg, err := env.Load("../../../../../.env_hp")
	if err != nil {
		log.Fatal(err)
	}

	// ========================================
	// Base de datos
	// ========================================

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	if err := db.Seed(db.DB); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Environment loaded successfully.")
	fmt.Println("Database connected successfully.")

	// ========================================
	// API & Tareas en Segundo Plano
	// ========================================

	// Iniciar la tarea de descarga de correos en segundo plano
	go descarga_correos.Start(cfg, database)

	router := api.New(database)

	log.Println("Servidor iniciado en http://localhost:8080")

	log.Fatal(router.Run(":8080"))
}
