package main

import (
	"code_assistant/src/cmd"
	"code_assistant/src/db"
	"log"
)

func main() {

	// Create a new instance of the Database
	database, err := db.NewDatabase("./local.db")
	if err != nil {
		log.Panic(err)
	}
	defer database.Close()

	// Create Tables

	err = database.CreateTable("files",
		`id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path TEXT NOT NULL UNIQUE,
		hash TEXT NOT NULL,
		last_update_datetime TEXT NOT NULL,
		rescan_required INT NOT NULL`)

	if err != nil {
		log.Panic(err)
	}

	err = database.CreateTable("functions",
		`id INTEGER PRIMARY KEY AUTOINCREMENT, 
		function_name TEXT NOT NULL UNIQUE, 
		arguments TEXT NOT NULL, 
		return TEXT NOT NULL, 
		namespace TEXT NOT NULL,
		description TEXT NOT NULL,
		file_id INT NOT NULL,
		line_start INT NOT NULL,
		line_end INT NOT NULL,
		hash TEXT NOT NULL,
		FOREIGN KEY(file_id) REFERENCES files(id)`)

	if err != nil {
		log.Panic(err)
	}

	// Insert a sample user
	_, err = database.Execute("INSERT INTO files (file_path, hash, last_update_datetime, rescan_required) VALUES (?, ?, ?, ?)", "sample_user", "sample_password", "date", 0)
	if err != nil {
		log.Println(err)
	}

	// Start the command-line interface
	cmd.StartCLI()
}
