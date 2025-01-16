# tdd-with-llm-go

Blog Post:
[https://blog.yfzhou.fyi/posts/tdd-llm/](https://blog.yfzhou.fyi/posts/tdd-llm/)

Usage:
1. Provide your own OpenAI / Anthropic API keys in config.toml (refer to config.example.toml).
2. Run main.go
```sh
% go run main.go -h
Usage of /var/folders/8q/tbphznl55hv1651h4d2td2d80000gn/T/go-build1989261082/b001/exe/main:
  -config string
    	configuration file (default "./config.toml")
  -diff
    	show code diff for each AI edit
  -impl-file string
    	AI-generated implementation code path (default "./sandbox/main.go")
  -iterate
    	iterate on existing code
  -provider string
    	oai | anthropic (default "anthropic")
  -sig string
    	function signature to be implemented
  -spec string
    	specification for the function to be implemented
  -test-file string
    	AI-generated test code path (default "./sandbox/main_test.go")

% go run main.go \
--spec 'develop a function to take in a large text, recognize and parse any and all ipv4 and ipv6 addresses and CIDRs contained within it (these may be surrounded by random words or symbols like commas), then return them as a list' \
--sig 'func ParseCidrs(input string) ([]*net.IPNet, error)'
```

Requires:
 - goimports: `go install golang.org/x/tools/cmd/goimports@latest`

TODOs:
 - support shared context file (e.g. for custom types, other setup). Put in sandbox/shared.go by default.
 - support ollama for local models.
 - support complex project structure.
