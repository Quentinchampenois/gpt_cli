package main

import (
	"context"
	"gpt_api/cmd"
	"os"
)

func main() {
	ctx := context.Background()
	os.Exit(int(cmd.Execute(ctx)))
}
