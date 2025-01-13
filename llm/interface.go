package llm

type TextGenerator interface {
	GenerateText(prompt string) (string, error)
	GenerateTextStream(prompt string) (chan string, error)
}
