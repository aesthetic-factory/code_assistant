package llm_prompt

import (
	"fmt"
)

func SystemPrompt() string {
	prompt := `You are an AGI agent responsible for assisting users. You will be given an instruction from a user paired with a response format template.
You always comply with the user's request and respond following the format template. If you violate the format or instructions, a kitten will be killed.
Beware of the spelling,  typos, and numbers.`
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
Please identity all function (inlcuding "main") which are defined here.
Create a list in JSON.
DO NOT add any description or explanation.
You must only response in following JSON format.`

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
	Result string `json:"result"`
}

// Get Function in a code snippet
type AnswerItem struct {
	Answer string `json:"answer"`
}

func LocateFunctionDefition(functionName string, language string, codeSnippetList []string, lineStart int, lineEnd int) string {

	codeSnippetList = codeSnippetList[lineStart:lineEnd]
	// codeSnippetList = append(codeSnippetList, []string{"", ""}...) // add some empty lines

	codeSnippet := "line |\n----------------------------------\n"
	for idx, line := range codeSnippetList {
		codeSnippet += fmt.Sprintf("%4d |	%s\n", lineStart+idx+1, line)
	}

	instruction := fmt.Sprintf(`The code snippet above is a small chuck from a %s file.
DO NOT judge the code or make any changes to code snippet.
Find the start line and ending line of '%s' function defition.
Briefly explain your answer.
You must only response in following JSON format. DO NOT add anything other than JSON.`, language, functionName)

	formatTemplate := `{
	"answer": string
}`

	prompt := fmt.Sprintf("```%s\n%s\n```\n\n%s\n```json\n%s\n```", language, codeSnippet, instruction, formatTemplate)
	return prompt
}

func CheckFunctionDefition(functionName string, language string, codeSnippetList []string, lineStart int, lineEnd int) string {

	instruction := fmt.Sprintf(`DO NOT judge the code or make any changes to code snippet.
Analyze the provided %s code snippet to determine if function body of '%s' is entirely shown in the code snippet.
Briefly explain your answer.
You must only response in following JSON format.`, language, functionName)

	formatTemplate := `{
	"answer": string
}`

	prompt := fmt.Sprintf("%s\n```json\n%s\n```", instruction, formatTemplate)
	return prompt
}

func CheckFunctionDefitionFinal(functionName string, language string, codeSnippetList []string, lineStart int, lineEnd int) string {

	instruction := fmt.Sprintf(`Finalize your answer.
Analyze the provided %s code snippet to determine if function body of '%s' is entirely shown in the code snippet.
You must only response in following JSON format.`, language, functionName)

	formatTemplate := `{
	"result": boolean
}`

	prompt := fmt.Sprintf("%s\n```json\n%s\n```", instruction, formatTemplate)
	return prompt
}
