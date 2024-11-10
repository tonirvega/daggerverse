package main

import (
	"context"
	"dagger/kubernetes/internal/dagger"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Kubernetes struct {
	DockerSocket *dagger.Socket
	KindSvc      *dagger.Service
	KindPort     int
	Container    *dagger.Container
}

func New(
	ctx context.Context,
	// Docker socket path. E.g. /var/run/docker.sock
	// How to use it:
	// dagger call --docker-sock=/var/run/docker.sock kind --kind-svc=tcp://127.0.0.1:3000
	// +required
	dockerSocket *dagger.Socket,

	// It should be the tcp://127.0.0.1 followed by any port. E.g. tcp://127.0.0.1:3000
	// Before launch this function, make sure that you have configured in your /etc/hosts file
	// an entry for localhost 127.0.0.1 . Otherwise, the alpine container will not be able to connect to the kind cluster.
	// +required
	kindSvc *dagger.Service,

) *Kubernetes {

	ep, err := kindSvc.Endpoint(ctx)

	if err != nil {

		panic(err)

	}

	port, err := strconv.Atoi(strings.Split(ep, ":")[1])

	if err != nil {

		panic(err)

	}

	if port < 1024 || port > 65535 {

		panic(fmt.Sprintf("Invalid port number: %d, it should be between 1024 and 65535", port))

	}

	container := dag.Container().
		From("alpine").
		WithUnixSocket("/var/run/docker.sock", dockerSocket).
		WithFile("kind.yaml", dag.CurrentModule().Source().File("kind.yaml")).
		WithExec([]string{"apk", "add", "docker", "kubectl", "k9s", "curl"}).
		WithExec([]string{"curl", "-Lo", "./kind", "https://kind.sigs.k8s.io/dl/v0.25.0/kind-linux-amd64"}).
		WithExec([]string{"chmod", "+x", "./kind"}).
		WithExec([]string{"mv", "./kind", "/usr/local/bin/kind"}).
		WithEnvVariable("BUST", time.Now().String()).
		WithExec([]string{"kind", "delete", "cluster"}).
		WithExec([]string{
			"kind", "create", "cluster",
			"--config", "kind.yaml",
			"--wait", "1m",
		}).
		WithServiceBinding("localhost", kindSvc).
		WithExec([]string{
			"kubectl", "config",
			"set-cluster", "kind-kind", fmt.Sprintf("--server=https://localhost:%d", port)},
		)

	return &Kubernetes{
		DockerSocket: dockerSocket,
		KindSvc:      kindSvc,
		KindPort:     port,
		Container:    container,
	}
}

func (m *Kubernetes) LoadContainerOnKind(

	ctx context.Context,

	container *dagger.Container,

	tag string,

) *dagger.Container {

	containerName := fmt.Sprintf("%s.tar", tag)

	tarball := container.
		// This is the image name that will be loaded in the kind cluster
		WithAnnotation(
			"org.opencontainers.image.ref.name",
			fmt.Sprintf("%s:latest", tag),
		).

		// Kind requires the docker.io/library prefix, otherwise it will load the image
		// This a fake image name in docker.io, it is not a real image.
		// You should user imagePullPolicy: Never in your Kubernetes manifests.
		WithAnnotation(
			"io.containerd.image.name",
			fmt.Sprintf("docker.io/library/%s:latest", tag),
		).
		AsTarball()

	return m.Container.
		WithFile(containerName, tarball).
		WithEnvVariable("BUST", time.Now().String()).
		WithExec([]string{"kind", "load", "image-archive", fmt.Sprintf("%s.tar", tag)})

}

func (m *Kubernetes) K9s(

	ctx context.Context,

) *dagger.Container {

	return m.Container.WithExec([]string{"k9s"}).Terminal()

}

func (m *Kubernetes) Inspect(

	ctx context.Context,

) *dagger.Container {

	return m.Container.Terminal()

}

func (m *Kubernetes) Run(

	ctx context.Context,

) *dagger.Container {

	return m.LoadContainerOnKind(ctx, dag.Container().From("alpine"), "alpinesito")

}
