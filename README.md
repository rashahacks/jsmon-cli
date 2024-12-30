GoLang API Client
This repository contains a Go-based API client for interacting with various endpoints. The client is designed to be lightweight, easy to use, and extendable for your specific needs.

Installation
Follow the steps below to install the jsmon-cli binary:

Open your terminal and run the following command to install the client using Go:

bash >> 
go install github.com/rashahacks/jsmon-cli@latest


Once the installation is complete, the binary will be available in your $HOME/go/bin folder. You can verify its existence by navigating to the directory:

bash>>
cd $HOME/go/bin
ls | grep jsmon-cli

Usage
After installation, you can start using the jsmon-cli command-line tool.

To check the available commands and options, run:

bash>> 

jsmon-cli --help
This will display a list of all available flags and commands.

Example Command
Here’s an example of how to use the client:

bash>>
jsmon-cli -u <JS-URL-HERE>

For more details about specific commands, refer to the tool's help output or documentation.

