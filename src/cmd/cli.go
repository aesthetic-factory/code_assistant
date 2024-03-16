package cmd

import (
	"bufio"
	"code_assistant/src/http_client"
	"fmt"
	"log"
	"os"
	"strings"
)

// StartCLI starts the command-line interface
func StartCLI() {
	fmt.Println("Welcome to Code Assistant!")
	fmt.Println("Type 'help' to see available options or 'exit' to quit.")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "exit" {
			fmt.Println("Exiting Code Assistant. Goodbye!")
			break
		}

		handleCommand(input)
	}
}

// handleCommand parses and handles user commands
func handleCommand(input string) {

	if input == "" {
		// overwrite empty input
		input = "help"
	}

	switch strings.ToLower(input) {
	case "help":
		fmt.Println("Available options:")
		fmt.Println(" - scan code")
		fmt.Println(" - list file")
		fmt.Println(" - code explanation")
		fmt.Println(" - exit")
	case "scan code":
		// Add implementation for scanning code
		fmt.Println("Scanning code...")
		ollamaTest()
	case "list file":
		// Add implementation for listing files
		fmt.Println("Listing files...")
		ollamaTest2()
	case "code explanation":
		// Add implementation for code explanation
		fmt.Println("Explaining code...")
	default:
		fmt.Println("Invalid command. Type 'help' to see available options.")
	}
}

func ollamaTest() {
	// Example input data
	req := http_client.TextGenRequest{
		URL:         "http://192.168.1.141:11434/api/generate",
		Model:       "dolphin-phi",
		Temperature: 0.1,
		Prompt:      "Hello, World!",
		Stream:      false,
	}

	// Call GenerateRemote function
	resp, err := http_client.TextGenerateRemote(req)
	if err != nil {
		log.Fatalf("Error calling GenerateRemote: %v", err)
	}

	// Print response
	fmt.Println("Result:", resp.Result)
	fmt.Println("Token:", resp.Token)
}

func ollamaTest2() {
	// Example input data
	req := http_client.EmbeddingRequest{
		URL:    "http://192.168.1.141:11434/api/embeddings",
		Model:  "nomic-embed-text",
		Prompt: "Hello, World!",
	}

	// Call GenerateRemote function
	resp, err := http_client.EmbeddingGenerateRemote(req)
	if err != nil {
		log.Fatalf("Error calling GenerateRemote: %v", err)
	}

	// Print response
	fmt.Println("Result:", resp.Result)
}
