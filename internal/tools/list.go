package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	pipelineinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/pipeline"
	pipelineruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/pipelinerun"
	taskinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/task"
	taskruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/taskrun"
	stepactioninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1beta1/stepaction"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type listParams struct {
	Namespace     string `json:"namespace"`
	LabelSelector string `json:"labelSelector"`
	Prefix        string `json:"prefix"`
}

func parseLabelSelector(lselector string) (labels.Selector, error) {
	if lselector == "" {
		return labels.NewSelector(), nil
	}
	return labels.Parse(lselector)
}

func filterList[T metav1.Object](in []T, prefix string) []T {
	out := make([]T, 0, len(in))
	for _, item := range in {
		if strings.HasPrefix(item.GetName(), prefix) {
			out = append(out, item)
		}
	}
	return out
}

func listTasks() *mcp.ServerTool {
	return mcp.NewServerTool(
		"list_tasks",
		"List tasks in the cluster with filtering options",
		handlerListTasks,
	)
}

func handlerListTasks(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[listParams],
) (*mcp.CallToolResultFor[string], error) {
	namespace := params.Arguments.Namespace
	lselector := params.Arguments.LabelSelector
	prefix := params.Arguments.Prefix

	selector, err := parseLabelSelector(lselector)
	if err != nil {
		return nil, err
	}

	taskInformer := taskinformer.Get(ctx)

	var trs []*v1.Task

	if namespace == "" {
		// No namespace, searching all PipelineRuns
		trs, err = taskInformer.Lister().List(selector)
		if err != nil {
			return nil, err
		}
	} else {
		trs, err = taskInformer.Lister().Tasks(namespace).List(selector)
		if err != nil {
			return nil, err
		}
	}

	// Filter after the fact
	if prefix != "" {
		trs = filterList(trs, prefix)
	}

	jsonData, err := json.Marshal(trs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}

	return result(string(jsonData)), nil
}

func listTaskRuns() *mcp.ServerTool {
	return mcp.NewServerTool(
		"list_taskruns",
		"List taskruns in the cluster with filtering options",
		handlerListTaskRuns,
	)
}

func handlerListTaskRuns(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[listParams],
) (*mcp.CallToolResultFor[string], error) {
	namespace := params.Arguments.Namespace
	lselector := params.Arguments.LabelSelector
	prefix := params.Arguments.Prefix

	selector, err := parseLabelSelector(lselector)
	if err != nil {
		return nil, err
	}

	taskRunInformer := taskruninformer.Get(ctx)
	var trs []*v1.TaskRun

	if namespace == "" {
		// No namespace, searching all PipelineRuns
		trs, err = taskRunInformer.Lister().List(selector)
		if err != nil {
			return nil, err
		}
	} else {
		trs, err = taskRunInformer.Lister().TaskRuns(namespace).List(selector)
		if err != nil {
			return nil, err
		}
	}

	// Filter after the fact
	if prefix != "" {
		trs = filterList(trs, prefix)
	}

	jsonData, err := json.Marshal(trs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}

	return result(string(jsonData)), nil
}

func listStepactions() *mcp.ServerTool {
	return mcp.NewServerTool(
		"list_stepactions",
		"List stepactions in the cluster with filtering options",
		handlerListStepactions,
	)
}

func handlerListStepactions(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[listParams],
) (*mcp.CallToolResultFor[string], error) {
	namespace := params.Arguments.Namespace
	lselector := params.Arguments.LabelSelector
	prefix := params.Arguments.Prefix

	selector, err := parseLabelSelector(lselector)
	if err != nil {
		return nil, err
	}

	stepactionInformer := stepactioninformer.Get(ctx)
	var trs []*v1beta1.StepAction

	if namespace == "" {
		// No namespace, searching all PipelineRuns
		trs, err = stepactionInformer.Lister().List(selector)
		if err != nil {
			return nil, err
		}
	} else {
		trs, err = stepactionInformer.Lister().StepActions(namespace).List(selector)
		if err != nil {
			return nil, err
		}
	}

	// Filter after the fact
	if prefix != "" {
		trs = filterList(trs, prefix)
	}

	jsonData, err := json.Marshal(trs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}

	return result(string(jsonData)), nil
}

func listPipelines() *mcp.ServerTool {
	return mcp.NewServerTool(
		"list_pipelines",
		"List pipelines in the cluster with filtering options",
		handlerListPipelines,
	)
}

func handlerListPipelines(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[listParams],
) (*mcp.CallToolResultFor[string], error) {
	namespace := params.Arguments.Namespace
	lselector := params.Arguments.LabelSelector
	prefix := params.Arguments.Prefix

	selector, err := parseLabelSelector(lselector)
	if err != nil {
		return nil, err
	}

	pipelineInformer := pipelineinformer.Get(ctx)
	var prs []*v1.Pipeline

	if namespace == "" {
		// No namespace, searching all Pipelines
		prs, err = pipelineInformer.Lister().List(selector)
		if err != nil {
			return nil, err
		}
	} else {
		prs, err = pipelineInformer.Lister().Pipelines(namespace).List(selector)
		if err != nil {
			return nil, err
		}
	}

	// Filter after the fact
	if prefix != "" {
		prs = filterList(prs, prefix)
	}

	jsonData, err := json.Marshal(prs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}

	return result(string(jsonData)), nil
}

func listPipelineRuns() *mcp.ServerTool {
	return mcp.NewServerTool(
		"list_pipelineruns",
		"List pipelineruns in the cluster with filtering options",
		handlerListPipelineRuns,
	)
}

func handlerListPipelineRuns(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[listParams],
) (*mcp.CallToolResultFor[string], error) {
	namespace := params.Arguments.Namespace
	lselector := params.Arguments.LabelSelector
	prefix := params.Arguments.Prefix

	selector, err := parseLabelSelector(lselector)
	if err != nil {
		return nil, err
	}

	pipelineRunInformer := pipelineruninformer.Get(ctx)
	var prs []*v1.PipelineRun

	if namespace == "" {
		// No namespace, searching all PipelineRuns
		prs, err = pipelineRunInformer.Lister().List(selector)
		if err != nil {
			return nil, err
		}
	} else {
		prs, err = pipelineRunInformer.Lister().PipelineRuns(namespace).List(selector)
		if err != nil {
			return nil, err
		}
	}

	// Filter after the fact
	if prefix != "" {
		prs = filterList(prs, prefix)
	}

	jsonData, err := json.Marshal(prs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}

	return result(string(jsonData)), nil
}
