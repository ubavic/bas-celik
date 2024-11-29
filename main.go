package main

import (
	"os"

	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/internal"
	"github.com/ubavic/bas-celik/internal/logger"

	"embed"
	"fmt"
)

//go:embed assets/version
var version string

//go:embed embed/*
var embedFS embed.FS

func main() {
	internal.SetVersion(version)

	cfg, exit := internal.ProcessFlags()
	if exit {
		return
	}

	if !cfg.Verbose {
		logger.DisableLog()
	}

	cfg.EmbedDirectory = embedFS

	err := configDocumentPackage()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	err = internal.Run(cfg)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(2)
	}
}

func configDocumentPackage() error {
	documentConfig := document.DocumentConfig{}
	var err error

	documentConfig.FontRegular, err = embedFS.ReadFile("embed/liberationSansRegular.ttf")
	if err != nil {
		return fmt.Errorf("reading font: %w", err)
	}

	documentConfig.FontBold, err = embedFS.ReadFile("embed/liberationSansBold.ttf")
	if err != nil {
		return fmt.Errorf("reading font: %w", err)
	}

	documentConfig.RfzoLogo, err = embedFS.ReadFile("embed/rfzo.png")
	if err != nil {
		return fmt.Errorf("reading font: %w", err)
	}

	err = document.Configure(documentConfig)
	if err != nil {
		return fmt.Errorf("setup error: %w", err)
	}

	return nil
}
