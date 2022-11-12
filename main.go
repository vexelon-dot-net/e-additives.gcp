package main

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vexelon-dot-net/e-additives.gcp/config"
	"github.com/vexelon-dot-net/e-additives.gcp/service"
)

func main() {
	const HEART = "\u2764"

	fmt.Printf("%s e-additives API service v%s %s\n\n", HEART, config.VERSION, HEART)

	config, err := config.CreateFromEnv()
	if err != nil {
		log.Fatalf("Config error: %v\n", err)
	}

	if config.IsDevMode {
		fmt.Println("DEV mode enabled.")
	}

	server := service.NewServer(config)
	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
