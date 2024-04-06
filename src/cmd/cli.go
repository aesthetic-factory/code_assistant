package cmd

import (
	"bufio"
	"code_assistant/src/code_analyzer"
	"code_assistant/src/config"
	"code_assistant/src/db"
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
		fmt.Println(" - list function")
		fmt.Println(" - code explanation")
		fmt.Println(" - exit")

	case "scan code":
		var directory string
		fmt.Print("Enter directory to scan: ")
		if _, err := fmt.Scanln(&directory); err != nil {
			directory = config.AppConfig.WorkingDir
		}
		if directory == "" {
			fmt.Println("directory cannot be empty")
			return
		}
		fmt.Printf("Scanning directory %s ...\n", directory)
		code_analyzer.AnalyzeDirectory(directory)

	case "list file":
		fmt.Println("Listing files...")
		listFiles()

	case "list function":
		fmt.Println("Listing functions...")
		listFunctions()

	case "code explanation":
		// Add implementation for code explanation
		fmt.Println("Explaining code...")

	default:
		fmt.Println("Invalid command. Type 'help' to see available options.")
	}
}

func listFiles() {

	rows, err := db.GetDatabase().Query("SELECT id, file_path, last_update_datetime FROM files")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	// Iterate over the rows and print out the values
	for rows.Next() {
		var id int
		var file_path string
		var last_update_datetime string
		err := rows.Scan(&id, &file_path, &last_update_datetime)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, file_path: %s, last_update_datetime: %s\n\n",
			id, file_path, last_update_datetime)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}

func listFunctions() {

	rows, err := db.GetDatabase().Query("SELECT a.id, function_name, signature, arguments, return, description, b.file_path, line_start, line_end FROM functions a JOIN files b ON a.file_id = b.id")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	// Iterate over the rows and print out the values
	for rows.Next() {
		var id int
		var function_name string
		var signature string
		var arguments string
		var return_type string
		var description string
		var file_path string
		var line_start int
		var line_end int
		err := rows.Scan(&id, &function_name, &signature, &arguments, &return_type, &description, &file_path, &line_start, &line_end)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, function_name: %s, file_path: %s\ndescription: %s\nsignature: %s\narguments: %s\nreturn_type: %s\nline_start: %d\nline_end: %d\n\n",
			id, function_name, file_path, description, signature, arguments, return_type, line_start, line_end)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}
