package main

import (
	"context"
	"log/slog"
	"os"
	"tdd-go/config"
	"tdd-go/provider"

	"tdd-go/runner"
)

func main() {
	opts, config, err := config.Parse()
	if err != nil {
		slog.Error("failed to parse options", "err", err.Error())
		os.Exit(1)
	}

	gen, err := provider.NewTextGenerator(opts.Provider, config)
	if err != nil {
		slog.Error("failed to initialize text generator", "err", err.Error())
		os.Exit(1)
	}

	runner := runner.TestRunner{
		Config:    &config,
		Opts:      &opts,
		Generator: gen,
		Sandbox:   runner.LocalSandbox{ImplCodePath: opts.ImplCodePath, TestCodePath: opts.TestCodePath},
	}

	if err := runner.Run(context.Background()); err != nil {
		slog.Error("runner failed", "err", err)
		os.Exit(1)
	}

}
