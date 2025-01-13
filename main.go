package main

import (
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

	oai := llm.OAI{
		BaseURL: config.OpenAI.BaseURL,
		Key:     config.OpenAI.Key,
	}

	iteration := 0
	for {
		var (
			testCode, implCode, output string
		)

		iteration++

		slog.Info(fmt.Sprintf("generating draft #%d", iteration))

		prompt := ""
		if iteration == 1 {
			prompt = prompts.RequirementPrompt(opts.Requirements, opts.FunctionSignature)
		} else {
			prompt = prompts.IteratePrompt(opts.Requirements, opts.FunctionSignature, testCode, implCode, output)
		}
		output, err := oai.GenerateText(prompt)
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
		if !confirmIterate() {
			fmt.Println("aborting")
			return
		}
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

func confirmIterate() bool {
	fmt.Print("\nDo you want to proceed with another iteration? (y/N): ")
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
}

func parseFlags() (Options, Config) {
	opts := Options{}
	flag.StringVar(&opts.Requirements, "r", "", "requirements for the function")
	flag.StringVar(&opts.FunctionSignature, "s", "", "function signature")
	flag.StringVar(&opts.ConfigPath, "c", "./config.toml", "configuration file")
	flag.Parse()

	if opts.Requirements == "" || opts.FunctionSignature == "" {
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
	Requirements      string
	FunctionSignature string
}

type Config struct {
	OpenAI struct {
		BaseURL string `toml:"base_url"`
		Key     string `toml:"key"`
	} `toml:"openai"`
}
