package code_analyzer

import (
	"code_assistant/src/http_client"
	"code_assistant/src/llm_prompt"
	"fmt"
	"log"
)

func AnalyzeFunction(fa *FileAnalyzer, functionName string, language string) {

	// Get a default TextGenRequest struct
	chatReq := http_client.NewChatRequest()

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
		chatReq.Messages = append(chatReq.Messages, resp.Result)
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
	}
}
