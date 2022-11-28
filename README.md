# gh-release-with-wf-dispatch

gh-release-with-wf-dispatch crates github tag with gh-release

Using this function with workflow dispatch, the RELEASE file can be rewritten.


## Usage

### Creating a personal access token

Creating your personal access token and save as a repository secret. (e.g. GH_PAT)

Personal access tokens (classic): https://github.com/settings/tokens

> **Note**
>
> Fine-grained personal access tokens are not supported currently.

### Workflow

```yaml
name: Release

on:
  workflow_dispatch:
    inputs:
      release_file_path:
        description: "path to RELEASE file"
        required: true
        type: choice
        options:
          - app1/RELEASE
          - app2/RELEASE
      next_semver_level:
        description: "semver level to bump"
        required: true
        type: choice
        options:
          - patch
          - minor
          - major

jobs:
  main:
    runs-on: ubuntu-22.04
    permissions:
      contents: read
    timeout-minutes: 5
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
      - name: release
        uses: peaceiris/gh-release-with-wf-dispatch@feat-semver-input
        with:
          github_token: ${{ secrets.GH_PAT }}
          release_file_path: ${{ github.event.inputs.release_file_path }}
          next_semver_level: ${{ github.event.inputs.next_semver_level }}
          base_branch: "develop"
          user_name: "username"
          user_email: "username@example.com"
```
