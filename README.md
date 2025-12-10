# Tektoncd Model Context Protocol server

*This project is in its early stages, and the README is currently minimal.*

This project provides a [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server for the tektoncd projects.
It initially focuses on [`tektoncd/pipeline`](https://github.com/tektoncd/pipeline) objects but will over time add support for other tektoncd projects.

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
