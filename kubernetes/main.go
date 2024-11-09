package main

import (
	"context"
	"dagger/kubernetes/internal/dagger"
	"time"
)

type Kubernetes struct{}

// dagger call run --docker-sock=/var/run/docker.sock --kind-svc=tcp://127.0.0.1:3000 terminal stdout
func (m *Kubernetes) Run(ctx context.Context, dockerSock *dagger.Socket, kindSvc *dagger.Service) *dagger.Container {

	return dag.Container().
		From("alpine").
		WithUnixSocket("/var/run/docker.sock", dockerSock).
		WithFile("kind.yaml", dag.CurrentModule().Source().File("kind.yaml")).
		WithExec([]string{"apk", "add", "docker", "kubectl", "k9s", "curl"}).
		WithExec([]string{"curl", "-Lo", "./kind", "https://kind.sigs.k8s.io/dl/v0.25.0/kind-linux-amd64"}).
		WithExec([]string{"chmod", "+x", "./kind"}).
		WithExec([]string{"mv", "./kind", "/usr/local/bin/kind"}).
		WithEnvVariable("BUST", time.Now().String()).
		WithExec([]string{"kind", "delete", "cluster"}).
		WithExec([]string{
			"kind", "create", "cluster",
			"--config", "kind.yaml", "--wait", "30s",
		}, dagger.ContainerWithExecOpts{
			InsecureRootCapabilities: true,
		}).
		WithServiceBinding("localhost", kindSvc).
		WithExec([]string{"kubectl", "get", "nodes", "--server=https://localhost:3000"})
}
