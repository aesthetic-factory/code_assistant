package main

import (
	"code_assistant/src/cmd"
	"code_assistant/src/config"
	"code_assistant/src/db"
	"log"
)

func main() {

	// Load config
	config.LoadConfig()

	// Create a new instance of the Database
	database, err := db.NewDatabase(config.AppConfig.DbFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer database.Close()

	// Create Tables

	err = database.CreateTable("files",
		`id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path TEXT NOT NULL UNIQUE,
		sha256 TEXT NOT NULL,
		last_update_datetime DATETIME NOT NULL,
		rescan_required INT NOT NULL`)

	if err != nil {
		log.Panic(err)
	}

	err = database.CreateTable("functions",
		`id INTEGER PRIMARY KEY AUTOINCREMENT, 
		function_name TEXT NOT NULL UNIQUE, 
		signature TEXT NOT NULL,
		arguments TEXT NOT NULL,
		return TEXT NOT NULL, 
		namespace TEXT NOT NULL,
		description TEXT NOT NULL,
		file_id INT NOT NULL,
		line_start INT NOT NULL,
		line_end INT NOT NULL,
		FOREIGN KEY(file_id) REFERENCES files(id)`)

	if err != nil {
		log.Panic(err)
	}

	// Start the command-line interface
	cmd.StartCLI()
}
