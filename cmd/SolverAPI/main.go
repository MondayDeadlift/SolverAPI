package main

import (
	"SolverAPI/internal/app"
	"log"
)

func main() {
	server, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	if err := server.SetupDependencies(); err != nil {
		log.Fatalf("Failed to setup dependencies: %v", err)
	}

	log.Printf("Starting server on port %s", server.Config.HTTPServer.Port)
	if err := server.Start(); err != nil {
		log.Fatal("Server error: ", err)
	}
}

//curl http://localhost:8080/users/JanKlodVamBan - Информация пользователя (в конце твой ник на кодварс)
//pg_ctl restart -D "C:\Program Files\PostgreSQL\17\data" - рестарт PostgreSQL (перезапустить сервер без перезапуска PostgreSQL не получается)
//curl http://localhost:8080/katas/random - Случайная задача
