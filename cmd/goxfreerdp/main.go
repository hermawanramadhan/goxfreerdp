package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"goxfreerdp/internal/config"
	"goxfreerdp/internal/ui"
)

func main() {
	_, err := config.EnsureConfigDir()
	if err != nil {
		log.Fatalf("Error ensuring config directory: %v\n", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	var initialRDPFile string
	if len(os.Args) > 1 {
		initialRDPFile = os.Args[1]
		if absPath, err := filepath.Abs(initialRDPFile); err == nil {
			initialRDPFile = absPath
		}
	}

	if ui.TrySingleInstance(initialRDPFile) {
		fmt.Println("Sent connection request to existing instance. Exiting.")
		os.Exit(0)
	}

	// Save default configuration if the file does not exist yet
	filePath, err := config.GetConfigFilePath()
	if err != nil {
		log.Fatalf("Error getting config path: %v\n", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = config.SaveConfig(cfg)
		if err != nil {
			log.Printf("Failed to save initial config: %v\n", err)
		}
	}

	// Start GTK3 GUI Mode (with optional initial RDP file)
	err = ui.StartApp(&cfg, initialRDPFile)
	if err != nil {
		log.Fatalf("Failed to build main UI: %v\n", err)
	}
}
