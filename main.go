package main

import (
	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/internal"

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

	err := document.SetData(embedFS)
	if err != nil {
		fmt.Println("Setup error:", err)
		return
	}

	err = internal.Run(cfg)
	if err != nil {
		fmt.Println("Error saving document:", err)
	}
}
