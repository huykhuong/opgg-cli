package gemini

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/huykhuong/lol/internal/opgg"
	"google.golang.org/genai"
)

type Client struct {
	Ctx    context.Context
	Client *genai.Client
}

func NewClient(ctx context.Context, apiKey string) (*Client, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return nil, fmt.Errorf("an error occurred while creating Gemini client: %w", err)
	}

	return &Client{
		Ctx:    ctx,
		Client: client,
	}, nil
}

func buildTools(mcpTools []opgg.Tool) []*genai.Tool {
	var declarations []*genai.FunctionDeclaration

	for _, t := range mcpTools {
		decl := &genai.FunctionDeclaration{
			Name:        t.Name,
			Description: t.Description,
		}

		if len(t.InputSchema) > 0 {
			var schema any
			if err := json.Unmarshal(t.InputSchema, &schema); err == nil {
				decl.ParametersJsonSchema = schema
			}
		}

		declarations = append(declarations, decl)
	}

	return []*genai.Tool{{FunctionDeclarations: declarations}}
}

func (c *Client) Execute(prompt string, opggClient *opgg.Client) (string, error) {
	mcpTools, err := opggClient.ListTools(c.Ctx)
	if err != nil {
		return "", fmt.Errorf("failed to discover tools: %w", err)
	}

	chat, err := c.Client.Chats.Create(c.Ctx, "gemini-2.5-flash", &genai.GenerateContentConfig{
		Tools: buildTools(mcpTools),
	}, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create chat session: %w", err)
	}

	resp, err := chat.SendMessage(c.Ctx, genai.Part{Text: prompt})
	if err != nil {
		return "", err
	}

	for {
		part := resp.Candidates[0].Content.Parts[0]
		if part.FunctionCall == nil {
			break
		}

		funcCall := part.FunctionCall
		fmt.Printf("🤖 AI wants to call: %s\n", funcCall.Name)

		toolResult, err := opggClient.CallTool(c.Ctx, funcCall.Name, funcCall.Args)
		if err != nil {
			return "", fmt.Errorf("tool execution failed: %w", err)
		}

		resp, err = chat.SendMessage(c.Ctx, genai.Part{
			FunctionResponse: &genai.FunctionResponse{
				Name:     funcCall.Name,
				Response: toolResult,
			},
		})
		if err != nil {
			return "", err
		}
	}

	return resp.Candidates[0].Content.Parts[0].Text, nil
}
