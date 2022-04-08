# gh-release-with-wf-dispatch

gh-release-with-wf-dispatch crates github tag with gh-release

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
        options:
          - RELEASE
          - testdata/RELEASE
      newTag:
        description: 'new tag'
        required: 'false'
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
    "github_token": ${{ secrets.GITHUB_TOKEN }}
    "owner": ${{ github.event.repository.owner.login }} 
    "repo": ${{ github.event.repository.name }}
    "releae_file_path": ${{ github.event.inputs.releaeFilePath }}
    "base_branch": ${{ github.event.inputs.baseBranch }}
    "new_tag": ${{ github.event.inputs.newTag }}
```
