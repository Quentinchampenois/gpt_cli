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

const apiEndpoint = "https://api.openai.com/v1/images/generations"

type ImageResponse struct {
	Created int                 `json:"created"`
	Data    []map[string]string `json:"data"`
}

type ImageCommand struct{}

func (cmd *ImageCommand) GetApiKey() string {
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		return apiKey
	}

	return ""
}

func (*ImageCommand) Name() string     { return "image" }
func (*ImageCommand) Synopsis() string { return "Generate a new image from the given prompt" }
func (*ImageCommand) Usage() string {
	return `[HELP] image <prompt>
Generate an image from the given prompt.

example : 
$ image "a horse in a house"	
`
}

func (cmd *ImageCommand) SetFlags(fs *flag.FlagSet) {}

func (cmd *ImageCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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
	requestBody := strings.NewReader(fmt.Sprintf(`{"model": "image-alpha-001", "prompt": "%s"}`, prompt))
	request, err := http.NewRequest(http.MethodPost, apiEndpoint, requestBody)
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

	filename := fmt.Sprintf("./api/images/res-%s.json", time.Now().Format("2006-02-01-15-04-05"))
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
