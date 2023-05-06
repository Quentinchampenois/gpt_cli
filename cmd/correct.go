package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const correctApiEndpoint = "https://api.openai.com/v1/edits"

type CorrectCommand struct{}

func (cmd *CorrectCommand) GetApiKey() string {
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		return apiKey
	}

	return ""
}

func (*CorrectCommand) Name() string     { return "correct" }
func (*CorrectCommand) Synopsis() string { return "Generate a new correct from the given prompt" }
func (*CorrectCommand) Usage() string {
	return `[HELP] correct <prompt>
Generate an correct from the given prompt.

example : 
$ correct "a horse in a house"	
`
}

func (cmd *CorrectCommand) SetFlags(fs *flag.FlagSet) {}

func (cmd *CorrectCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Missing prompt argument, please refer to the help section : ")
		fmt.Fprintln(os.Stderr, cmd.Usage())
		return subcommands.ExitUsageError
	}

	prompt := f.Arg(0)

	if cmd.GetApiKey() == "" {
		fmt.Fprintln(os.Stderr, "Error: OPENAI_API_KEY environment variable not set.")
		return subcommands.ExitFailure
	}
	requestBody := strings.NewReader(fmt.Sprintf(`{"model": "text-davinci-edit-001", "input": "%s", "instruction": "Corrige les erreurs d'orthographe et de grammaire'"}`, prompt))
	request, err := http.NewRequest(http.MethodPost, correctApiEndpoint, requestBody)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating API request: %v\n", err)
		return subcommands.ExitFailure
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cmd.GetApiKey()))

	client := &http.Client{}
	log.Println("Sending POST request on OpenAI Image Endpoint...")
	response, err := client.Do(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending API request: %v\n", err)
		return subcommands.ExitFailure
	}
	defer response.Body.Close()

	// Print API response
	if response.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "API request failed with status code %d\n", response.StatusCode)
		return subcommands.ExitFailure
	}

	log.Println("Response successfully received")

	filename := fmt.Sprintf("./api/corrects/res-%s.json", time.Now().Format("2006-02-01-15-04-05"))
	resFile, err := os.Create(filename)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		return subcommands.ExitFailure
	}
	defer func(resFile *os.File) {
		err := resFile.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		}
	}(resFile)

	_, err = io.Copy(resFile, response.Body)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		return subcommands.ExitFailure
	}
	log.Println(fmt.Sprintf("Saving response in file : %s", filename))
	log.Println("End of command")

	return subcommands.ExitSuccess
}
