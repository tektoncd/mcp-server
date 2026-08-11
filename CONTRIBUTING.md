# Contributing to Tekton MCP Server

Thank you for contributing time and expertise to Tekton MCP Server.

Before contributing, read and follow the project
[Code of Conduct](./code-of-conduct.md).

## Develop and test changes

See [`DEVELOPMENT.md`](./DEVELOPMENT.md) for repository setup, build and test
commands, local and cluster execution, dependency updates, and debugging.

The [Tekton community repository](https://github.com/tektoncd/community)
defines the shared contribution process:

- [development standards](https://github.com/tektoncd/community/blob/main/standards.md)
- [contacting the community](https://github.com/tektoncd/community/blob/main/contact.md)
- [finding work and proposing changes](https://github.com/tektoncd/community/tree/main/process)
- [code review](https://github.com/tektoncd/community/blob/main/process/README.md#reviews)
- [contributor ladder](https://github.com/tektoncd/community/blob/main/process/contributor-ladder.md)

The MCP Server project was proposed in
[`tektoncd/community#1194`](https://github.com/tektoncd/community/issues/1194).
Repository reviewers and approvers are listed in [`OWNERS`](./OWNERS).

## Submit a pull request

1. Work in a focused branch on a personal fork.
2. Add or update tests for behavior changes.
3. Run the checks documented in [`DEVELOPMENT.md`](./DEVELOPMENT.md).
4. Use the pull request template and include an appropriate release-note block.
5. Address reviewer feedback and keep the branch current with `main`.

Use [GitHub issues](https://github.com/tektoncd/mcp-server/issues) for bugs and
feature proposals. Report security vulnerabilities privately through the
[project security policy](https://github.com/tektoncd/mcp-server/security/policy),
not through a public issue.

## Documentation

User-facing changes should update the relevant repository documentation. For
changes to the Tekton website, follow the
[Tekton documentation contributor guide](https://github.com/tektoncd/website/blob/main/content/en/docs/Contribute/_index.md).

Tekton automation and repository infrastructure are maintained in
[`tektoncd/plumbing`](https://github.com/tektoncd/plumbing).
