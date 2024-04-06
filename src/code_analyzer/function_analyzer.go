package code_analyzer

import (
	"code_assistant/src/http_client"
	"code_assistant/src/llm_prompt"
	"code_assistant/src/util"
	"fmt"
	"log"
)

type Function struct {
	LineStart int
	LineEnd   int
}

func IdentifyFunction(fa *FileAnalyzer, functionName string, language string) (bool, int, int) {

	// Get a default TextGenRequest struct
	chatReq := http_client.NewChatRequest()

	startLine := 0
	endLine := 0

	// Locate Function Prompt
	{
		prompt := llm_prompt.LocateFunctionDefition(functionName, language, fa.CodeSnippet, fa.LineStart, fa.LineEnd)
		fmt.Printf("LocateFunctionDefition\n%s\n\n", prompt) //DEBUG

		chatReq.Messages = append(chatReq.Messages, http_client.Chat{Role: "user", Content: prompt})

		// Call ChatGenerateRemote function
		resp, err := http_client.ChatGenerateRemote(chatReq)
		if err != nil {
			log.Fatalf("Error calling ChatGenerateRemote: %v", err)
		}
		fmt.Println("Role:", resp.Result.Role)
		fmt.Println("Content:", resp.Result.Content)

		chatReq.Messages = append(chatReq.Messages, resp.Result) // append to the messages
	}

	// Locate Function Prompt
	{
		prompt := llm_prompt.LocateFunctionDefitionFinal(functionName, language, fa.CodeSnippet, fa.LineStart, fa.LineEnd)
		fmt.Printf("LocateFunctionDefitionFinal\n%s\n\n", prompt) //DEBUG
		chatReq.Messages = append(chatReq.Messages, http_client.Chat{Role: "user", Content: prompt})

		// Call ChatGenerateRemote function
		resp, err := http_client.ChatGenerateRemote(chatReq)
		if err != nil {
			log.Fatalf("Error calling ChatGenerateRemote: %v", err)
		}

		// Print response
		fmt.Println("Role:", resp.Result.Role)
		fmt.Println("Content:", resp.Result.Content)

		// no need to append to the messages
		chatReq.Messages = append(chatReq.Messages, resp.Result) // append to the messages

		res, err := util.ParseJsonObject[llm_prompt.LocateFunctionResponse](resp.Result.Content)
		if err != nil {
			log.Fatalf("Error LocateFunctionDefitionFinal ParseJsonObject: %v", err)
			return false, startLine, endLine
		}
		startLine = res.StartLine
		endLine = res.EndLine
	}

	// Check Function Defition Prompt
	{
		prompt := llm_prompt.CheckFunctionDefition(functionName, language, fa.CodeSnippet, fa.LineStart, fa.LineEnd)
		fmt.Printf("CheckFunctionDefition\n%s\n\n", prompt) //DEBUG
		chatReq.Messages = append(chatReq.Messages, http_client.Chat{Role: "user", Content: prompt})

		// Call ChatGenerateRemote function
		resp, err := http_client.ChatGenerateRemote(chatReq)
		if err != nil {
			log.Fatalf("Error calling ChatGenerateRemote: %v", err)
		}

		// Print response
		fmt.Println("Role:", resp.Result.Role)
		fmt.Println("Content:", resp.Result.Content)
		chatReq.Messages = append(chatReq.Messages, resp.Result)
	}

	// Check Function Defition Prompt
	{
		prompt := llm_prompt.CheckFunctionDefitionFinal(functionName, language, fa.CodeSnippet, fa.LineStart, fa.LineEnd)
		fmt.Printf("CheckFunctionDefitionFinal\n%s\n\n", prompt) //DEBUG
		chatReq.Messages = append(chatReq.Messages, http_client.Chat{Role: "user", Content: prompt})

		// Call ChatGenerateRemote function
		resp, err := http_client.ChatGenerateRemote(chatReq)
		if err != nil {
			log.Fatalf("Error calling ChatGenerateRemote: %v", err)
		}

		// Print response
		fmt.Println("Role:", resp.Result.Role)
		fmt.Println("Content:", resp.Result.Content)
		res, err := util.ParseJsonObject[llm_prompt.BooleanItem](resp.Result.Content)
		if err != nil {
			log.Fatalf("Error CheckFunctionDefitionFinal ParseJsonObject: %v", err)
			return false, startLine, endLine
		}
		return res.Result, startLine, endLine
	}
}

func AnalyzeFunction(fa *FileAnalyzer, language string, functionName string, functionStartLine int, functionEndLine int) (*llm_prompt.AnalyzeFunctionResponse, error) {

	// Get a default TextGenRequest struct
	chatReq := http_client.NewChatRequest()

	{
		prompt := llm_prompt.AnalyzeFunction(functionName, language, fa.CodeSnippet, fa.LineStart, fa.LineEnd)
		fmt.Printf("AnalyzeFunction\n%s\n\n", prompt) //DEBUG

		chatReq.Messages = append(chatReq.Messages, http_client.Chat{Role: "user", Content: prompt})

		// Call ChatGenerateRemote function
		resp, err := http_client.ChatGenerateRemote(chatReq)
		if err != nil {
			log.Fatalf("Error calling ChatGenerateRemote: %v", err)
		}
		fmt.Println("Role:", resp.Result.Role)
		fmt.Println("Content:", resp.Result.Content)

		chatReq.Messages = append(chatReq.Messages, resp.Result)
	}

	{
		prompt := llm_prompt.AnalyzeFunctionFinal(functionName, language, fa.CodeSnippet, fa.LineStart, fa.LineEnd)
		fmt.Printf("AnalyzeFunction\n%s\n\n", prompt) //DEBUG

		chatReq.Messages = append(chatReq.Messages, http_client.Chat{Role: "user", Content: prompt})

		// Call ChatGenerateRemote function
		resp, err := http_client.ChatGenerateRemote(chatReq)
		if err != nil {
			log.Fatalf("Error calling ChatGenerateRemote: %v", err)
		}
		fmt.Println("Role:", resp.Result.Role)
		fmt.Println("Content:", resp.Result.Content)

		res, err := util.ParseJsonObject[llm_prompt.AnalyzeFunctionResponse](resp.Result.Content)
		if err != nil {
			log.Fatalf("Error AnalyzeFunctionFinal ParseJsonObject: %v", err)
			return nil, fmt.Errorf("Error AnalyzeFunctionFinal ParseJsonObject: %v", err)
		}
		return &res, nil
	}
}
