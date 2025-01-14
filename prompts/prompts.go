package prompts

import (
	"bytes"
	"html/template"
)

type RequirementReq struct {
	Requirements string
	Signature    string
}

const requirementTmpl = `You are a software engineer who will work on a challenging go function.
Requirements: {{.Requirements}}
Function Signature: {{.Signature}}
You should use Test-Driven Development to approach a robust solution.
Write a good test first, be sure to consider as many common edge cases as you can.
Then provide an implementation.
This is only your first draft, so focus on breaking down the problem well and designing your overall code structure, instead of trying to arrive at a perfect solution.
Provide your solution strictly in this format, clearly marking the 2 code blocks for parsing:
== Test:
` + "```" + `go
// YOUR TEST CODE
` + "```" + `
== Implementation:
` + "```" + `go
// YOUR IMPLEMENTATION CODE
` + "```"

type IterateReq struct {
	Requirements string
	Signature    string
	PrevTest     string
	PrevImpl     string
	PrevOutput   string
	Hint         string
}

const iterateTmpl = `You are a software engineer implementing a challenging go function using Test-Driven Development.
Requirements: {{.Requirements}}
Function Signature: {{.Signature}}
You've written a unit test and a draft implementation, but some tests have failed.
== Previous Test:
` + "```" + `go
{{.PrevTest}}
` + "```" + `
== Previous Implementation:
` + "```" + `go
{{.PrevImpl}}
` + "```" + `
== Output:
` + "```" + `
{{.PrevOutput}}
` + "```" + `
Your task is to iterate towards a correct unit test and implementation. Think about why the given solution failed and how you can fix it.
Note that it may be the previous unit test that was incorrect (as may be the implementation), so look for problems in both places.
{{if .Hint}}Hint: {{.Hint}}{{end}}
Provide an updated version of your solution strictly in this format, clearly marking the 2 code blocks for parsing:
== Updated Test:
` + "```" + `go
// YOUR UPDATED TEST CODE
` + "```" + `
== Updated Implementation:
` + "```" + `go
// YOUR UPDATED IMPLEMENTATION CODE
` + "```"

var (
	requirementTemplate = template.Must(template.New("requirement").Parse(requirementTmpl))
	iterateTemplate     = template.Must(template.New("iterate").Parse(iterateTmpl))
)

func RequirementPrompt(req RequirementReq) string {
	var buf bytes.Buffer
	if err := requirementTemplate.Execute(&buf, req); err != nil {
		return "Error generating prompt: " + err.Error()
	}
	return buf.String()
}

func IteratePrompt(req IterateReq) string {
	var buf bytes.Buffer
	if err := iterateTemplate.Execute(&buf, req); err != nil {
		return "Error generating prompt: " + err.Error()
	}
	return buf.String()
}
