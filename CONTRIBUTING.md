## Contributing to KubeVirt Tekton Tasks

### Our workflow

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

### Getting started

To make yourself comfortable with the code, you might want to work on some
Issues marked with one or more of the following labels:
[good-first-issue](https://github.com/kubevirt/kubevirt-tekton-tasks/labels/good-first-issue),
[help wanted](https://github.com/kubevirt/kubevirt-tekton-tasks/labels/help%20wanted)
or [kind/bug](https://github.com/kubevirt/kubevirt-tekton-tasks/labels/kind%2Fbug).
Any help is highly appreciated.

### Testing

**Untested features do not exist**. To ensure that what the code really works,
relevant flows should be covered via unit tests and functional tests. So when
thinking about a contribution, also think about testability. All tests can be
run local without the need of CI. Have a look at the
[Testing](docs/getting-started.md#testing)
section in the [Developer Guide](docs/getting-started.md).

### Contributor compliance with Developer Certificate Of Origin (DCO)
We require every contributor to certify that they are legally permitted to contribute to our project. A contributor expresses this by consciously signing their commits, and by this act expressing that they comply with the Developer Certificate Of Origin

A signed commit is a commit where the commit message contains the following content:

Signed-off-by: John Doe <jdoe@example.org>
This can be done by adding --signoff to your git command line.


### PR Checklist
Before your PR can be merged it must meet the following criteria:

- README.md has been updated if new task is added or functionality of existing tasks is affected.
- [Onboarding a New Task](docs/onboarding-new-task.md) has been taken into account when introducing a new task
- Functionality must be fully tested.
