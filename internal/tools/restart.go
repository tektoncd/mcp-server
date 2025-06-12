package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tektoncd/mcp-server/internal/params"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	pipelineclient "github.com/tektoncd/pipeline/pkg/client/injection/client"
	pipelineruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/pipelinerun"
	taskruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/taskrun"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func restartPipelineRun() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("restart_pipelinerun",
			mcp.WithDescription("Restart a PipelineRun"),
			mcp.WithString("name", mcp.Required(),
				mcp.Description("Name or Reference of the PipelineRun to restart"),
			),
			mcp.WithString("namespace",
				mcp.Description("Namespace where the PipelineRun is located"),
				mcp.DefaultString("default"),
			),
			// TODO add "parameters" objects
		),
		Handler: handlerRestartPipelineRun,
	}
}

func handlerRestartPipelineRun(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.GetArguments()["name"].(string)
	namespace, err := params.Optional[string](request, "namespace")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("namespace must be a string", err), nil
	}

	pipelinerunInformer := pipelineruninformer.Get(ctx)
	pipelineclientset := pipelineclient.Get(ctx)

	usepr, err := pipelinerunInformer.Lister().PipelineRuns(namespace).Get(name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to get PipelineRun %s/%s", namespace, name), err), nil
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
		return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to create PipelineRun %s/%s", namespace, pr.ObjectMeta.Name), err), nil
	}

	return result(fmt.Sprintf("Restarting pipelinerun %s as %s in namespace %s", name, pr.ObjectMeta.Name, namespace)), nil
}

func restartTaskRun() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("restart_taskrun",
			mcp.WithDescription("Restart a TaskRun"),
			mcp.WithString("name", mcp.Required(),
				mcp.Description("Name or Reference of the TaskRun to restart"),
			),
			mcp.WithString("namespace",
				mcp.Description("Namespace where the TaskRun is located"),
				mcp.DefaultString("default"),
			),
			// TODO add "parameters" objects
		),
		Handler: handlerRestartTaskRun,
	}
}

func handlerRestartTaskRun(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.GetArguments()["name"].(string)
	namespace, err := params.Optional[string](request, "namespace")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("namespace must be a string", err), nil
	}

	taskrunInformer := taskruninformer.Get(ctx)
	pipelineclientset := pipelineclient.Get(ctx)

	usetr, err := taskrunInformer.Lister().TaskRuns(namespace).Get(name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to get TaskRun %s/%s", namespace, name), err), nil
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
		return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to create TaskRun %s/%s", namespace, tr.ObjectMeta.Name), err), nil
	}

	return result(fmt.Sprintf("Restarting taskrun %s as %s in namespace %s", name, tr.ObjectMeta.Name, namespace)), nil
}
