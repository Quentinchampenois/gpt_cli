package main

import (
	"context"
	"github.com/joho/godotenv"
	"gpt_api/cmd"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error while loading '.env': %s", err)
	}

	ctx := context.Background()
	os.Exit(int(cmd.Execute(ctx)))
}
