package main

import (
	"context"
	"encoding/json"
	"strings"
)

func (m *AiDockerfile) GetProjectFiles(ctx context.Context, projectDir *Directory) ([]*File, error) {

	files := []*File{}

	entries, err := projectDir.Entries(ctx)

	if err != nil {

		return nil, err
	}

	for _, entry := range entries {

		if m.IsDir(ctx, projectDir, entry) {
			continue
		}

		files = append(files, projectDir.File(entry))

	}

	return files, nil

}

type Generate struct {
	Model string `json:"model"`

	Prompt string `json:"prompt"`

	Stream bool `json:"stream"`
}

func (m *AiDockerfile) WrapContentFiles(ctx context.Context, files []*File) (string, error) {

	content := ""

	for _, file := range files {

		fileContent, err := file.Contents(ctx)

		if err != nil {
			return "", err
		}

		fileName, err := file.Name(ctx)

		if err != nil {

			return "", err

		}

		content += "# " + fileName + "\n"
		content += fileContent + "\n"

	}

	generate, err := json.Marshal(Generate{
		Model:  "aidockerfile",
		Prompt: content,
		Stream: false,
	})

	if err != nil {

		return "", err

	}

	return string(generate), nil

}

func (m *AiDockerfile) IsDir(ctx context.Context, dir *Directory, path string) bool {

	_, err := dir.Directory(path).Sync(ctx)

	if err != nil {

		if strings.Contains(err.Error(), "not a directory") {
			return false
		}

		panic(err)

	}

	return true
}
