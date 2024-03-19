package code_analyzer

import (
	"code_assistant/src/fileutil"
	"code_assistant/src/http_client"
	"code_assistant/src/llm_prompt"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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
		fa.ScanFile()
		break
	}
}

// Helper function
// GetCodeLanguage returns the file extension for the given file path.
// It maps common file extensions to their corresponding programming languages.
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

// Helper function
func ParseJsonObject[T any](input string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		return result, fmt.Errorf("error parsing JSON %s", err)
	}
	return result, nil
}

// Helper function
func ParseJsonArray[T any](input string) ([]T, error) {

	// Parse the JSON data into a slice of structs
	var results []T
	if err := json.Unmarshal([]byte(input), &results); err != nil {
		return results, fmt.Errorf("error parsing JSON %s", err)
	}
	return results, nil
}

type FileAnalyzer struct {
	FilePath    string
	CodeSnippet []string
	LineStart   int
	LineEnd     int
	StepSize    int
}

func NewFunctionAnalyzer(filePath string) (*FileAnalyzer, error) {
	if !fileutil.FileExists(filePath) {
		return nil, fmt.Errorf("file does not exist %s", filePath)
	}

	// Read file to memory
	codeSnippet, _ := fileutil.ReadFileLines(filePath)
	fa := &FileAnalyzer{FilePath: filePath, CodeSnippet: codeSnippet, LineStart: 0, LineEnd: 0, StepSize: 42}

	// Set LineEnd
	if len(codeSnippet) < fa.StepSize {
		fa.LineEnd = len(codeSnippet)
	} else {
		fa.LineEnd = fa.StepSize
	}

	return fa, nil
}

func (fa *FileAnalyzer) SlideWindow(step int) {
	fa.LineStart += step
	fa.LineEnd += step
}

func (fa *FileAnalyzer) EnlargeWindow(step int) {
	fa.LineEnd += step
}

// Entry point in FileAnalyzer
func (fa *FileAnalyzer) ScanFile() {

	fa.ScanContent()
}

// Scan Content in a window
func (fa *FileAnalyzer) ScanContent() {

	// 1 Get Code Language
	language := GetCodeLanguage(fa.FilePath)

	// 2 Search For Class or Namespace

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

	functions, _ := ParseJsonArray[llm_prompt.FunctionListItem](resp.Result)

	// Access the parsed objects
	for _, f := range functions {

		// if something goes wrong when parsing
		if f.FunctionName == "" {
			continue
		}

		AnalyzeFunction(fa, f.FunctionName, language)
	}

	// re := regexp.MustCompile("```json\\s+(.+?)\\s+```")
	// match := re.FindStringSubmatch(resp.Result)
	// if len(match) >= 2 {
	// err := json.Unmarshal([]byte(match[1]), &jsonObj)

	// 2 Search For Functions
}
