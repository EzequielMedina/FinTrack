package main

import (
    "log"
    "os"

    "github.com/fintrack/chatbot-service/internal/app"
    "github.com/fintrack/chatbot-service/internal/config"
)

func main() {
    // health probe for Dockerfile
    if len(os.Args) > 1 && os.Args[1] == "health" {
        os.Exit(0)
    }

    cfg := config.Load()
    application, err := app.Initialize(cfg)
    if err != nil {
        log.Fatalf("failed to initialize application: %v", err)
    }

    if err := application.Start(); err != nil {
        log.Fatalf("application stopped with error: %v", err)
    }
}