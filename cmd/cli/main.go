/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"os"

	"github.com/hculpan/areweplaying/cmd/cli/cmd"
	"github.com/hculpan/areweplaying/pkg/utils"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	logFile := os.Getenv("CLI_LOG_FILE")
	if len(logFile) > 0 {
		utils.SetLogFile(logFile)
	}

	cmd.Execute()
}
