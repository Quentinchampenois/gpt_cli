# OpenAI CLI

Command Line Interface that allows to communicate with OpenAI API.

I am currently learning the Go programming language, it may contain not recommended practice, don't hesitate to open a pull request if you want !

## Overview

This project provide a CLI to communicate with OpenAI API. 

### Commands

CLI allows to interact with different API endpoint : 

**Images**
  ```
    gpt_api image "<TYPE YOUR PROMPT>"
  ```
It returns a list of urls of your generated image and save in folder `api/images/`

Example:
```bash
./gpt_api image "A mouse driving a car"
```


**Completion**
  ```
    gpt_api program "<TYPE YOUR PROMPT>"
  ```
It returns a list of urls of your generated program and save in folder `api/completions/`

Example:
```bash
./gpt_api program "Create a Ruby REST API using Sinatra with JWT authentication"
```

**Correction**
  ```
    gpt_api correct "<TYPE YOUR PROMPT>"
  ```
It returns a edited version of given input in `api/corrects/`

Example:
```bash
./gpt_api correct "I wuld like to corect this input"
```

## Getting started

### Requirements
* OpenAI API Key, [see documentation](https://platform.openai.com/overview)
* Go version `> 1.20`

### Run locally

Build project and use cli 

`go build . && ./gpt_api`

## Resources

* ChatGPT
* [Osinet CLI with google/subcommands](https://osinet.fr/go/articles/cli-google-subcommands/)
* OpenAI Documentation
  * [Models](https://platform.openai.com/docs/models/model-endpoint-compatibility)
  * [API Reference](https://platform.openai.com/docs/api-reference/chat/create)