package main

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func getClient() *openai.Client {
	var opts []option.RequestOption

	if config.service.baseURL != "" {
		opts = append(opts, option.WithBaseURL(config.service.baseURL))
	}
	if config.service.key != "" {
		opts = append(opts, option.WithAPIKey(config.service.key))
	}

	return openai.NewClient(opts...)
}

func getRequestParams(prompt string) openai.ChatCompletionNewParams {
	model := openai.ChatModelGPT4o
	if config.service.model != "" {
		model = config.service.model
	}

	params := openai.ChatCompletionNewParams{
		Model: openai.F(openai.ChatModel(model)),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
	}

	if config.service.seed != 0 {
		params.Seed = openai.F(config.service.seed)
	}
	if config.service.temperature != 0 {
		params.Temperature = openai.F(config.service.temperature)
	}
	if config.service.topP != 0 {
		params.TopP = openai.F(config.service.topP)
	}

	return params
}
