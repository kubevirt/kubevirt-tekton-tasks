# Creating and Onboarding a New Task
Firstly, an image stub has to be committed and registered in the CI,
so the real image can be tested by it.

1. Run `make onboard-new-task-with-ci-stub` and fill the name of the task and the name of its ENV variable.
    - A new module for this task will be created with a simple stub `Dockerfile`.
      You can modify the `Dockerfile` to include all the images required.
    - A new config will be created in `./configs`
    - A new ENV variable will be registered in `scripts/common.sh`.
      The ENV variable will be used by the CI to deploy the tasks.
      You can modify these, for example if one image is used by multiple tasks.
2. Commit these changes and make a PR against https://github.com/kubevirt/kubevirt-tekton-tasks.
3. Create a PR against https://github.com/openshift/release and update `kubevirt-kubevirt-tekton-tasks-main.yaml`.
    - All new base images have to be added.
    - A new task image has to be added & pointing to the stub `Dockerfile`.
    - The name of this image has to be passed in the already registered ENV variable.
4. Implement the task functionality and create a new PR against https://github.com/kubevirt/kubevirt-tekton-tasks.
5. The CI should build the task image before the tests run and then it should be ready to use.
