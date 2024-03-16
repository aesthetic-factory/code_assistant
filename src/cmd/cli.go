package cmd

import (
	"bufio"
	"fmt"
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
	case "list file":
		// Add implementation for listing files
		fmt.Println("Listing files...")
	case "code explanation":
		// Add implementation for code explanation
		fmt.Println("Explaining code...")
	default:
		fmt.Println("Invalid command. Type 'help' to see available options.")
	}
}
