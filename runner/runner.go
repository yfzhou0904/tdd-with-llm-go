package runner

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"tdd-go/config"
	"tdd-go/llm"
	"tdd-go/prompts"
)

type TestRunner struct {
	Opts      *config.Options
	Config    *config.Config
	Generator llm.TextGenerator
	Sandbox   Sandbox
}

func (r *TestRunner) Run(ctx context.Context) error {
	prompt := r.createInitialPrompt()

	for iteration := 1; ; iteration++ {
		testCode, implCode, err := r.generateDraft(ctx, iteration, prompt)
		if err != nil {
			return fmt.Errorf("generation failed: %w", err)
		}

		passed, output, err := r.runTests(iteration, testCode, implCode)
		if err != nil {
			return fmt.Errorf("test execution failed: %w", err)
		}

		if passed {
			slog.Info("test passed", "iteration", iteration)
			return nil
		}

		proceed, hint := confirmIterate()
		if !proceed {
			return fmt.Errorf("user abort")
		} else if hint != "" {
			slog.Info("using hint", "hint", hint)
		}

		prompt = prompts.IteratePrompt(prompts.IterateReq{
			Requirements: r.Opts.Specification,
			Signature:    r.Opts.FunctionSignature,
			PrevTest:     testCode,
			PrevImpl:     implCode,
			PrevOutput:   output,
			Hint:         hint,
		})
	}
}

func (r *TestRunner) generateDraft(ctx context.Context, iteration int, prompt string) (
	testCode, implCode string, err error,
) {
	if iteration == 1 && r.Opts.Iterate {
		return r.Sandbox.ReadExistingCode()
	}

	slog.Info("generating draft", "iteration", iteration, "context_words", len(strings.Fields(prompt)))
	output, err := r.Generator.GenerateText(prompt)
	if err != nil {
		return "", "", fmt.Errorf("text generation failed: %w", err)
	}

	slog.Info("parsing llm output", "words", len(strings.Fields(output)))
	testCode, implCode, err = prompts.ParseTestAndImpl(output)
	if err != nil {
		slog.Error("failed to parse llm output", "err", err)
		fmt.Println("\n" + output + "\n")
		return "", "", err
	}

	return testCode, implCode, nil
}

func (r *TestRunner) runTests(iteration int, testCode, implCode string) (passed bool, output string, err error) {
	slog.Info("testing draft", "iteration", iteration)
	output, err = r.Sandbox.WriteAndTest(testCode, implCode, r.Opts.ShowDiff)
	if err == nil {
		return true, output, nil
	}

	slog.Warn("test failed", "err", err)
	fmt.Println("\n" + output + "\n")
	return false, output, nil
}

func (r *TestRunner) createInitialPrompt() string {
	return prompts.RequirementPrompt(prompts.RequirementReq{
		Requirements: r.Opts.Specification,
		Signature:    r.Opts.FunctionSignature,
	})
}

func confirmIterate() (proceed bool, hint string) {
	fmt.Print("Test failed, proceed to the next iteration? Perhaps give a hint (y/n/<hint>): ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userResponse := scanner.Text()

	if len(userResponse) == 1 {
		if strings.ToLower(userResponse) == "y" {
			return true, ""
		}
		return false, ""
	}
	return true, userResponse
}
