//go:build e2e
// +build e2e

/*
Copyright 2024 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
)

func TestPipelineExecution(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	c, namespace := setup(ctx, t)

	pipelineName := generateName("exec-pipeline")
	
	// First create a pipeline to execute
	t.Run("Setup_Pipeline", func(t *testing.T) {
		pipelineYAML := fmt.Sprintf(`
apiVersion: tekton.dev/v1
kind: Pipeline
metadata:
  name: %s
  namespace: %s
spec:
  params:
  - name: message
    type: string
    default: "Hello from Pipeline"
  tasks:
  - name: first-task
    taskSpec:
      params:
      - name: message
        type: string
      results:
      - name: output
        description: The output message
      steps:
      - name: write-result
        image: mirror.gcr.io/ubuntu
        script: |
          #!/bin/bash
          echo -n "$(params.message)" > $(results.output.path)
          echo "Wrote: $(params.message)"
    params:
    - name: message
      value: "$(params.message)"
  - name: second-task
    runAfter: ["first-task"]
    taskSpec:
      params:
      - name: input-message
        type: string
      steps:
      - name: echo
        image: mirror.gcr.io/ubuntu
        script: |
          #!/bin/bash
          echo "Received: $(params.input-message)"
    params:
    - name: input-message
      value: "$(tasks.first-task.results.output)"
`, pipelineName, namespace)

		result, err := c.MCPSession.CallTool(ctx, "create_pipeline", map[string]interface{}{
			"namespace": namespace,
			"yaml":      pipelineYAML,
		})
		if err != nil {
			t.Fatalf("Failed to create pipeline: %v", err)
		}
		if !isSuccessResult(result) {
			t.Fatalf("Create pipeline failed: %v", result)
		}
	})

	var pipelineRunName string

	t.Run("Start_Pipeline", func(t *testing.T) {
		// Start pipeline via MCP
		result, err := c.MCPSession.CallTool(ctx, "start_pipeline", map[string]interface{}{
			"name":      pipelineName,
			"namespace": namespace,
		})
		if err != nil {
			t.Fatalf("Failed to start pipeline via MCP: %v", err)
		}

		content := getResultContent(result)
		t.Logf("Start pipeline result: %s", content)

		// Extract the PipelineRun name from the result
		// The result should contain the created PipelineRun name
		pipelineRuns, err := c.V1PipelineRunClient.List(ctx, metav1.ListOptions{
			LabelSelector: fmt.Sprintf("tekton.dev/pipeline=%s", pipelineName),
		})
		if err != nil {
			t.Fatalf("Failed to list PipelineRuns: %v", err)
		}
		if len(pipelineRuns.Items) == 0 {
			t.Fatal("No PipelineRun created after start_pipeline")
		}
		pipelineRunName = pipelineRuns.Items[0].Name
		t.Logf("Created PipelineRun: %s", pipelineRunName)
	})

	t.Run("Monitor_PipelineRun", func(t *testing.T) {
		if pipelineRunName == "" {
			t.Skip("No PipelineRun to monitor")
		}

		// Wait for PipelineRun to complete
		err := WaitForPipelineRunState(ctx, c, pipelineRunName, func(pr *v1.PipelineRun) bool {
			return pr.IsDone()
		}, "PipelineRunComplete")
		if err != nil {
			t.Fatalf("Error waiting for PipelineRun to complete: %v", err)
		}

		// Get the final status via MCP
		result, err := c.MCPSession.CallTool(ctx, "get_pipelinerun", map[string]interface{}{
			"name":      pipelineRunName,
			"namespace": namespace,
			"output":    "yaml",
		})
		if err != nil {
			t.Fatalf("Failed to get PipelineRun via MCP: %v", err)
		}

		content := getResultContent(result)
		if !contains(content, pipelineRunName) {
			t.Errorf("PipelineRun name not found in response")
		}

		// Verify via K8s API that it succeeded
		pr, err := c.V1PipelineRunClient.Get(ctx, pipelineRunName, metav1.GetOptions{})
		if err != nil {
			t.Fatalf("Failed to get PipelineRun: %v", err)
		}

		if !IsSuccessful(t, pr.Status.Status) {
			t.Errorf("PipelineRun did not succeed: %v", pr.Status.GetCondition(apis.ConditionSucceeded))
		}
	})
}

func TestTaskExecution(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	c, namespace := setup(ctx, t)

	taskName := generateName("exec-task")

	// Create a task to execute
	t.Run("Setup_Task", func(t *testing.T) {
		taskYAML := fmt.Sprintf(`
apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: %s
  namespace: %s
spec:
  params:
  - name: message
    type: string
    default: "Hello from Task"
  results:
  - name: timestamp
    description: Current timestamp
  steps:
  - name: echo-and-timestamp
    image: mirror.gcr.io/ubuntu
    script: |
      #!/bin/bash
      echo "Message: $(params.message)"
      date +%%s > $(results.timestamp.path)
`, taskName, namespace)

		result, err := c.MCPSession.CallTool(ctx, "create_task", map[string]interface{}{
			"namespace": namespace,
			"yaml":      taskYAML,
		})
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}
		if !isSuccessResult(result) {
			t.Fatalf("Create task failed: %v", result)
		}
	})

	var taskRunName string

	t.Run("Start_Task", func(t *testing.T) {
		// Start task via MCP
		result, err := c.MCPSession.CallTool(ctx, "start_task", map[string]interface{}{
			"name":      taskName,
			"namespace": namespace,
		})
		if err != nil {
			t.Fatalf("Failed to start task via MCP: %v", err)
		}

		content := getResultContent(result)
		t.Logf("Start task result: %s", content)

		// Find the created TaskRun
		taskRuns, err := c.V1TaskRunClient.List(ctx, metav1.ListOptions{
			LabelSelector: fmt.Sprintf("tekton.dev/task=%s", taskName),
		})
		if err != nil {
			t.Fatalf("Failed to list TaskRuns: %v", err)
		}
		if len(taskRuns.Items) == 0 {
			t.Fatal("No TaskRun created after start_task")
		}
		taskRunName = taskRuns.Items[0].Name
		t.Logf("Created TaskRun: %s", taskRunName)
	})

	t.Run("Get_TaskRun_Logs", func(t *testing.T) {
		if taskRunName == "" {
			t.Skip("No TaskRun to get logs from")
		}

		// Wait for TaskRun to complete
		err := WaitForTaskRunState(ctx, c, taskRunName, func(tr *v1.TaskRun) bool {
			return tr.IsDone()
		}, "TaskRunComplete")
		if err != nil {
			t.Fatalf("Error waiting for TaskRun to complete: %v", err)
		}

		// Get logs via MCP
		result, err := c.MCPSession.CallTool(ctx, "get_taskrun_logs", map[string]interface{}{
			"name":      taskRunName,
			"namespace": namespace,
		})
		if err != nil {
			t.Fatalf("Failed to get TaskRun logs via MCP: %v", err)
		}

		logs := getResultContent(result)
		if logs == "" {
			t.Error("No logs returned from get_taskrun_logs")
		}
		t.Logf("TaskRun logs: %s", logs)

		// Verify logs contain expected message
		if !contains(logs, "Message:") {
			t.Error("Logs don't contain expected output")
		}
	})
}

func TestPipelineRunRestart(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	c, namespace := setup(ctx, t)

	// Create a pipeline that will fail
	failingPipelineName := generateName("failing-pipeline")
	
	t.Run("Setup_Failing_Pipeline", func(t *testing.T) {
		pipelineYAML := fmt.Sprintf(`
apiVersion: tekton.dev/v1
kind: Pipeline
metadata:
  name: %s
  namespace: %s
spec:
  tasks:
  - name: fail-task
    taskSpec:
      steps:
      - name: fail
        image: mirror.gcr.io/ubuntu
        script: |
          #!/bin/bash
          echo "This task will fail"
          exit 1
`, failingPipelineName, namespace)

		result, err := c.MCPSession.CallTool(ctx, "create_pipeline", map[string]interface{}{
			"namespace": namespace,
			"yaml":      pipelineYAML,
		})
		if err != nil {
			t.Fatalf("Failed to create failing pipeline: %v", err)
		}
		if !isSuccessResult(result) {
			t.Fatalf("Create pipeline failed: %v", result)
		}
	})

	var failedRunName string

	t.Run("Create_Failed_Run", func(t *testing.T) {
		// Create a PipelineRun that will fail
		pipelineRunYAML := fmt.Sprintf(`
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  generateName: %s-run-
  namespace: %s
spec:
  pipelineRef:
    name: %s
`, failingPipelineName, namespace, failingPipelineName)

		result, err := c.MCPSession.CallTool(ctx, "create_pipelinerun", map[string]interface{}{
			"namespace":    namespace,
			"yaml":         pipelineRunYAML,
			"generateName": fmt.Sprintf("%s-run-", failingPipelineName),
		})
		if err != nil {
			t.Fatalf("Failed to create PipelineRun: %v", err)
		}

		// Find the created PipelineRun
		time.Sleep(2 * time.Second) // Give it a moment to be created
		pipelineRuns, err := c.V1PipelineRunClient.List(ctx, metav1.ListOptions{
			LabelSelector: fmt.Sprintf("tekton.dev/pipeline=%s", failingPipelineName),
		})
		if err != nil {
			t.Fatalf("Failed to list PipelineRuns: %v", err)
		}
		if len(pipelineRuns.Items) == 0 {
			t.Fatal("No PipelineRun created")
		}
		failedRunName = pipelineRuns.Items[0].Name

		// Wait for it to fail
		err = WaitForPipelineRunState(ctx, c, failedRunName, func(pr *v1.PipelineRun) bool {
			return pr.IsDone() && !pr.IsSuccessful()
		}, "PipelineRunFailed")
		if err != nil {
			t.Fatalf("Error waiting for PipelineRun to fail: %v", err)
		}
	})

	t.Run("Restart_Failed_PipelineRun", func(t *testing.T) {
		if failedRunName == "" {
			t.Skip("No failed PipelineRun to restart")
		}

		// Restart the failed PipelineRun via MCP
		result, err := c.MCPSession.CallTool(ctx, "restart_pipelinerun", map[string]interface{}{
			"name":      failedRunName,
			"namespace": namespace,
		})
		if err != nil {
			t.Fatalf("Failed to restart PipelineRun via MCP: %v", err)
		}

		content := getResultContent(result)
		t.Logf("Restart result: %s", content)

		// Verify a new PipelineRun was created
		time.Sleep(2 * time.Second)
		pipelineRuns, err := c.V1PipelineRunClient.List(ctx, metav1.ListOptions{
			LabelSelector: fmt.Sprintf("tekton.dev/pipeline=%s", failingPipelineName),
		})
		if err != nil {
			t.Fatalf("Failed to list PipelineRuns after restart: %v", err)
		}
		if len(pipelineRuns.Items) < 2 {
			t.Error("No new PipelineRun created after restart")
		}
	})
}

func TestListOperationsWithFilters(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	c, namespace := setup(ctx, t)

	// Create multiple pipelines with different labels
	t.Run("Setup_Multiple_Resources", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			pipelineYAML := fmt.Sprintf(`
apiVersion: tekton.dev/v1
kind: Pipeline
metadata:
  name: filter-test-%d
  namespace: %s
  labels:
    environment: %s
    team: platform
spec:
  tasks:
  - name: dummy
    taskSpec:
      steps:
      - name: echo
        image: mirror.gcr.io/ubuntu
        script: echo "test"
`, i, namespace, map[int]string{0: "dev", 1: "staging", 2: "dev"}[i])

			result, err := c.MCPSession.CallTool(ctx, "create_pipeline", map[string]interface{}{
				"namespace": namespace,
				"yaml":      pipelineYAML,
			})
			if err != nil {
				t.Fatalf("Failed to create pipeline %d: %v", i, err)
			}
			if !isSuccessResult(result) {
				t.Fatalf("Create pipeline %d failed", i)
			}
		}
	})

	t.Run("List_With_Prefix_Filter", func(t *testing.T) {
		result, err := c.MCPSession.CallTool(ctx, "list_pipelines", map[string]interface{}{
			"namespace": namespace,
			"prefix":    "filter-test",
		})
		if err != nil {
			t.Fatalf("Failed to list pipelines with prefix: %v", err)
		}

		content := getResultContent(result)
		// Should contain all 3 pipelines
		for i := 0; i < 3; i++ {
			if !contains(content, fmt.Sprintf("filter-test-%d", i)) {
				t.Errorf("Pipeline filter-test-%d not found in filtered list", i)
			}
		}
	})

	t.Run("List_With_Label_Selector", func(t *testing.T) {
		result, err := c.MCPSession.CallTool(ctx, "list_pipelines", map[string]interface{}{
			"namespace":      namespace,
			"label-selector": "environment=dev",
		})
		if err != nil {
			t.Fatalf("Failed to list pipelines with label selector: %v", err)
		}

		content := getResultContent(result)
		// Should contain only the dev pipelines (0 and 2)
		if !contains(content, "filter-test-0") {
			t.Error("Pipeline filter-test-0 not found in label-filtered list")
		}
		if !contains(content, "filter-test-2") {
			t.Error("Pipeline filter-test-2 not found in label-filtered list")
		}
		if contains(content, "filter-test-1") {
			t.Error("Pipeline filter-test-1 should not be in dev-filtered list")
		}
	})
}