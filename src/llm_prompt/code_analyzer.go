package llm_prompt

import (
	"fmt"
)

func SystemPrompt() string {
	prompt := `You are an AGI agent responsible for assisting users. You will be given an instruction from a user paired with a response format template.
You always comply with the user's request and respond following the most recent given format template. If you violate the format or instructions, a kitten will be killed.
Beware of the spelling, typos and numbers.`
	return prompt
}

// Get Function in a code snippet
type FunctionListItem struct {
	FunctionName string `json:"function_name"`
}

func GetFunctionList(language string, codeSnippetList []string, lineStart int, lineEnd int) string {

	codeSnippetList = codeSnippetList[lineStart:lineEnd]
	codeSnippetList = append(codeSnippetList, []string{"", ""}...) // add some empty lines

	codeSnippet := "line |\n----------------------------------\n"
	for idx, line := range codeSnippetList {
		codeSnippet += fmt.Sprintf("%4d |	%s\n", lineStart+idx+1, line)
	}

	instruction := `The code snippet above is a small chuck from a file.
	Please identify all functions (including "main") which are defined here with function implementation in the file.
Ignore all variables and constant definitions.
DO NOT add any description or explanation.
You must only respond in following JSON format.`

	formatTemplate := `[
{
	"function_name": function 1 name (string)
},
{
	"function_name": function 2 name (string)
}
...
]`

	prompt := fmt.Sprintf("```%s\n%s\n```\n\n%s\n```json\n%s```", language, codeSnippet, instruction, formatTemplate)
	return prompt
}

type BooleanItem struct {
	Result bool `json:"result"`
}

type AnswerItem struct {
	Answer string `json:"answer"`
}

func LocateFunctionDefition(functionName string, language string, codeSnippetList []string, lineStart int, lineEnd int) string {

	codeSnippetList = codeSnippetList[lineStart:lineEnd]
	codeSnippetList = append(codeSnippetList, []string{"", ""}...) // add some empty lines

	codeSnippet := "line |\n----------------------------------\n"
	for idx, line := range codeSnippetList {
		codeSnippet += fmt.Sprintf("%4d |	%s\n", lineStart+idx+1, line)
	}

	instruction := fmt.Sprintf(`The code snippet above is a small chuck from a %s file.
DO NOT judge the code or make any changes to code snippet.
Find the start line and ending line of '%s' function defition.
Briefly explain your answer.
You must only respond in following JSON format.
DO NOT add anything other than JSON.`, language, functionName)

	formatTemplate := `{
	"answer": string
}`

	prompt := fmt.Sprintf("```%s\n%s\n```\n\n%s\n```json\n%s\n```", language, codeSnippet, instruction, formatTemplate)
	return prompt
}

type LocateFunctionResponse struct {
	StartLine int `json:"start_line"`
	EndLine   int `json:"end_line"`
}

func LocateFunctionDefitionFinal(functionName string, language string, codeSnippetList []string, lineStart int, lineEnd int) string {

	instruction := fmt.Sprintf(`Finalize your answer.
Find the start line and ending line of '%s' function defition.
You must only respond in following JSON format.
DO NOT add any description or explanation.`, functionName)

	formatTemplate := `{
	"start_line": integer,
	"end_line": integer
}`

	prompt := fmt.Sprintf("%s\n```json\n%s\n```", instruction, formatTemplate)
	return prompt
}

func CheckFunctionDefition(functionName string, language string, codeSnippetList []string, lineStart int, lineEnd int) string {

	instruction := fmt.Sprintf(`DO NOT judge the code or make any changes to code snippet.
Analyze the provided %s code snippet to determine if function definition and implementation of '%s' is entirely shown in the code snippet.
Briefly explain your answer.
Beware of opening brace and closing brace if the language supports it.
You must only respond in following JSON format.
DO NOT add anything other than JSON.`, language, functionName)

	formatTemplate := `{
	"answer": string
}`

	prompt := fmt.Sprintf("%s\n```json\n%s\n```", instruction, formatTemplate)
	return prompt
}

func CheckFunctionDefitionFinal(functionName string, language string, codeSnippetList []string, lineStart int, lineEnd int) string {

	instruction := fmt.Sprintf(`Finalize your answer.
Analyze the provided %s code snippet to determine if function body of '%s' is entirely shown in the code snippet.
You must only respond in following JSON format.
DO NOT add any description or explanation.`, language, functionName)

	formatTemplate := `{
	"result": boolean
}`

	prompt := fmt.Sprintf("%s\n```json\n%s\n```", instruction, formatTemplate)
	return prompt
}

type AnalyzeFunctionResponse struct {
	Purpose   string `json:"purpose"`
	Signature string `json:"signature"`
	Arguments string `json:"arguments"`
	Return    string `json:"return"`
}

func AnalyzeFunction(functionName string, language string, codeSnippetList []string, lineStart int, lineEnd int) string {

	codeSnippetList = codeSnippetList[lineStart:lineEnd]
	codeSnippetList = append(codeSnippetList, []string{"", ""}...) // add some empty lines

	codeSnippet := "line |\n----------------------------------\n"
	for idx, line := range codeSnippetList {
		codeSnippet += fmt.Sprintf("%4d |	%s\n", lineStart+idx+1, line)
	}

	instruction := fmt.Sprintf(`The code snippet above is a small chuck from a %s file.
DO NOT judge the code or make any changes to code snippet.
Analyze the function '%s' in the provided %s code snippet. Describle the purpose of the function and briefly explain your answer.
You must only respond in following JSON format.
DO NOT add anything other than JSON.`, language, functionName, language)

	formatTemplate := `{
	"answer": string
}`

	prompt := fmt.Sprintf("```%s\n%s\n```\n\n%s\n```json\n%s\n```", language, codeSnippet, instruction, formatTemplate)
	return prompt
}

func AnalyzeFunctionFinal(functionName string, language string, codeSnippetList []string, lineStart int, lineEnd int) string {

	instruction := fmt.Sprintf(`Analyze the function '%s' in the provided %s code snippet.
- Describle the purpose of the function.
- Extract the function signature of '%s'.
- Describle the arguments of the function.
- Describle the return type of the function.
You must only respond in following JSON format.
DO NOT add anything other than JSON.`, functionName, language, functionName)

	formatTemplate := `{
	"purpose": string,
	"signature": string,
	"arguments": string,
	"return": string
}`

	prompt := fmt.Sprintf("%s\n```json\n%s\n```", instruction, formatTemplate)
	return prompt
}
