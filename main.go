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

	var jsonPath, pdfPath string
	var verboseFlag, getValidUntilFromRfzo bool
	var readerIndex uint
	exit := internal.ProcessFlags(&jsonPath, &pdfPath, &verboseFlag, &getValidUntilFromRfzo, &readerIndex)
	if exit {
		return
	}

	err := document.SetData(embedFS)
	if err != nil {
		fmt.Println("Setup error:", err)
		return
	}

	err = internal.Run(pdfPath, jsonPath, verboseFlag, getValidUntilFromRfzo, readerIndex)
	if err != nil {
		fmt.Println("Error saving document:", err)
	}
}
