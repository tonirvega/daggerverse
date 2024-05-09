package main

import (
	"context"
	"fmt"
)

type AiDockerfile struct{}

func (m *AiDockerfile) GuessDockerfile(ctx context.Context, projectDir *Directory) *File {

	files, err := m.GetProjectFiles(ctx, projectDir)

	content, err := m.WrapContentFiles(ctx, files)

	if err != nil {

		panic(err)

	}

	fmt.Println(content)

	return dag.Container().From("alpine").File("hola")

}

func (m *AiDockerfile) CreateAIModel(ctx context.Context) (string, error) {

	return dag.Container().
		From("alpine").
		WithExec([]string{"apk", "add", "curl"}).
		WithServiceBinding("ollama", m.GetOllamaSvc(ctx)).
		WithExec([]string{
			"curl",
			"ollama:11434/api/create",
			"-d",
			m.GetModelFileData(ctx),
		}).
		Stdout(ctx)

}

func (m *AiDockerfile) GetOllamaSvc(ctx context.Context) *Service {

	return dag.Container().
		From("ollama/ollama").
		WithExposedPort(11434).
		AsService()
}
