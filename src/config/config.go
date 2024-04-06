package config

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Ollama struct {

	// Define Ollama configuration variables here
	BaseUrl string

	TextGenModel   string
	ChatModel      string
	EmbeddingModel string
}

type Config struct {

	// Define base configuration variables here
	Ollama     Ollama
	DebugMode  bool
	DbFilePath string

	WorkingDir string
}

var AppConfig Config

// LoadEnv loads environment variables from a .env file.
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		// .env file not found or cannot be read
		log.Printf("Failed to load .env file: %v", err)
		return
	}
}

// LoadConfig loads the configuration settings from command-line flags and environment variables.
//
// No parameters.
// No return value.
func LoadConfig() {

	// Load environment variables from a .env file
	loadEnv()

	// Define flags for command-line arguments
	ollamaBaseUrl := flag.String("ollama_base_url", getEnv("OLLAMA_BASE_URL", "http://127.0.0.1:11434"), "Ollama Base URL")
	ollamaTextGenModel := flag.String("ollama_textgen_model", getEnv("OLLAMA_TEXTGEN_MODEL", "mistral:instruct"), "Ollama Text Generation Model")
	ollamaChatModel := flag.String("ollama_chat_model", getEnv("OLLAMA_CHAT_MODEL", "mistral:instruct"), "Ollama Chat Model")
	ollamaEmbeddingModel := flag.String("ollama_embedding_model", getEnv("OLLAMA_EMBEDDING_MODEL", "nomic-embed-text:latest"), "Ollama Embedding Model")

	debugMode := flag.Bool("debug", getBoolEnv("DEBUG_MODE", false), "Enable debug mode")
	dbFilePath := flag.String("db_filepath", getEnv("DB_FILEPATH", "./local.db"), "Database File Path")
	workingDir := flag.String("working_dir", getEnv("WORKING_DIR", ""), "Working Directory for Code Base")

	// Parse command-line arguments
	flag.Parse()

	// Set configuration from command-line flags
	AppConfig.Ollama.BaseUrl = *ollamaBaseUrl
	AppConfig.Ollama.TextGenModel = *ollamaTextGenModel
	AppConfig.Ollama.ChatModel = *ollamaChatModel
	AppConfig.Ollama.EmbeddingModel = *ollamaEmbeddingModel

	AppConfig.DebugMode = *debugMode
	AppConfig.DbFilePath = *dbFilePath
	AppConfig.WorkingDir = *workingDir
}

// getEnv gets the value of the environment variable with the specified key.
// If the variable is not set, it returns the default value.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// getIntEnv gets the value of the environment variable with the specified key
// and converts it to an integer. If the variable is not set or cannot be parsed
// as an integer, it returns the default value.
func getIntEnv(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getBoolEnv gets the value of the environment variable with the specified key
// and converts it to a boolean. If the variable is not set or cannot be parsed
// as a boolean, it returns the default value.
func getBoolEnv(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
