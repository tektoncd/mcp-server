package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	pipelineclient "github.com/tektoncd/pipeline/pkg/client/injection/client"
	pipelineruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/pipelinerun"
	taskruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/taskrun"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type restartParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func restartSchema() mcp.ToolOption {
	scheme, err := jsonschema.For[restartParams]()
	if err != nil {
		panic(err)
	}

	scheme.Properties["name"].Description = "Name or referece of the object"
	scheme.Properties["namespace"].Description = "Namespace of the object"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)

	return mcp.Input(mcp.Schema(scheme))
}

func restartPipelineRun() *mcp.ServerTool {
	return mcp.NewServerTool(
		"restart_pipelinerun",
		"Restart a PipelineRun",
		handlerRestartPipelineRun,
		restartSchema(),
	)
}

func handlerRestartPipelineRun(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[restartParams],
) (*mcp.CallToolResultFor[string], error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace

	pipelinerunInformer := pipelineruninformer.Get(ctx)
	pipelineclientset := pipelineclient.Get(ctx)

	usepr, err := pipelinerunInformer.Lister().PipelineRuns(namespace).Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get PipelineRun %s/%s: %w", namespace, name, err)
	}
	pr := &v1.PipelineRun{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "tekton.dev/v1",
			Kind:       "PipelineRun",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    namespace,
			GenerateName: name + "-",
		},
		Spec: usepr.Spec,
	}
	pr.Spec.Status = ""
	if len(usepr.ObjectMeta.GenerateName) > 0 {
		pr.ObjectMeta.GenerateName = usepr.ObjectMeta.GenerateName
	}

	pr, err = pipelineclientset.TektonV1().PipelineRuns(namespace).Create(ctx, pr, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create PipelineRun %s/%s: %w", namespace, pr.ObjectMeta.Name, err)
	}

	return result(fmt.Sprintf("Restarting pipelinerun %s as %s in namespace %s", name, pr.ObjectMeta.Name, namespace)), nil
}

func restartTaskRun() *mcp.ServerTool {
	return mcp.NewServerTool(
		"restart_taskrun",
		"Restart a TaskRun",
		handlerRestartTaskRun,
		restartSchema(),
	)
}

func handlerRestartTaskRun(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[restartParams],
) (*mcp.CallToolResultFor[string], error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace

	taskrunInformer := taskruninformer.Get(ctx)
	pipelineclientset := pipelineclient.Get(ctx)

	usetr, err := taskrunInformer.Lister().TaskRuns(namespace).Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get TaskRun %s/%s: %w", namespace, name, err)
	}

	tr := &v1.TaskRun{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "tekton.dev/v1",
			Kind:       "TaskRun",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    namespace,
			GenerateName: name + "-",
		},
		Spec: usetr.Spec,
	}
	tr.Spec.Status = ""
	if len(usetr.ObjectMeta.GenerateName) > 0 {
		tr.ObjectMeta.GenerateName = usetr.ObjectMeta.GenerateName
	}

	tr, err = pipelineclientset.TektonV1().TaskRuns(namespace).Create(ctx, tr, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create TaskRun %s/%s: %w", namespace, tr.ObjectMeta.Name, err)
	}

	return result(fmt.Sprintf("Restarting taskrun %s as %s in namespace %s", name, tr.ObjectMeta.Name, namespace)), nil
}
