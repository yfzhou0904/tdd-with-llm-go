package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

func Parse() (opts Options, config Config, err error) {
	opts = Options{}
	// required:
	flag.StringVar(&opts.Specification, "spec", "", "specification for the function to be implemented")
	flag.StringVar(&opts.FunctionSignature, "sig", "", "function signature to be implemented")
	// optional:
	flag.StringVar(&opts.ConfigPath, "config", "./config.toml", "configuration file")
	flag.StringVar(&opts.Provider, "provider", "anthropic", "oai | anthropic")
	flag.StringVar(&opts.ImplCodePath, "impl-file", "./sandbox/main.go", "AI-generated implementation code path")
	flag.StringVar(&opts.TestCodePath, "test-file", "./sandbox/main_test.go", "AI-generated test code path")
	flag.BoolVar(&opts.ShowDiff, "diff", false, "show code diff for each AI edit")
	flag.BoolVar(&opts.Iterate, "iterate", false, "iterate on existing code")
	flag.Parse()

	if opts.Specification == "" || opts.FunctionSignature == "" {
		err = fmt.Errorf("--spec and --sig are required")
		return
	}
	if _, err = os.Stat(opts.ConfigPath); err != nil {
		err = fmt.Errorf("can't stat config file: %w", err)
		return
	}
	_, err = toml.DecodeFile(opts.ConfigPath, &config)
	if err != nil {
		err = fmt.Errorf("can't decode config file: %w", err)
		return
	}
	return opts, config, nil
}

type Options struct {
	ConfigPath        string
	ImplCodePath      string
	TestCodePath      string
	Specification     string
	FunctionSignature string
	Iterate           bool
	Provider          string
	ShowDiff          bool
}

type Config struct {
	OpenAI    OpenAIConfig    `toml:"openai"`
	Anthropic AnthropicConfig `toml:"anthropic"`
}

type OpenAIConfig struct {
	BaseURL string `toml:"base_url"`
	Key     string `toml:"key"`
}

type AnthropicConfig struct {
	BaseURL string `toml:"base_url"`
	Key     string `toml:"key"`
}
