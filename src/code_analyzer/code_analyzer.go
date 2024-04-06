package code_analyzer

import (
	"code_assistant/src/db"
	"code_assistant/src/fileutil"
	"code_assistant/src/http_client"
	"code_assistant/src/llm_prompt"
	"code_assistant/src/util"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"
)

// Entry point in code_analyzer
func AnalyzeDirectory(directory string) {
	ext := []string{".go", ".h", ".cpp", ".hpp", ".js", ".ts", ".py", ".java"}
	codeFilePaths, _ := fileutil.ScanFiles(directory, ext)
	// fmt.Println(codeFilePaths) // DEBUG

	for _, path := range codeFilePaths {
		fa, err := NewFunctionAnalyzer(path)
		if err != nil {
			log.Panicln(err)
		}
		if fa == nil {
			continue
		}
		fa.ScanFile()
	}
}

// GetCodeLanguage returns the programming language based on the file extension.
//
// filePath: a string representing the file path.
// string: the programming language based on the file extension.
func GetCodeLanguage(filePath string) string {
	switch {
	case strings.HasSuffix(filePath, ".cpp"):
		return "cpp"
	case strings.HasSuffix(filePath, ".h"):
		return "cpp"
	case strings.HasSuffix(filePath, ".hpp"):
		return "cpp"
	case strings.HasSuffix(filePath, ".js"):
		return "javascript"
	case strings.HasSuffix(filePath, ".ts"):
		return "typescript"
	case strings.HasSuffix(filePath, ".go"):
		return "golang"
	case strings.HasSuffix(filePath, ".java"):
		return "java"
	case strings.HasSuffix(filePath, ".py"):
		return "python"
	default:
		return "unkown type"
	}
}

func GetFileFromDb(filePath string) (int, string, error) {
	// get row from db where functionName matches
	rows, _ := db.GetDatabase().Query("SELECT id, sha256 FROM files WHERE file_path = ? LIMIT 1", filePath)

	defer rows.Close()

	// Iterate over the rows and print out the values
	for rows.Next() {
		var id int = -1
		var hash string = ""

		err := rows.Scan(&id, &hash)
		if err != nil {
			log.Fatal(err)
		}
		if id > -1 {
			return id, hash, nil
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return 0, "", err
	}
	return 0, "", fmt.Errorf("file not found")
}

type FileAnalyzer struct {
	FileId      int
	FilePath    string
	CodeSnippet []string
	LineStart   int
	LineEnd     int
	StepSize    int
	SHA256      string
}

// NewFunctionAnalyzer creates a new FileAnalyzer instance for the given file path.
//
// It checks if the file exists and returns an error if it does not. Then, it reads the file contents into memory.
// The file contents are concatenated into a single string and a SHA256 hash is generated for it.
// The hash is compared with the hash stored in the database to check if the file has already been analyzed.
// If the file has been analyzed, it returns nil. Otherwise, it inserts the file into the database and creates a new FileAnalyzer instance.
// The LineEnd field of the FileAnalyzer instance is set based on the StepSize field.
// If the length of the code snippet is less than the StepSize, LineEnd is set to the length of the code snippet.
// Otherwise, LineEnd is set to the StepSize.
// The function returns the created FileAnalyzer instance and nil error, or nil and an error if the file does not exist.
func NewFunctionAnalyzer(filePath string) (*FileAnalyzer, error) {
	if !fileutil.FileExists(filePath) {
		return nil, fmt.Errorf("file does not exist %s", filePath)
	}

	// Read file to memory
	codeSnippet, _ := fileutil.ReadFileLines(filePath)

	// Generate SHA256 hash
	concatenated := strings.Join(codeSnippet, "\n")
	hash := sha256.New()
	hash.Write([]byte(concatenated))
	hashedString := hex.EncodeToString(hash.Sum(nil))

	_, dbHash, _ := GetFileFromDb(filePath)
	if dbHash == hashedString {
		// file exists and already analyzed
		return nil, nil
	}

	// Get the current date and time
	currentTime := time.Now()

	// Try to insert the file into the database
	// function_name is unique, so there should only be one row
	db.GetDatabase().Execute("INSERT INTO files (file_path, sha256, last_update_datetime, rescan_required) VALUES (?, ?, ?, ?)", filePath, hashedString, currentTime, 0)

	// Get the ID of the inserted record
	dbFileId, _, _ := GetFileFromDb(filePath)

	fa := &FileAnalyzer{FileId: dbFileId, FilePath: filePath, CodeSnippet: codeSnippet, LineStart: 0, LineEnd: 0, StepSize: 100, SHA256: hashedString}

	// Set LineEnd
	if len(codeSnippet) < fa.StepSize {
		fa.LineEnd = len(codeSnippet)
	} else {
		fa.LineEnd = fa.StepSize
	}

	return fa, nil
}

// SlideWindow moves the window of the FileAnalyzer by the specified step.
//
// It takes an integer parameter for the step and does not return anything.
func (fa *FileAnalyzer) SlideWindow(step int) {
	fa.LineStart += step
	fa.LineEnd += step
}

// EnlargeWindow description of the Go function.
//
// It takes an int parameter 'step' and does not return anything.
func (fa *FileAnalyzer) EnlargeWindow(step int) {
	fa.LineEnd += step
}

// Entry point in FileAnalyzer
func (fa *FileAnalyzer) ScanFile() {

	fmt.Printf("Scanning file %s\n", fa.FilePath)
	fa.ScanContent()

	// Update the file in the database
	db.GetDatabase().Execute("UPDATE files SET sha256 = ?, last_update_datetime = ? WHERE id = ?", fa.SHA256, time.Now(), fa.FileId)
}

// Scan Content in a window
func (fa *FileAnalyzer) ScanContent() {

	// 1 Get Code Language
	language := GetCodeLanguage(fa.FilePath)

	// 2 Search For Class or Namespace
	// TODO: Implement

	// Get a default TextGenRequest struct
	req := http_client.NewTextGenRequest()

	// 3 Search For Functions
	prompt := llm_prompt.GetFunctionList(language, fa.CodeSnippet, fa.LineStart, fa.LineEnd)
	fmt.Printf("GetFunctionList\n%s\n\n", prompt) //DEBUG
	req.Prompt = prompt

	// Call TextGenerateRemote function
	resp, err := http_client.TextGenerateRemote(req)
	if err != nil {
		log.Fatalf("Error calling TextGenerateRemote: %v", err)
	}

	// Print response
	fmt.Println("Result:", resp.Result)
	fmt.Println("Token:", resp.Token)

	functions, _ := util.ParseJsonArray[llm_prompt.FunctionListItem](resp.Result)

	// Access the parsed objects
	for _, f := range functions {

		// if something goes wrong when parsing
		if f.FunctionName == "" {
			continue
		}

		validFunction, startLine, endLine := IdentifyFunction(fa, f.FunctionName, language)
		if !validFunction {
			continue
		}

		functionInfo, err := AnalyzeFunction(fa, language, f.FunctionName, startLine, endLine)
		if err != nil {
			continue
		}

		// Remove old record if exists
		db.GetDatabase().Execute(`DELETE FROM functions WHERE function_name = ? AND file_id = ?`, f.FunctionName, fa.FileId)

		db.GetDatabase().Execute(`INSERT INTO functions (function_name, signature, arguments, return, namespace, description, file_id, line_start, line_end) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			f.FunctionName, functionInfo.Signature, functionInfo.Arguments, functionInfo.Return, "NONE", functionInfo.Purpose, fa.FileId, startLine, endLine)
	}
}
