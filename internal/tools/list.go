package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tektoncd/mcp-server/internal/params"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	pipelineinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/pipeline"
	pipelineruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/pipelinerun"
	taskinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/task"
	taskruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/taskrun"
	stepactioninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1beta1/stepaction"
	"k8s.io/apimachinery/pkg/labels"
)

func listTasks() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("list_tasks",
			mcp.WithDescription("List tasks in the cluster with filtering options"),
			mcp.WithString("namespace", mcp.Description("Which namespace to use to look for Task")),
			mcp.WithString("prefix", mcp.Description("Name prefix to filter Task")),
			mcp.WithString("label-selector", mcp.Description("Label selector to filter Task")),
		),
		Handler: handlerListTask,
	}
}

func handlerListTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskInformer := taskinformer.Get(ctx)
	namespace, err := params.Optional[string](request, "namespace")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	lselector, err := params.Optional[string](request, "label-selector")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	prefix, err := params.Optional[string](request, "prefix")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var selector labels.Selector
	if lselector != "" {
		selector, err = labels.Parse(lselector)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	} else {
		selector = labels.NewSelector()
	}

	var trs []*v1.Task

	if namespace == "" {
		// No namespace, searching all PipelineRuns
		trs, err = taskInformer.Lister().List(selector)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	} else {
		trs, err = taskInformer.Lister().Tasks(namespace).List(selector)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	}

	// Filter after the fact
	if prefix != "" {
		filteredTRs := []*v1.Task{}
		for _, pr := range trs {
			if strings.HasPrefix(pr.Name, prefix) {
				filteredTRs = append(filteredTRs, pr)
			}
		}
		trs = filteredTRs
	}

	jsonData, err := json.Marshal(trs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func listTaskRuns() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("list_taskruns",
			mcp.WithDescription("List taskruns in the cluster with filtering options"),
			mcp.WithString("namespace", mcp.Description("Which namespace to use to look for Taskruns")),
			mcp.WithString("prefix", mcp.Description("Name prefix to filter Taskruns")),
			mcp.WithString("label-selector", mcp.Description("Label selector to filter Taskruns")),
		),
		Handler: handlerListTaskRun,
	}
}

func handlerListTaskRun(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskRunInformer := taskruninformer.Get(ctx)
	namespace, err := params.Optional[string](request, "namespace")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	lselector, err := params.Optional[string](request, "label-selector")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	prefix, err := params.Optional[string](request, "prefix")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var selector labels.Selector
	if lselector != "" {
		selector, err = labels.Parse(lselector)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	} else {
		selector = labels.NewSelector()
	}

	var trs []*v1.TaskRun

	if namespace == "" {
		// No namespace, searching all PipelineRuns
		trs, err = taskRunInformer.Lister().List(selector)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	} else {
		trs, err = taskRunInformer.Lister().TaskRuns(namespace).List(selector)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	}

	// Filter after the fact
	if prefix != "" {
		filteredTRs := []*v1.TaskRun{}
		for _, pr := range trs {
			if strings.HasPrefix(pr.Name, prefix) {
				filteredTRs = append(filteredTRs, pr)
			}
		}
		trs = filteredTRs
	}

	jsonData, err := json.Marshal(trs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func listStepactions() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("list_stepactions",
			mcp.WithDescription("List stepactions in the cluster with filtering options"),
			mcp.WithString("namespace", mcp.Description("Which namespace to use to look for Stepactions")),
			mcp.WithString("prefix", mcp.Description("Name prefix to filter Stepactions")),
			mcp.WithString("label-selector", mcp.Description("Label selector to filter Stepactions")),
		),
		Handler: handlerListStepaction,
	}
}

func handlerListStepaction(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	stepactionInformer := stepactioninformer.Get(ctx)
	namespace, err := params.Optional[string](request, "namespace")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	lselector, err := params.Optional[string](request, "label-selector")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	prefix, err := params.Optional[string](request, "prefix")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var selector labels.Selector
	if lselector != "" {
		selector, err = labels.Parse(lselector)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	} else {
		selector = labels.NewSelector()
	}

	var trs []*v1beta1.StepAction

	if namespace == "" {
		// No namespace, searching all PipelineRuns
		trs, err = stepactionInformer.Lister().List(selector)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	} else {
		trs, err = stepactionInformer.Lister().StepActions(namespace).List(selector)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	}

	// Filter after the fact
	if prefix != "" {
		filteredTRs := []*v1beta1.StepAction{}
		for _, pr := range trs {
			if strings.HasPrefix(pr.Name, prefix) {
				filteredTRs = append(filteredTRs, pr)
			}
		}
		trs = filteredTRs
	}

	jsonData, err := json.Marshal(trs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func listPipelines() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("list_pipelines",
			mcp.WithDescription("List pipelines in the cluster with filtering options"),
			mcp.WithString("namespace", mcp.Description("Which namespace to use to look for Pipeline")),
			mcp.WithString("prefix", mcp.Description("Name prefix to filter Pipeline")),
			mcp.WithString("label-selector", mcp.Description("Label selector to filter Pipeline")),
		),
		Handler: handlerListPipeline,
	}
}

func handlerListPipeline(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pipelineInformer := pipelineinformer.Get(ctx)
	namespace, err := params.Optional[string](request, "namespace")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	lselector, err := params.Optional[string](request, "label-selector")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	prefix, err := params.Optional[string](request, "prefix")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var selector labels.Selector
	if lselector != "" {
		selector, err = labels.Parse(lselector)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	} else {
		selector = labels.NewSelector()
	}

	var prs []*v1.Pipeline

	if namespace == "" {
		// No namespace, searching all Pipelines
		prs, err = pipelineInformer.Lister().List(selector)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	} else {
		prs, err = pipelineInformer.Lister().Pipelines(namespace).List(selector)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	}

	// Filter after the fact
	if prefix != "" {
		filteredPRs := []*v1.Pipeline{}
		for _, pr := range prs {
			if strings.HasPrefix(pr.Name, prefix) {
				filteredPRs = append(filteredPRs, pr)
			}
		}
		prs = filteredPRs
	}

	jsonData, err := json.Marshal(prs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func listPipelineRuns() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("list_pipelineruns",
			mcp.WithDescription("List pipelineruns in the cluster with filtering options"),
			mcp.WithString("namespace", mcp.Description("Which namespace to use to look for PipelineRuns")),
			mcp.WithString("prefix", mcp.Description("Name prefix to filter PipelineRuns")),
			mcp.WithString("label-selector", mcp.Description("Label selector to filter PipelineRuns")),
		),
		Handler: handlerListPipelineRun,
	}
}

func handlerListPipelineRun(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pipelineRunInformer := pipelineruninformer.Get(ctx)
	namespace, err := params.Optional[string](request, "namespace")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	lselector, err := params.Optional[string](request, "label-selector")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	prefix, err := params.Optional[string](request, "prefix")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var selector labels.Selector
	if lselector != "" {
		selector, err = labels.Parse(lselector)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	} else {
		selector = labels.NewSelector()
	}

	var prs []*v1.PipelineRun

	if namespace == "" {
		// No namespace, searching all PipelineRuns
		prs, err = pipelineRunInformer.Lister().List(selector)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	} else {
		prs, err = pipelineRunInformer.Lister().PipelineRuns(namespace).List(selector)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	}

	// Filter after the fact
	if prefix != "" {
		filteredPRs := []*v1.PipelineRun{}
		for _, pr := range prs {
			if strings.HasPrefix(pr.Name, prefix) {
				filteredPRs = append(filteredPRs, pr)
			}
		}
		prs = filteredPRs
	}

	jsonData, err := json.Marshal(prs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}
