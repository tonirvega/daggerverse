// A generated module for Kubernetes functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/kubernetes/internal/dagger"
)

type Kubernetes struct{}

// Returns a container that echoes whatever string argument is provided
func (m *Kubernetes) Run(ctx context.Context) *dagger.Container {

	dindSvc := dag.Container().
		From("docker:dind").
		WithUser("root").
		WithEnvVariable("DOCKER_TLS_CERTDIR", "").
		WithExec([]string{"-H", "tcp://0.0.0.0:2375"}, dagger.ContainerWithExecOpts{
			UseEntrypoint:            true,
			InsecureRootCapabilities: true}).
		WithExposedPort(2375).AsService()

	endpoint, err := dindSvc.Endpoint(ctx, dagger.ServiceEndpointOpts{Scheme: "tcp"})

	if err != nil {

		panic(err)

	}

	return dag.Container().
		From("alpine").
		WithEnvVariable("DOCKER_HOST", endpoint).
		WithServiceBinding("dind", dindSvc).
		// WithFile("cgroups.sh", dag.CurrentModule().Source().File("cgroups.sh")).
		// WithFile("Dockerfile", dag.CurrentModule().Source().File("Dockerfile")).
		WithExec([]string{"apk", "add", "docker", "kubectl", "k9s", "curl"}).
		WithExec([]string{"curl", "-Lo", "./kind", "https://kind.sigs.k8s.io/dl/v0.25.0/kind-linux-amd64"}).
		WithExec([]string{"chmod", "+x", "./kind"}).
		WithExec([]string{"mv", "./kind", "/usr/local/bin/kind"}, dagger.ContainerWithExecOpts{InsecureRootCapabilities: true}).
		WithExec([]string{"kind", "create", "cluster"}, dagger.ContainerWithExecOpts{InsecureRootCapabilities: true})
}

// Returns lines that match a pattern in the files of the provided Directory
func (m *Kubernetes) GrepDir(ctx context.Context, directoryArg *dagger.Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}
