package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"tdd-go/llm"
	"tdd-go/prompts"

	"github.com/BurntSushi/toml"
)

func main() {
	opts, config := parseFlags()

	var gen llm.TextGenerator

	if opts.Provider == "oai" {
		gen = llm.OAI{
			BaseURL: config.OpenAI.BaseURL,
			Key:     config.OpenAI.Key,
		}
	} else if opts.Provider == "anthropic" {
		gen = llm.Claude{
			BaseURL: config.Anthropic.BaseURL,
			Key:     config.Anthropic.Key,
		}
	} else {
		log.Fatalln("unknown provider")
	}

	prompt := prompts.RequirementPrompt(prompts.RequirementReq{
		Requirements: opts.Specification,
		Signature:    opts.FunctionSignature,
	})

	iteration := 0
	for {
		iteration++

		var (
			testCode, implCode, output string
		)

		slog.Info(fmt.Sprintf("generating draft #%d", iteration))
		output, err := gen.GenerateText(prompt)
		if err != nil {
			log.Fatalln(err.Error())
		}

		slog.Info("parsing llm output", "words", len(strings.Fields(output)))
		testCode, implCode, err = prompts.ParseTestAndImpl(output)
		if err != nil {
			slog.Error("failed to parse llm output", "err", err.Error())
			fmt.Println("\n" + output + "\n")
			return
		}

		slog.Info(fmt.Sprintf("testing draft #%d", iteration))
		output, err = testImplementation(testCode, implCode)
		if err == nil {
			fmt.Println("test passed")
			return
		}

		slog.Error("test failed", "err", err.Error())
		fmt.Println("\n" + output + "\n")

		proceed, hint := false, ""
		if proceed, hint = confirmIterate(); !proceed {
			fmt.Println("aborting")
			return
		}

		prompt = prompts.IteratePrompt(prompts.IterateReq{
			Requirements: opts.Specification,
			Signature:    opts.FunctionSignature,
			PrevTest:     testCode,
			PrevImpl:     implCode,
			PrevOutput:   output,
			Hint:         hint,
		})
	}
}

func testImplementation(testCode, implCode string) (output string, err error) {
	// write code to go files under ./sandbox/
	os.MkdirAll("./sandbox", os.ModePerm)

	testFile, err := os.Create("./sandbox/main_test.go.new")
	if err != nil {
		return "", err
	}
	defer testFile.Close()
	_, err = testFile.WriteString(testCode)
	if err != nil {
		return "", err
	}

	implFile, err := os.Create("./sandbox/main.go.new")
	if err != nil {
		return "", err
	}
	_, err = implFile.WriteString(implCode)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("diff", "-y", "./sandbox/main.go", "./sandbox/main.go.new")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = exec.Command("diff", "-y", "./sandbox/main_test.go", "./sandbox/main_test.go.new")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	if err = exec.Command("mv", "./sandbox/main.go.new", "./sandbox/main.go").Run(); err != nil {
		log.Fatalln("failed to mv", err)
	}
	if err = exec.Command("mv", "./sandbox/main_test.go.new", "./sandbox/main_test.go").Run(); err != nil {
		log.Fatalln("failed to mv", err)
	}

	os.Chdir("./sandbox")
	defer os.Chdir("..")
	outputBytes, err := exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		return string(outputBytes), err
	}
	exec.Command("gofmt", "-w", ".").Run()
	exec.Command("goimports", "-w", ".").Run()
	outputBytes, err = exec.Command("go", "test", ".", "-v").CombinedOutput()
	if err != nil {
		return string(outputBytes), err
	}
	return string(outputBytes), nil
}

func confirmIterate() (proceed bool, hint string) {
	fmt.Print("\nTest failed, proceed to the next iteration? Perhaps give a hint (y/n/<hint>): ")

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

func parseFlags() (Options, Config) {
	opts := Options{}
	flag.StringVar(&opts.ConfigPath, "config", "./config.toml", "configuration file")
	flag.StringVar(&opts.Specification, "spec", "", "detailed spec & requirements for the function to be implemented")
	flag.StringVar(&opts.FunctionSignature, "sig", "", "function signature")
	flag.StringVar(&opts.Provider, "provider", "anthropic", "oai / anthropic")
	flag.StringVar(&opts.Hint, "hint", "", "hint for llm")
	flag.Parse()

	if opts.Specification == "" || opts.FunctionSignature == "" {
		log.Fatalln("requirements and function signature are required")
	}

	var config Config
	if _, err := os.Stat(opts.ConfigPath); err == nil {
		_, err = toml.DecodeFile(opts.ConfigPath, &config)
		if err != nil {
			log.Fatalln(err)
		}

	} else {
		log.Fatalln("can't stat config file")
	}

	return opts, config
}

type Options struct {
	ConfigPath        string
	Specification     string
	FunctionSignature string
	Hint              string
	Provider          string
}

type Config struct {
	OpenAI struct {
		BaseURL string `toml:"base_url"`
		Key     string `toml:"key"`
	} `toml:"openai"`
	Anthropic struct {
		BaseURL string `toml:"base_url"`
		Key     string `toml:"key"`
	} `toml:"anthropic"`
}
