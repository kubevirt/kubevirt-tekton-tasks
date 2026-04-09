# Contributing to KubeVirt Tekton Tasks

## Our workflow

Contributing to KubeVirt Tekton Tasks should be as simple as possible. Have a question? Want
to discuss something? Want to contribute something? Just open an
[Issue](https://github.com/kubevirt/kubevirt-tekton-tasks/issues), a [Pull
Request](https://github.com/kubevirt/kubevirt-tekton-tasks/pulls), or send a mail to our
[Google Group](https://groups.google.com/forum/#!forum/kubevirt-dev).

If you spot a bug or want to change something pretty simple, just go
ahead and open an Issue and/or a Pull Request, including your changes
at [kubevirt/kubevirt-tekton-tasks](https://github.com/kubevirt/kubevirt-tekton-tasks).

For bigger changes, please create a tracker Issue, describing what you want to
do. Make sure that all your Pull Requests link back to the
relevant Issues.

## Getting started

To make yourself comfortable with the code, you might want to work on some
Issues marked with one or more of the following labels:
[good-first-issue](https://github.com/kubevirt/kubevirt-tekton-tasks/labels/good%20first%20issue),
[help wanted](https://github.com/kubevirt/kubevirt-tekton-tasks/labels/help%20wanted)
or [kind/bug](https://github.com/kubevirt/kubevirt-tekton-tasks/labels/kind%2Fbug).
Any help is highly appreciated.

## Testing

**Untested features do not exist**. To ensure that what the code really works,
relevant flows should be covered via unit tests and functional tests. So when
thinking about a contribution, also think about testability. All tests can be
run local without the need of CI. Have a look at the
[Testing](testing.md) documentation.

## Contributor compliance with Developer Certificate Of Origin (DCO)

We require every contributor to certify that they are legally permitted to contribute to our project. A contributor expresses this by consciously signing their commits, and by this act expressing that they comply with the Developer Certificate Of Origin.

A signed commit is a commit where the commit message contains the following content:

`Signed-off-by: John Doe <jdoe@example.org>`

This can be done by adding `--signoff` to your git command line.

## PR checklist

Before your PR can be merged it must meet the following criteria:

- README.md has been updated if new task is added or functionality of existing tasks is affected.
- [Adding a New Task](#adding-a-new-task) process has been followed when introducing a new task.
- Functionality must be fully tested.

## Code quality workflow

### After editing Go files

1. `make lint-fix` to auto-format
2. `make test` to run unit tests

### Before committing

1. `make lint` to verify formatting
2. `make test` to verify all tests pass
3. `make test-yaml-consistency` to verify generated YAML is up to date
4. All commits require DCO sign-off: `git commit --signoff`
5. If AI tools were used, add to commit message: `Co-authored-by: <AI tool name>`

## Adding a new task

An image stub must be committed and registered in CI before the real image can be tested.

1. Run `make onboard-new-task-with-ci-stub` and fill the name of the task and the name of its ENV variable.
    - A new module for this task will be created with a simple stub `Dockerfile`. You can modify the `Dockerfile` to include all the images required.
    - A new config will be created in `./configs`.
    - A new ENV variable will be registered in `scripts/common.sh`. The ENV variable will be used by the CI to deploy the tasks. You can modify these, for example if one image is used by multiple tasks.
2. Commit these changes and make a PR against https://github.com/kubevirt/kubevirt-tekton-tasks.
3. Create a PR against https://github.com/openshift/release and update `kubevirt-kubevirt-tekton-tasks-main.yaml`.
    - All new base images have to be added.
    - A new task image has to be added and must point to the stub `Dockerfile`.
    - The name of this image has to be passed in the already registered ENV variable.
4. Implement the task functionality and create a new PR against https://github.com/kubevirt/kubevirt-tekton-tasks.
5. The CI should build the task image before the tests run and then it should be ready to use.

---
<- Back to [AGENTS.md](../AGENTS.md) | [Documentation Index](../AGENTS.md#documentation)
