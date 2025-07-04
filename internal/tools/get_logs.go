package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
	taskruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/taskrun"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
)

type getLogsParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func getTaskRunLogsSchema() (mcp.ToolOption, error) {
	scheme, err := jsonschema.For[getLogsParams]()
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name or referece of the object"
	scheme.Properties["namespace"].Description = "Namespace of the object"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)

	return mcp.Input(mcp.Schema(scheme)), nil
}

func getTaskRunLogs() (*mcp.ServerTool, error) {
	schema, err := getTaskRunLogsSchema()
	if err != nil {
		return nil, err
	}
	return mcp.NewServerTool(
		"get_taskrun_logs",
		"Get the logs for a given TaskRun",
		handlerGetTaskRunLogs,
		schema,
	), nil
}

func handlerGetTaskRunLogs(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[getLogsParams],
) (*mcp.CallToolResultFor[string], error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace

	taskrunInformer := taskruninformer.Get(ctx)
	kubeclientset := kubeclient.Get(ctx)

	task, err := taskrunInformer.Lister().TaskRuns(namespace).Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get TaskRun %s/%s: %w", namespace, name, err)
	}

	podName := task.Status.PodName
	if podName == "" {
		return nil, fmt.Errorf("podName not set for TaskRun %s/%s", namespace, name)
	}

	logs, err := getLogs(ctx, kubeclientset.CoreV1().Pods(namespace), podName)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs for TaskRun %s/%s: %w", namespace, name, err)
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
