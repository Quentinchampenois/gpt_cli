# OpenAI CLI

Command Line Interface that allows to communicate with OpenAI API.

I am currently learning the Go programming language, it may contain not recommended practice, don't hesitate to open a pull request if you want !

## Overview

This project provide a CLI to communicate with OpenAI API. 

### Commands

CLI allows to interact with different API endpoint : 

**Images**
  ```
    gpt_api image generate "<TYPE YOUR PROMPT>"
  ```
It returns a list of urls of your generated image and save in folder `api/images/`


## Getting started

### Requirements
* OpenAI API Key, [see documentation](https://platform.openai.com/overview)
* Go version `> 1.20`


## Resources

* ChatGPT
* [Osinet CLI with google/subcommands](https://osinet.fr/go/articles/cli-google-subcommands/)
