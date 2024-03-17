package prompt

import (
	"fmt"
	"strings"
)

func PromptGetFunctionList(language string, codeSnippetList []string, lineStart int, lineEnd int) string {

	codeSnippet := strings.Join(codeSnippetList[lineStart:lineEnd], "\n")

	instruction := `Please identity all functions in the Go Code above.
	Only include function with full function body and definition.
	To avoid confusion, DO NOT include function with partial body.
	Create a list in JSON.
	You must only response in following JSON format.
	DO NOT add any description.`

	formatTemplate := `[
		{
			"function_name": string
		}
		...
	]`
	prompt := fmt.Sprintf("```%s\n%s```\n\n%s\n```json\n%s```", language, codeSnippet, instruction, formatTemplate)

	return prompt
}
