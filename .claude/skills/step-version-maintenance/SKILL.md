---
name: step-version-maintenance
description: Instructions for upgrading the pinned major versions in starter workflows to the latest step major versions.
disable-model-invocation: true
allowed-tools: mcp__bitrise__step_search
---

### Context

This repo contains starter Bitrise workflows for various project types. The workflows naturally contain Bitrise steps, and as a best practice, the major version of those steps is pinned (automatically receive patch and minor updates, but not major updates). Therefore, when a new major version of a step is released, the pinned major version in the starter workflows needs to be updated to ensure that users get the latest features and improvements when they use the starter workflows.

Those major versions are defined in @steps/const.go.

### Instructions

Prerequisites:

- Access to the Bitrise MCP and its step search tool
- A proper Go environment according to @.tool-versions

If the above are not met, do not proceed, just flag the issue to the user.


1. Identify the steps used in starter workflows by reading @steps/const.go.
2. For each step, check the current latest major version. You can do this via the Bitrise MCP step search tool.
3. Update the major version in @steps/const.go to the latest major version for each step.
4. Run Go tests to verify your changes. There might be failing tests unrelated to the changes (mostly tooling issues). In this case, go ahead and let CI be the judge.
5. Create a new branch, commit your changes, and open a PR.
