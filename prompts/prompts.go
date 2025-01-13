package prompts

import "fmt"

func RequirementPrompt(req, sig string) string {
	return fmt.Sprintf(`You are a software engineer who will work on a challenging go function.
Requirements: %s
Function Signature: %s
You should use Test-Driven Development to approach a robust solution.
Write a good test first, be sure to consider as many common edge cases as you can.
Then provide an implementation.
This is only your first draft, so focus on breaking down the problem well and designing your overall code structure, instead of trying to arrive at a perfect solution.
Provide your solution strictly in this format, clearly marking the 2 code blocks for parsing:
== Test:
`+"```"+`go
// YOUR TEST CODE
`+"```"+`
== Implementation:
`+"```"+`go
// YOUR IMPLEMENTATION CODE
`+"```", req, sig)
}

func IteratePrompt(req, sig, test, impl, output string) string {
	return fmt.Sprintf(`You are a software engineer implementing a challenging go function using Test-Driven Development.
Requirements: %s
Function Signature: %s
You've written a unit test and a draft implementation, but some tests have failed.
== Previous Test:
`+"```"+`go
%s
`+"```"+`
== Previous Implementation:
`+"```"+`go
%s
`+"```"+`
== Output:
`+"```"+`
%s
`+"```"+`
Your task is to iterate towards a correct unit test and implementation. Think about why the given solution failed and how you can fix it.
Note that it may be the previous unit test that was incorrect (as may be the implementation), so look for problems in both.
Provide an updated version of your solution strictly in this format, clearly marking the 2 code blocks for parsing:
== Updated Test:
`+"```"+`go
// YOUR UPDATED TEST CODE
`+"```"+`
== Updated Implementation:
`+"```"+`go
// YOUR UPDATED IMPLEMENTATION CODE
`+"```", req, sig, test, impl, output)
}
