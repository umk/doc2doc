package main

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"strings"
)

//go:embed messages/create.tmpl
var createMessage string

var createMessageTmpl = template.Must(template.New("create").Parse(createMessage))

//go:embed messages/update.tmpl
var updateMessage string

var updateMessageTmpl = template.Must(template.New("update").Parse(updateMessage))

func generate(
	ctx context.Context,
	c *configService,
	prompt string,
	previousIn, previousOut *string,
	in string,
	outputPath string,
) (string, error) {
	var sb strings.Builder
	if previousIn == nil || previousOut == nil {
		if err := createMessageTmpl.Execute(&sb, map[string]any{
			"Prompt":      prompt,
			"In":          in,
			"OutputPath":  outputPath,
			"PreviousOut": resolvePtrOrDefault(previousOut),
		}); err != nil {
			return "", fmt.Errorf("failed to generate a message: %w", err)
		}
	} else {
		if err := updateMessageTmpl.Execute(&sb, map[string]any{
			"Prompt":      prompt,
			"In":          in,
			"OutputPath":  outputPath,
			"PreviousIn":  *previousIn,
			"PreviousOut": *previousOut,
		}); err != nil {
			return "", fmt.Errorf("failed to generate a message: %w", err)
		}
	}

	client := getClient(c)

	params := getRequestParams(c, sb.String())

	r, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", fmt.Errorf("service request failed: %w", err)
	}

	choice := r.Choices[0]

	if choice.Message.Refusal != "" {
		return "", fmt.Errorf("refused: %s", choice.Message.Refusal)
	} else {
		return choice.Message.Content, nil
	}
}
