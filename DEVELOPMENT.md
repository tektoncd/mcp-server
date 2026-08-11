# Developing Tekton MCP Server

## Requirements

Install:

- [Go](https://go.dev/doc/install) at the version declared in [`go.mod`](./go.mod)
- [Git](https://git-scm.com/)

To run all presubmit checks locally, also install
[golangci-lint](https://golangci-lint.run/docs/welcome/install/) and
[yamllint](https://yamllint.readthedocs.io/en/stable/quickstart.html).

For cluster development, also install:

- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [ko](https://ko.build/install/)
- a Kubernetes cluster with [Tekton Pipelines](https://tekton.dev/docs/installation/pipelines/) installed

The server uses the standard Kubernetes client configuration. Local commands
therefore use the current kubeconfig context unless `KUBECONFIG` selects a
different file.

## Check out a fork

```shell
mkdir -p tektoncd
cd tektoncd
git clone https://github.com/YOUR-GITHUB-USERNAME/mcp-server.git
cd mcp-server
git remote add upstream https://github.com/tektoncd/mcp-server.git
```

Keep feature branches current with `upstream/main` and submit changes from your
fork.

## Build and test

Build the server from vendored dependencies:

```shell
mkdir -p bin
go build -mod=vendor -o bin/tekton-mcp-server ./cmd/tekton-mcp-server
```

Run the same unit-test command used by CI:

```shell
go test -v -race -timeout 5m ./...
```

CI also deploys the server and Tekton Pipelines to kind. To run its MCP smoke
test against an existing deployment:

```shell
kubectl create namespace mcp-e2e
MCP_SERVER_URL=http://127.0.0.1:8080 E2E_NAMESPACE=mcp-e2e \
  go test -v -tags=e2e -timeout=3m ./test/e2e
```

Before submitting a pull request, also run:

```shell
go build -mod=vendor ./...
gofmt -w $(find . -name '*.go' ! -path './vendor/*')
golangci-lint run --timeout=10m
yamllint -c .yamllint $(find . -path ./vendor -prune -o -type f -regex '.*y[a]ml' -print)
go install github.com/google/go-licenses@v1.0.0
go-licenses check ./...
```

CI currently uses the golangci-lint version configured in
[`.github/workflows/ci.yaml`](./.github/workflows/ci.yaml).

## Run locally

Build the binary first, then choose a transport.

For an MCP client that starts the server as a subprocess:

```shell
./bin/tekton-mcp-server -transport=stdio
```

For Streamable HTTP development:

```shell
./bin/tekton-mcp-server -transport=http -address=127.0.0.1:8080
```

The HTTP endpoint is `http://127.0.0.1:8080`. It does not provide built-in
authentication or TLS, so keep it bound to localhost during local development.

## Deploy a development build

Set `KO_DOCKER_REPO` to a registry accessible to the cluster and apply the
manifests:

```shell
export KO_DOCKER_REPO=registry.example.com/YOUR-USER/mcp-server
ko apply -R -f config/
kubectl -n tekton-mcp rollout status deployment/tekton-mcp-server
```

Connect to the in-cluster HTTP server without exposing it publicly:

```shell
kubectl -n tekton-mcp port-forward service/tekton-mcp-server 8080:8080
```

Inspect logs with:

```shell
kubectl -n tekton-mcp logs deployment/tekton-mcp-server
```

`ko delete -R -f config/` removes every included resource, including the
`tekton-mcp` namespace and anything else stored in it. Use that command only
when destructive cleanup of the namespace is intended.

The development manifests grant cluster-wide access because the server can
manage Tekton resources across namespaces. Review `config/300-rbac/` before
using them on a shared cluster.

## Update dependencies

Do not edit `vendor/` directly. After changing a module dependency, run:

```shell
go mod tidy
go mod vendor
go test ./...
git diff --check
git status --short
```

Commit `go.mod`, `go.sum`, and the resulting `vendor/` changes together. The
repository currently has no separate project-owned code-generation step.

## Debugging

- Confirm the intended cluster with `kubectl config current-context`.
- Check authorization with `kubectl auth can-i` for the resource and verb being
  exercised.
- Run the server directly with `go run ./cmd/tekton-mcp-server` to keep build
  and runtime errors in the foreground.
- For cluster runs, inspect the Deployment, Pod events, and server logs in the
  `tekton-mcp` namespace.
- Use a namespace with disposable Tekton resources when exercising create,
  update, patch, start, restart, or delete tools.

See [`CONTRIBUTING.md`](./CONTRIBUTING.md) for the project contribution process.
