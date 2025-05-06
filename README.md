# Tektoncd Model Context Protocol server

*This project is in its early stages, and the README is currently minimal.*.

This project provides a [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server for the tektoncd projects.
It initially focuses on [`tektoncd/pipeline`](https://github.com/tektoncd/pipeline) objects but will over time add support for other tektoncd projects.

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

#### `get_taskrun_logs` - Get the logs for a given TaskRun
- `name`: Name or reference of the TaskRun to restart (string, required)
- `namespace`: Namespace where the TaskRun is located (string, optional, default: "default")
