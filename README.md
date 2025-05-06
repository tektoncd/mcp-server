# Tektoncd Model Context Protocol server

*This project is in its early stages, and the README is currently minimal.*.

This project provides a [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server for the tektoncd projects.
It initially focuses on [`tektoncd/pipeline`](https://github.com/tektoncd/pipeline) objects but will over time add support for other tektoncd projects.

## Use Cases

- Inspect Tekton Resources: Quickly list Pipelines, Tasks, PipelineRuns, TaskRuns, and StepActions with filtering options (by namespace, name prefix, or label selectors).
- Trigger Executions: Start new Pipeline or Task executions with ease using simple parameter inputs.
- Re-run Executions: Restart existing PipelineRuns or TaskRuns for retry workflows or debugging.

## Prerequisites

1. You’ll need access to a Kubernetes cluster where Tekton Pipelines is already installed and running. The MCP server interacts with Tekton resources inside the cluster. so make sure you have access to a valid KUBECONFIG file that can connect to the cluster.
1. Make sure you have Docker installed on your machine if you plan to run the server using a container, and verify that Docker is running properly. 

## Installation

### Usage with mcp-cli

Add the following JSON block to server_config_json.

```json
{
  "mcpServers": {
    "tekton": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-v<KUBE_CONFIG_PATH>:/home/nonroot/.kube/config",
        "quay.io/jkhelil/tektoncd/tekton-mcp-server"
      ]
    }
  }
}
```

then run the following command to start the MCP server:

```bash
uv run mcp-cli chat --server tekton --provider ollama --model llama3.2
```
This will start the MCP server and allow you to interact with it using the mcp-cli tool.


### Usage with Claude Desktop

Add the following JSON block to Claude's Developer settings.

```json
{
  "mcpServers": {
    "tekton": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-v<KUBE_CONFIG_PATH>:/home/nonroot/.kube/config",
        "quay.io/jkhelil/tektoncd/tekton-mcp-server"
      ]
    }
  }
}
```

### Build from source

You can use go build to build the binary in the cmd/tekton-mcp-server directory, and configure your server to use the built executable as its command. For example

```json
{
  "mcpServers": {
    "tekton": {
        "command": "/path/to/tekton-mcp-server",
        "args": ["stdio"]
    }
  }
}
```

## Tools

#### `list_pipelines` – List Pipelines in the Cluster with Filtering Options  
- `namespace`: Namespace to list Pipelines from (string, required)  
- `prefix`: Name prefix to filter Pipelines (string, optional)  
- `label-selector`: Label selector to filter Pipelines (string, optional)  

#### `start_pipeline` – Start a Pipeline  
- `name`: Name or reference of the Pipeline to start (string, required)  
- `namespace`: Namespace where the Pipeline is located (string, optional, default: "default")  


#### `list_pipeline_runs` – List PipelineRuns in the Cluster with Filtering Options  
- `namespace`: Namespace to list PipelineRuns from (string, required)  
- `prefix`: Name prefix to filter PipelineRuns (string, optional)  
- `label-selector`: Label selector to filter PipelineRuns (string, optional)  

#### `restart_pipelinerun` – Restart a PipelineRun  
- `name`: Name or reference of the PipelineRun to restart (string, required)  
- `namespace`: Namespace where the PipelineRun is located (string, optional, default: "default")  

#### `list_tasks` – List Tasks in the Cluster with Filtering Options  
- `namespace`: Namespace to list Tasks from (string, required)  
- `prefix`: Name prefix to filter Tasks (string, optional)  
- `label-selector`: Label selector to filter Tasks (string, optional)  

#### `start_task` – Start a Task  
- `name`: Name or reference of the Task to start (string, required)  
- `namespace`: Namespace where the Task is located (string, optional, default: "default")  

#### `list_task_runs` – List TaskRuns in the Cluster with Filtering Options  
- `namespace`: Namespace to list TaskRuns from (string, required)  
- `prefix`: Name prefix to filter TaskRuns (string, optional)  
- `label-selector`: Label selector to filter TaskRuns (string, optional)  

#### `restart_taskrun` – Restart a TaskRun  
- `name`: Name or reference of the TaskRun to restart (string, required)  
- `namespace`: Namespace where the TaskRun is located (string, optional, default: "default")  

#### `list_stepactions` – List Step Actions in the Cluster with Filtering Options  
- `namespace`: Namespace to list Step Actions from (string, required)  
- `prefix`: Name prefix to filter Step Actions (string, optional)  
- `label-selector`: Label selector to filter Step Actions (string, optional)
