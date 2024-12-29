package main

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

func getClient(c *configService) *openai.Client {
	var opts []option.RequestOption

	if c.BaseURL != nil {
		opts = append(opts, option.WithBaseURL(*c.BaseURL))
	}
	if c.Key != nil {
		opts = append(opts, option.WithAPIKey(*c.Key))
	}

	return openai.NewClient(opts...)
}

func getRequestParams(c *configService, prompt string) openai.CompletionNewParams {
	s := shared.UnionString(prompt)

	model := openai.ChatModelGPT4o
	if c.Model != nil {
		model = *c.Model
	}

	params := openai.CompletionNewParams{
		Model:  openai.F(openai.CompletionNewParamsModel(model)),
		Prompt: openai.F(openai.CompletionNewParamsPromptUnion(s)),
	}

	if c.Seed != nil {
		params.Seed = openai.F(*c.Seed)
	}
	if c.Temperature != nil {
		params.Temperature = openai.F(*c.Temperature)
	}
	if c.TopP != nil {
		params.TopP = openai.F(*c.TopP)
	}

	return params
}
