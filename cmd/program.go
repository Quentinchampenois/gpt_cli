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

const completionApiEndpoint = "https://api.openai.com/v1/completions"

type ProgramCommand struct{}

func (cmd *ProgramCommand) GetApiKey() string {
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		return apiKey
	}

	return ""
}

func (*ProgramCommand) Name() string     { return "program" }
func (*ProgramCommand) Synopsis() string { return "Generate a new program from the given prompt" }
func (*ProgramCommand) Usage() string {
	return `[HELP] program <prompt>
Generate an image from the given prompt.

example : 
$ image "a horse in a house"	
`
}

func (cmd *ProgramCommand) SetFlags(fs *flag.FlagSet) {}

func (cmd *ProgramCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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
	requestBody := strings.NewReader(fmt.Sprintf(`{"model": "text-davinci-003", "prompt": "%s",  "temperature":0,
  "max_tokens":1000,
  "top_p":1.0,
  "frequency_penalty":0.2,
  "presence_penalty":0.0}`, prompt))
	request, err := http.NewRequest(http.MethodPost, completionApiEndpoint, requestBody)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating API request: %v\n", err)
		return subcommands.ExitFailure
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cmd.GetApiKey()))

	client := &http.Client{}
	log.Println("Sending POST request on OpenAI Completion Endpoint...")
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

	filename := fmt.Sprintf("./api/completions/res-%s.json", time.Now().Format("2006-02-01-15-04-05"))
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
