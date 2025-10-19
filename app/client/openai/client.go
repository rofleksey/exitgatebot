package openai

import (
	"context"
	_ "embed"
	"exitgatebot/app/config"
	"fmt"
	"strings"

	"github.com/samber/do"
	"github.com/sashabaranov/go-openai"
)

//go:embed SYSTEM_PROMPT.txt
var systemPrompt string

const resultPrefix = "Result:"

type Client struct {
	cfg    *config.Config
	client *openai.Client
}

func NewClient(di *do.Injector) (*Client, error) {
	cfg := do.MustInvoke[*config.Config](di)
	clientConfig := openai.DefaultConfig(cfg.OpenAI.Token)
	clientConfig.BaseURL = cfg.OpenAI.BaseURL
	client := openai.NewClientWithConfig(clientConfig)

	return &Client{
		cfg:    cfg,
		client: client,
	}, nil
}

func (c *Client) SummarizeComment(ctx context.Context, commentText string) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.cfg.OpenAI.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: commentText,
				},
			},
			Temperature:         0.1,
			MaxCompletionTokens: 1000,
		},
	)
	if err != nil {
		return "", fmt.Errorf("CreateChatCompletion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("empty openai response")
	}

	rawResult := strings.TrimSpace(resp.Choices[0].Message.Content)
	if !strings.HasPrefix(rawResult, resultPrefix) {
		return "", fmt.Errorf("invalid openai response: %s", rawResult)
	}

	return strings.TrimSpace(strings.TrimPrefix(rawResult, resultPrefix)), nil
}
