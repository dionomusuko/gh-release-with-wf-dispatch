# gh-release-with-wf-dispatch

gh-release-with-wf-dispatch crates github tag with gh-release

Using this function with workflow dispatch, the RELEASE file can be rewritten.


## Usage
1. Setting workflow dispatch
e.g.
```yaml
on:
  workflow_dispatch:
    inputs:
      releaseFilePath:
        description: 'Define the path of RELEASE file'
        required: 'true'
        default: 'RELEASE'
        type: choice
        options: # Set RELEASE file path
          - RELEASE
          - testdata/RELEASE
      newTag:
        description: 'new tag'
        required: 'false' # If you want to do PATCH version release, no input is required
      baseBranch:
        description: 'base branch e.g. master'
        require: 'true'
        default: 'master'
```

2. Write your workflow file
```yaml
- name: release
  uses: dionomusko/gh-release-with-wf-dispatch@master
  with:
    "github_token": ${{ secrets.GH_PAT }} # Set your GitHub personal access token (see: https://docs.github.com/ja/enterprise-cloud@latest/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
    "owner": ${{ github.event.repository.owner.login }}
    "repo": ${{ github.event.repository.name }}
    "releae_file_path": ${{ github.event.inputs.releaeFilePath }}
    "base_branch": ${{ github.event.inputs.baseBranch }}
    "new_tag": ${{ github.event.inputs.newTag }}
```
