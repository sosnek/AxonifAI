package chatgpt

import (
	"context"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
)

var client gpt3.Client
var ctx context.Context

func CreateChatGPTClient() {

	jwt := "openai-api-key"
	godotenv.Load()

	ctx = context.Background()
	client = gpt3.NewClient(jwt)
}

func AskChatGPT(message string) (output string, err error) {
	resp, err := client.CompletionWithEngine(ctx, "text-davinci-003", gpt3.CompletionRequest{
		Prompt:    []string{message},
		MaxTokens: gpt3.IntPtr(1512),
		Echo:      false,
	})

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Text, nil
}
