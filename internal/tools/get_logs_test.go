package tools

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mark3labs/mcp-go/mcp"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	ttesting "github.com/tektoncd/pipeline/pkg/reconciler/testing"
	"github.com/tektoncd/pipeline/test"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	fakerest "k8s.io/client-go/rest/fake"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
)

type fakeClient struct {
	*fake.Clientset
	logs map[string]map[string]string // podname:container:log
}

func (c *fakeClient) CoreV1() typedv1.CoreV1Interface {
	return &fakeCoreV1Client{
		CoreV1Interface: c.Clientset.CoreV1(),
		logs:            c.logs,
	}
}

type fakeCoreV1Client struct {
	typedv1.CoreV1Interface
	logs map[string]map[string]string
}

func (c *fakeCoreV1Client) Pods(namespace string) typedv1.PodInterface {
	return &fakePodV1Client{
		PodInterface: c.CoreV1Interface.Pods(namespace),
		logs:         c.logs,
	}
}

type fakePodV1Client struct {
	typedv1.PodInterface
	logs map[string]map[string]string
}

func (f *fakePodV1Client) GetLogs(name string, opts *corev1.PodLogOptions) *restclient.Request {
	statusCode := http.StatusOK
	pod, ok := f.logs[name]
	if !ok {
		statusCode = http.StatusNotFound
	}
	containerLogs, ok := pod[opts.Container]
	if !ok {
		statusCode = http.StatusNotFound
	}
	fakeClient := &fakerest.RESTClient{
		Client: fakerest.CreateHTTPClient(func(request *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: statusCode,
				Body: io.NopCloser(
					strings.NewReader(containerLogs)),
			}
			return resp, nil
		}),
	}
	return fakeClient.Request()
}

func TestHandlerGetTaskRunLogs(t *testing.T) {
	data := test.Data{
		TaskRuns: []*v1.TaskRun{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "hello-world", Namespace: "default"},
				Status:     v1.TaskRunStatus{TaskRunStatusFields: v1.TaskRunStatusFields{PodName: "hello-world"}},
			},
		},
		Pods: []*corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "hello-world", Namespace: "default"},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "hello"}, {Name: "goodbye"}}},
			},
		},
	}

	ctx, _ := ttesting.SetupFakeContext(t)
	clients, _ := test.SeedTestData(t, ctx, data)
	kubeclientset := &fakeClient{
		Clientset: clients.Kube,
		logs: map[string]map[string]string{
			"hello-world": {
				"hello":   "Hello, World!",
				"goodbye": "Goodbye!",
			},
		},
	}
	ctx = context.WithValue(ctx, kubeclient.Key{}, kubeclientset)

	request := newCallToolRequest(map[string]any{"name": "hello-world", "namespace": "default"})
	expected := mcp.NewTextContent(`
>>> Pod hello-world Container hello
Hello, World!
>>> Pod hello-world Container goodbye
Goodbye!`)

	result, err := handlerGetTaskRunLogs(ctx, request)
	if err != nil {
		t.Fatal(err)
	}

	received, _ := mcp.AsTextContent(result.Content[0])
	if diff := cmp.Diff(expected.Text, received.Text); diff != "" {
		t.Errorf("getLogs mismatch (-want +got):\n%s", diff)
	}
}
