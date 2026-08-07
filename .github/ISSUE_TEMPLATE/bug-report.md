---
name: Bug Report
about: Template for bug reports
labels: kind/bug
---

# Expected Behavior

# Actual Behavior

# Steps to Reproduce the Problem

1.
2.
3.

# Additional Info

- MCP Server version (image tag/digest or binary version/commit):

```
(paste your output here)
```

- MCP client and transport:

```
(paste the client name/version and transport here)
```

- Deployment mode and relevant configuration:

```
(local binary, container, or Kubernetes; remove sensitive values)
```

- Kubernetes version:

  **Output of `kubectl version`:**

```
(paste your output here)
```

- Tekton Pipeline version:

  **Output of `tkn version` or `kubectl get pods -n tekton-pipelines -l app=tekton-pipelines-controller -o=jsonpath='{.items[0].metadata.labels.version}'`**

```
(paste your output here)
```


<!-- Any other additional information -->
