package prompts

import (
	"fmt"
	"strings"
)

func ParseTestAndImpl(input string) (test, impl string, err error) {
	testStartIdx := strings.Index(input, "```go")
	if testStartIdx == -1 {
		return "", "", fmt.Errorf("failed to find test block")
	}

	testEndsIdx := strings.Index(input[testStartIdx+5:], "```")
	if testEndsIdx == -1 {
		return "", "", fmt.Errorf("failed to find test block")
	}

	testCode := input[testStartIdx+5 : testStartIdx+5+testEndsIdx]

	implStartIdx := strings.Index(input[testStartIdx+5+testEndsIdx:], "```go")
	if implStartIdx == -1 {
		return "", "", fmt.Errorf("failed to find implementation block")
	}

	implEndsIdx := strings.Index(input[testStartIdx+5+testEndsIdx+implStartIdx+5:], "```")
	if implEndsIdx == -1 {
		return "", "", fmt.Errorf("failed to find implementation block")
	}

	implCode := input[testStartIdx+5+testEndsIdx+implStartIdx+5 : testStartIdx+5+testEndsIdx+implStartIdx+5+implEndsIdx]

	return testCode, implCode, nil
}
