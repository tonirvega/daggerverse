package main

import (
	"context"
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

	return content, nil

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
