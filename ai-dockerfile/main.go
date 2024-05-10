package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type AiDockerfile struct{}

func (m *AiDockerfile) GuessDockerfile(ctx context.Context, projectDir *Directory) string {

	files, err := m.GetProjectFiles(ctx, projectDir)

	prompt, err := m.WrapContentFiles(ctx, files)

	if err != nil {

		panic(err)

	}

	response, err := m.CreateAIModelAndResponse(ctx, prompt)

	if err != nil {

		panic(err)

	}

	fmt.Printf(response)

	return response

}

type AIResponse struct {
	Response string `json:"response"`
}

func (m *AiDockerfile) CreateAIModelAndResponse(ctx context.Context, prompt string) (string, error) {

	response, err := dag.Container().
		From("alpine").
		WithExec([]string{"apk", "add", "curl"}).
		WithServiceBinding("ollama", m.GetOllamaSvc(ctx)).
		WithEnvVariable("CACHE_BUSTER", time.Now().String()).
		WithExec([]string{
			"curl",
			"ollama:11434/api/create",
			"-d",
			m.GetModelFileData(ctx),
		}).
		WithExec([]string{
			"curl",
			"ollama:11434/api/generate",
			"-d",
			prompt,
		}).
		Stdout(ctx)

	if err != nil {

		return "", err

	}

	parsedResponse := AIResponse{}

	json.Unmarshal([]byte(response), parsedResponse)

	fmt.Printf(parsedResponse.Response)

	return parsedResponse.Response, nil
}

func (m *AiDockerfile) GetOllamaSvc(ctx context.Context) *Service {

	return dag.Container().
		From("ollama/ollama").
		WithExposedPort(11434).
		AsService()
}
