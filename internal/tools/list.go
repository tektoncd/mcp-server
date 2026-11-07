package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	pipelineinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/pipeline"
	pipelineruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/pipelinerun"
	taskinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/task"
	taskruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/taskrun"
	stepactioninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1beta1/stepaction"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"knative.dev/pkg/apis"
)

const (
	listOutputFull    = "full"
	listOutputSummary = "summary"
)

type listParams struct {
	Namespace     string `json:"namespace,omitempty"`
	LabelSelector string `json:"labelSelector,omitempty"`
	Prefix        string `json:"prefix,omitempty"`
	Output        string `json:"output,omitempty"`
}

func listSchema() (*jsonschema.Schema, error) {
	scheme, err := jsonschema.For[listParams](nil)
	if err != nil {
		return nil, err
	}
	scheme.Properties["namespace"].Description = "Namespace to list from; empty lists all namespaces"
	scheme.Properties["labelSelector"].Description = "Label selector to filter results"
	scheme.Properties["prefix"].Description = "Name prefix to filter results"
	scheme.Properties["output"].Description = "\"full\" returns the complete objects; \"summary\" returns name, namespace, creation time and, for runs, status, reason, start and completion time and the referenced Pipeline or Task"
	scheme.Properties["output"].Enum = []any{listOutputFull, listOutputSummary}
	scheme.Properties["output"].Default = json.RawMessage(`"full"`)
	return scheme, nil
}

// listSummary is the compact form of a listed object (output "summary").
type listSummary struct {
	Name              string       `json:"name"`
	Namespace         string       `json:"namespace"`
	CreationTimestamp metav1.Time  `json:"creationTimestamp"`
	Description       string       `json:"description,omitempty"`
	Ref               string       `json:"ref,omitempty"`
	PipelineRun       string       `json:"pipelineRun,omitempty"`
	PipelineTask      string       `json:"pipelineTask,omitempty"`
	Status            string       `json:"status,omitempty"`
	Reason            string       `json:"reason,omitempty"`
	StartTime         *metav1.Time `json:"startTime,omitempty"`
	CompletionTime    *metav1.Time `json:"completionTime,omitempty"`
}

func newListSummary(o metav1.Object) listSummary {
	return listSummary{Name: o.GetName(), Namespace: o.GetNamespace(), CreationTimestamp: o.GetCreationTimestamp()}
}

func (s *listSummary) setCondition(c *apis.Condition) {
	if c != nil {
		s.Status = string(c.Status)
		s.Reason = c.Reason
	}
}

func summarizeTask(t *v1.Task) listSummary {
	s := newListSummary(t)
	s.Description = t.Spec.Description
	return s
}

func summarizePipeline(p *v1.Pipeline) listSummary {
	s := newListSummary(p)
	s.Description = p.Spec.Description
	return s
}

func summarizeStepAction(a *v1beta1.StepAction) listSummary {
	s := newListSummary(a)
	s.Description = a.Spec.Description
	return s
}

func summarizePipelineRun(pr *v1.PipelineRun) listSummary {
	s := newListSummary(pr)
	if pr.Spec.PipelineRef != nil {
		s.Ref = pr.Spec.PipelineRef.Name
	}
	s.setCondition(pr.Status.GetCondition(apis.ConditionSucceeded))
	s.StartTime = pr.Status.StartTime
	s.CompletionTime = pr.Status.CompletionTime
	return s
}

func summarizeTaskRun(tr *v1.TaskRun) listSummary {
	s := newListSummary(tr)
	if tr.Spec.TaskRef != nil {
		s.Ref = tr.Spec.TaskRef.Name
	}
	s.PipelineRun = tr.Labels[pipeline.PipelineRunLabelKey]
	s.PipelineTask = tr.Labels[pipeline.PipelineTaskLabelKey]
	s.setCondition(tr.Status.GetCondition(apis.ConditionSucceeded))
	s.StartTime = tr.Status.StartTime
	s.CompletionTime = tr.Status.CompletionTime
	return s
}

// listResult marshals the listed objects, in full or summarized.
func listResult[T metav1.Object](items []T, output string, summarize func(T) listSummary) (*mcp.CallToolResult, error) {
	var data any
	switch output {
	case "", listOutputFull:
		data = items
	case listOutputSummary:
		summaries := make([]listSummary, 0, len(items))
		for _, item := range items {
			summaries = append(summaries, summarize(item))
		}
		data = summaries
	default:
		return nil, fmt.Errorf("invalid output %q: must be %q or %q", output, listOutputFull, listOutputSummary)
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}
	return result(string(jsonData)), nil
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

func listTasks() (serverTool, error) {
	scheme, err := listSchema()
	if err != nil {
		return nil, err
	}
	return newServerTool(
		"list_tasks",
		"List tasks in the cluster with filtering options",
		handlerListTasks,
		scheme,
	), nil
}

func handlerListTasks(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *callToolParamsFor[listParams],
) (*mcp.CallToolResult, error) {
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

	return listResult(trs, params.Arguments.Output, summarizeTask)
}

func listTaskRuns() (serverTool, error) {
	scheme, err := listSchema()
	if err != nil {
		return nil, err
	}
	return newServerTool(
		"list_taskruns",
		"List taskruns in the cluster with filtering options",
		handlerListTaskRuns,
		scheme,
	), nil
}

func handlerListTaskRuns(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *callToolParamsFor[listParams],
) (*mcp.CallToolResult, error) {
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

	return listResult(trs, params.Arguments.Output, summarizeTaskRun)
}

func listStepactions() (serverTool, error) {
	scheme, err := listSchema()
	if err != nil {
		return nil, err
	}
	return newServerTool(
		"list_stepactions",
		"List stepactions in the cluster with filtering options",
		handlerListStepactions,
		scheme,
	), nil
}

func handlerListStepactions(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *callToolParamsFor[listParams],
) (*mcp.CallToolResult, error) {
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

	return listResult(trs, params.Arguments.Output, summarizeStepAction)
}

func listPipelines() (serverTool, error) {
	scheme, err := listSchema()
	if err != nil {
		return nil, err
	}
	return newServerTool(
		"list_pipelines",
		"List pipelines in the cluster with filtering options",
		handlerListPipelines,
		scheme,
	), nil
}

func handlerListPipelines(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *callToolParamsFor[listParams],
) (*mcp.CallToolResult, error) {
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

	return listResult(prs, params.Arguments.Output, summarizePipeline)
}

func listPipelineRuns() (serverTool, error) {
	scheme, err := listSchema()
	if err != nil {
		return nil, err
	}
	return newServerTool(
		"list_pipelineruns",
		"List pipelineruns in the cluster with filtering options",
		handlerListPipelineRuns,
		scheme,
	), nil
}

func handlerListPipelineRuns(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *callToolParamsFor[listParams],
) (*mcp.CallToolResult, error) {
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

	return listResult(prs, params.Arguments.Output, summarizePipelineRun)
}
