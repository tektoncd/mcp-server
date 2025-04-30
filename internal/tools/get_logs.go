package tools

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tektoncd/mcp-server/internal/params"
	taskruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/taskrun"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
)

func getTaskRunLogs() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("get_taskrun_logs",
			mcp.WithDescription("Get the logs for a given TaskRun"),
			mcp.WithString("name", mcp.Required(),
				mcp.Description("Name of the TaskRun to get logs for"),
			),
			mcp.WithString("namespace",
				mcp.Description("Namespace where the TaskRun is located"),
				mcp.DefaultString("default"),
			),
		),
		Handler: handlerGetTaskRunLogs,
	}
}

func handlerGetTaskRunLogs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.GetArguments()["name"].(string)
	if !ok {
		return mcp.NewToolResultError("name must be a string"), nil
	}
	namespace, err := params.Optional[string](request, "namespace")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("namespace must be a string", err), nil
	}

	taskrunInformer := taskruninformer.Get(ctx)
	kubeclientset := kubeclient.Get(ctx)

	task, err := taskrunInformer.Lister().TaskRuns(namespace).Get(name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to get TaskRun %s/%s", namespace, name), err), nil
	}

	podName := task.Status.PodName
	if podName == "" {
		return mcp.NewToolResultError(fmt.Sprintf("PodName not set for TaskRun %s/%s", namespace, name)), nil
	}

	logs, err := getLogs(ctx, kubeclientset.CoreV1().Pods(namespace), podName)
	if err != nil {
		return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to get logs for TaskRun %s/%s", namespace, name), err), nil
	}

	return result(logs), nil
}

func getLogs(ctx context.Context, client corev1.PodInterface, name string) (string, error) {
	pod, err := client.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get Pod %s: %w", name, err)
	}
	var sb strings.Builder
	for _, container := range pod.Spec.Containers {
		sb.WriteString(fmt.Sprintf("\n>>> Pod %s Container %s\n", pod.Name, container.Name))
		req := client.GetLogs(pod.Name, &v1.PodLogOptions{Follow: false, Container: container.Name})
		res, err := req.Stream(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to get container %q logs for Pod %s: %w", container.Name, name, err)
		}
		defer res.Close()
		data, err := io.ReadAll(res)
		if err != nil {
			return "", fmt.Errorf("failed to read response for container %q logs for Pod %s: %w", container.Name, name, err)
		}
		sb.Write(data)
	}
	return sb.String(), nil
}
