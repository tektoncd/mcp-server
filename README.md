# Tekton Model Context Protocol server
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ftektoncd%2Fmcp-server.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Ftektoncd%2Fmcp-server?ref=badge_shield)


Tekton MCP Server exposes Tekton resources to
[Model Context Protocol (MCP)](https://modelcontextprotocol.io) clients. It
currently focuses on [`tektoncd/pipeline`](https://github.com/tektoncd/pipeline)
and supports both stdio and Streamable HTTP transports.

> [!IMPORTANT]
> This project is under active development. Test it in a non-production cluster
> and review the access granted to both the server and the connected MCP client.

## Prerequisites

The server needs:

- access to a Kubernetes cluster with
  [Tekton Pipelines](https://tekton.dev/docs/installation/pipelines/) installed
- a kubeconfig for local use, or an in-cluster service account when deployed to
  Kubernetes
- [Go](https://go.dev/doc/install) at the version in [`go.mod`](./go.mod) when
  building from source

[`kubectl`](https://kubernetes.io/docs/tasks/tools/) and
[`ko`](https://ko.build/) are also required to deploy the manifests from a
source checkout.

## Build from source

```shell
git clone https://github.com/tektoncd/mcp-server.git
cd mcp-server
mkdir -p bin
go build -mod=vendor -o bin/tekton-mcp-server ./cmd/tekton-mcp-server
```

The server uses the standard Kubernetes client configuration. Confirm that
`kubectl config current-context` points to the intended cluster before running
it.

## Run locally

### stdio

Use stdio when an MCP client starts the server as a subprocess:

```shell
./bin/tekton-mcp-server -transport=stdio
```

A typical MCP client entry looks like this; adapt the surrounding configuration
to the client being used:

```json
{
  "mcpServers": {
    "tekton": {
      "command": "/absolute/path/to/bin/tekton-mcp-server",
      "args": ["-transport=stdio"]
    }
  }
}
```

### Streamable HTTP

Bind to localhost for local development:

```shell
./bin/tekton-mcp-server -transport=http -address=127.0.0.1:8080
```

Connect the MCP client to `http://127.0.0.1:8080`.

## Deploy to Kubernetes

Set `KO_DOCKER_REPO` to a registry accessible to the cluster, then build the
image and apply the manifests:

```shell
export KO_DOCKER_REPO=registry.example.com/YOUR-USER/mcp-server
ko apply -R -f config/
kubectl -n tekton-mcp rollout status deployment/tekton-mcp-server
```

For local access to the in-cluster service:

```shell
kubectl -n tekton-mcp port-forward service/tekton-mcp-server 8080:8080
```

`ko delete -R -f config/` removes every included resource, including the
`tekton-mcp` namespace and anything else stored in that namespace. Use it only
when that destructive cleanup is intended.

## Compatibility

The server uses Tekton `v1` resources and is built and tested against the
Kubernetes, Tekton Pipelines, and MCP Go SDK versions recorded in
[`go.mod`](./go.mod). Until a broader compatibility matrix is published, test
the server with the exact cluster versions on which it will run.

## Security

The server can read and modify Tekton resources using the permissions of its
kubeconfig or service account. Several tools create, patch, start, restart, and
delete resources, so use a dedicated least-privilege identity and a disposable
namespace while evaluating it.

The HTTP transport does not provide authentication or TLS. Do not expose it to
an untrusted network without an authenticating, TLS-terminating proxy. The
included development manifests grant cluster-wide permissions and should be
reviewed before use on a shared cluster.

Report vulnerabilities privately through the
[project security policy](https://github.com/tektoncd/mcp-server/security/policy).
Do not open a public issue for a suspected vulnerability.

## Development and contributing

See [`DEVELOPMENT.md`](./DEVELOPMENT.md) for build, test, dependency update,
cluster deployment, and debugging instructions. Contributions follow the
process in [`CONTRIBUTING.md`](./CONTRIBUTING.md).

The project was proposed and accepted in
[`tektoncd/community#1194`](https://github.com/tektoncd/community/issues/1194).

## Tools

### List Operations

#### `list_pipelines` – List Pipelines in the Cluster with Filtering Options
- `namespace`: Namespace to list Pipelines from (string, required)
- `prefix`: Name prefix to filter Pipelines (string, optional)
- `label-selector`: Label selector to filter Pipelines (string, optional)

#### `list_pipelineruns` – List PipelineRuns in the Cluster with Filtering Options
- `namespace`: Namespace to list PipelineRuns from (string, required)
- `prefix`: Name prefix to filter PipelineRuns (string, optional)
- `label-selector`: Label selector to filter PipelineRuns (string, optional)

#### `list_tasks` – List Tasks in the Cluster with Filtering Options
- `namespace`: Namespace to list Tasks from (string, required)
- `prefix`: Name prefix to filter Tasks (string, optional)
- `label-selector`: Label selector to filter Tasks (string, optional)

#### `list_taskruns` – List TaskRuns in the Cluster with Filtering Options
- `namespace`: Namespace to list TaskRuns from (string, required)
- `prefix`: Name prefix to filter TaskRuns (string, optional)
- `label-selector`: Label selector to filter TaskRuns (string, optional)

#### `list_stepactions` – List Step Actions in the Cluster with Filtering Options
- `namespace`: Namespace to list Step Actions from (string, required)
- `prefix`: Name prefix to filter Step Actions (string, optional)
- `label-selector`: Label selector to filter Step Actions (string, optional)

### Create Operations

#### `create_pipeline` – Create a new Pipeline from YAML definition
- `namespace`: Namespace where the Pipeline will be created (string, optional, default: "default")
- `yaml`: YAML definition of the Pipeline (string, required)

#### `create_task` – Create a new Task from YAML definition
- `namespace`: Namespace where the Task will be created (string, optional, default: "default")
- `yaml`: YAML definition of the Task (string, required)

#### `create_pipelinerun` – Create a new PipelineRun from YAML definition or generate from Pipeline
- `namespace`: Namespace where the PipelineRun will be created (string, optional, default: "default")
- `yaml`: YAML definition of the PipelineRun (string, optional)
- `generateName`: Generate name prefix for the PipelineRun (string, optional)

#### `create_taskrun` – Create a new TaskRun from YAML definition
- `namespace`: Namespace where the TaskRun will be created (string, optional, default: "default")
- `yaml`: YAML definition of the TaskRun (string, optional)
- `generateName`: Generate name prefix for the TaskRun (string, optional)

### Get Operations

#### `get_pipeline` – Get a specific Pipeline by name
- `name`: Name of the Pipeline to get (string, required)
- `namespace`: Namespace of the Pipeline (string, optional, default: "default")
- `output`: Output format - json or yaml (string, optional, default: "yaml")

#### `get_task` – Get a specific Task by name
- `name`: Name of the Task to get (string, required)
- `namespace`: Namespace of the Task (string, optional, default: "default")
- `output`: Output format - json or yaml (string, optional, default: "yaml")

#### `get_pipelinerun` – Get a specific PipelineRun by name
- `name`: Name of the PipelineRun to get (string, required)
- `namespace`: Namespace of the PipelineRun (string, optional, default: "default")
- `output`: Output format - json or yaml (string, optional, default: "yaml")

#### `get_taskrun` – Get a specific TaskRun by name
- `name`: Name of the TaskRun to get (string, required)
- `namespace`: Namespace of the TaskRun (string, optional, default: "default")
- `output`: Output format - json or yaml (string, optional, default: "yaml")

#### `get_taskrun_logs` - Get the logs for a given TaskRun
- `name`: Name or reference of the TaskRun to get logs from (string, required)
- `namespace`: Namespace where the TaskRun is located (string, optional, default: "default")

### Update Operations

#### `update_pipeline` – Update an existing Pipeline
- `name`: Name of the Pipeline to update (string, required)
- `namespace`: Namespace of the Pipeline (string, optional, default: "default")
- `yaml`: Updated YAML definition of the Pipeline (string, required)

#### `update_task` – Update an existing Task
- `name`: Name of the Task to update (string, required)
- `namespace`: Namespace of the Task (string, optional, default: "default")
- `yaml`: Updated YAML definition of the Task (string, required)

#### `patch_pipeline` – Apply a JSON patch to an existing Pipeline
- `name`: Name of the Pipeline to patch (string, required)
- `namespace`: Namespace of the Pipeline (string, optional, default: "default")
- `patch`: JSON patch to apply to the Pipeline (string, required)

### Delete Operations

#### `delete_pipeline` – Delete a Pipeline
- `name`: Name of the Pipeline to delete (string, required)
- `namespace`: Namespace of the Pipeline (string, optional, default: "default")

#### `delete_task` – Delete a Task
- `name`: Name of the Task to delete (string, required)
- `namespace`: Namespace of the Task (string, optional, default: "default")

#### `delete_pipelinerun` – Delete a PipelineRun
- `name`: Name of the PipelineRun to delete (string, required)
- `namespace`: Namespace of the PipelineRun (string, optional, default: "default")

#### `delete_taskrun` – Delete a TaskRun
- `name`: Name of the TaskRun to delete (string, required)
- `namespace`: Namespace of the TaskRun (string, optional, default: "default")

#### `delete_all_pipelineruns` – Delete multiple PipelineRuns based on selectors
- `namespace`: Namespace to delete PipelineRuns from (string, optional, default: "default")
- `labelSelector`: Label selector to filter PipelineRuns to delete (string, optional)
- `fieldSelector`: Field selector to filter PipelineRuns to delete (string, optional)

### Start/Restart Operations

#### `start_pipeline` – Start a Pipeline
- `name`: Name or reference of the Pipeline to start (string, required)
- `namespace`: Namespace where the Pipeline is located (string, optional, default: "default")

#### `start_task` – Start a Task
- `name`: Name or reference of the Task to start (string, required)
- `namespace`: Namespace where the Task is located (string, optional, default: "default")

#### `restart_pipelinerun` – Restart a PipelineRun
- `name`: Name or reference of the PipelineRun to restart (string, required)
- `namespace`: Namespace where the PipelineRun is located (string, optional, default: "default")

#### `restart_taskrun` – Restart a TaskRun
- `name`: Name or reference of the TaskRun to restart (string, required)
- `namespace`: Namespace where the TaskRun is located (string, optional, default: "default")

## Artifact Hub Integration

The MCP server provides integration with [Artifact Hub](https://artifacthub.io) to discover, install, and trigger Tekton tasks and pipelines from the community catalog.

### Artifact Hub Discovery Operations

#### `list_artifacthub_tasks` – List Tekton Tasks from Artifact Hub
- `query`: Search query to filter tasks (string, optional)
- `limit`: Maximum number of results to return (integer, optional, default: 20)

#### `list_artifacthub_pipelines` – List Tekton Pipelines from Artifact Hub
- `query`: Search query to filter pipelines (string, optional)
- `limit`: Maximum number of results to return (integer, optional, default: 20)

### Artifact Hub Installation Operations

#### `install_artifacthub_task` – Install a Tekton Task from Artifact Hub
- `packageId`: The Artifact Hub package ID of the task to install (string, required)
- `version`: Version of the task to install (string, optional)
- `namespace`: Namespace where the task will be installed (string, optional, default: "default")

#### `install_artifacthub_pipeline` – Install a Tekton Pipeline from Artifact Hub
- `packageId`: The Artifact Hub package ID of the pipeline to install (string, required)
- `version`: Version of the pipeline to install (string, optional)
- `namespace`: Namespace where the pipeline will be installed (string, optional, default: "default")

### Artifact Hub Trigger Operations

#### `trigger_artifacthub_task` – Trigger a Task installed from Artifact Hub
- `name`: Name of the installed task to trigger (string, required)
- `namespace`: Namespace where the task is located (string, optional, default: "default")
- `params`: Parameters to pass to the task (object, optional)

#### `trigger_artifacthub_pipeline` – Trigger a Pipeline installed from Artifact Hub
- `name`: Name of the installed pipeline to trigger (string, required)
- `namespace`: Namespace where the pipeline is located (string, optional, default: "default")
- `params`: Parameters to pass to the pipeline (object, optional)


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ftektoncd%2Fmcp-server.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Ftektoncd%2Fmcp-server?ref=badge_large)