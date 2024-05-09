package main

import (
	"context"
	"encoding/json"
)

type ModelFile struct {
	Name string `json:"name"`

	Modelfile string `json:"modelfile"`
}

func (m *AiDockerfile) GetModelFileData(ctx context.Context) string {

	modefileContents, err := dag.CurrentModule().Source().File("Modelfile").Contents(ctx)

	if err != nil {

		panic(err)

	}

	jsonData, err := json.Marshal(ModelFile{

		Name: "aidockerfile",

		Modelfile: modefileContents,
	})

	if err != nil {

		panic(err)

	}

	return string(jsonData)
}
