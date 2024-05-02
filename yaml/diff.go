package main

import (
	"context"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"sigs.k8s.io/yaml"
)

func (m *Yaml) Equal(ctx context.Context, yaml1 *File, yaml2 *File) bool {

	contentsFromYaml1, errF1 := yaml1.Contents(ctx)

	if errF1 != nil {

		panic(errF1)

	}

	contentsFromYaml2, errF2 := yaml2.Contents(ctx)

	if errF2 != nil {

		panic(errF2)

	}

	jsonString1, err1 := yaml.YAMLToJSON([]byte(contentsFromYaml1))

	jsonString2, err2 := yaml.YAMLToJSON([]byte(contentsFromYaml2))

	if err1 != nil {

		panic(err1)

	}

	if err2 != nil {

		panic(err2)

	}

	return jsonpatch.Equal(jsonString1, jsonString2)
}

// Diff compares two yaml files and returns a json patch string
// The output is a string that represents the json patch
func (m *Yaml) Diff(ctx context.Context, yaml1 *File, yaml2 *File) string {

	contentsFromYaml1, errF1 := yaml1.Contents(ctx)

	if errF1 != nil {

		panic(errF1)

	}

	contentsFromYaml2, errF2 := yaml2.Contents(ctx)

	if errF2 != nil {

		panic(errF2)

	}

	jsonString1, err1 := yaml.YAMLToJSON([]byte(contentsFromYaml1))

	jsonString2, err2 := yaml.YAMLToJSON([]byte(contentsFromYaml2))

	if err1 != nil {

		panic(err1)

	}

	if err2 != nil {

		panic(err2)

	}

	patch, err := jsonpatch.CreateMergePatch(jsonString1, jsonString2)

	if err != nil {

		panic(err)

	}

	return string(patch)
}

func (m *Yaml) GetFile(ctx context.Context, path string) *File {

	return dag.CurrentModule().Source().File(path)

}
