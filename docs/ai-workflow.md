# AI Agent Workflow Rules

<rules>
- You MUST follow the 2-step process (planning -> approval -> implementation) for ALL file modifications (code, docs, config, makefiles, scripts, YAML).
- NEVER hand-edit files under `release/tasks/` or `release/pipelines/`. Always modify templates and regenerate.
- NEVER skip the Review Checklist before committing.
- ALWAYS use `git commit --signoff` (DCO requirement).
- NEVER deviate from the approved plan without asking for permission first.
</rules>

<context>
**Exception for trivial changes**: For single-line fixes, typo corrections, or formatting-only changes, you may skip the `change.md` planning step if the user explicitly confirms the change is trivial. Still follow the Review Checklist and Git Workflow steps.
</context>

## STEP 1: Planning and Documentation

<instructions>
1. Before writing or modifying any project code, create or update a file named `change.md` in the root directory.
2. In `change.md`, clearly outline:
   - The goal of the task.
   - The files that will be modified or created.
   - A step-by-step technical plan of the exact changes you intend to make.
   - **Note:** `change.md` is gitignored and temporary - it's only used during the current task session for planning purposes.
3. Once you have saved `change.md`, you must **STOP**. Do not write any code. Ask the user: *"I have outlined the plan in change.md. Please review it. Reply with 'APPROVED' to proceed, or let me know what needs to be changed."*
4. Wait for the user's explicit approval.
</instructions>

## STEP 2: Implementation

<instructions>
1. ONLY after the user explicitly types "APPROVED" (or a clear affirmative), begin implementing the steps exactly as outlined in `change.md`.
2. Do not deviate from the approved plan without asking for permission first.
3. Before completing implementation, follow the Review Checklist below.
</instructions>

### Review Checklist

<checklist>
This checklist MUST be followed:
- When implementing changes (after coding, before committing)
- When performing standalone code reviews (when asked to review code without implementing)

#### Documentation Review

Check if changes require updates to `docs/` files:
- New features, tasks, or commands
- Modified behavior or configuration
- New dependencies or build requirements
- Changed workflows or processes
- Update relevant documentation files before proceeding
- If unsure whether docs need updating, ask the user

#### Manifest Regeneration Review

Check if changes require regenerating YAML manifests:
- **NEVER hand-edit files in `release/`** (see [YAML Generation](yaml-generation.md))
- Run `make generate-yaml-tasks` if changes affect:
  - Files in `templates/` directory
  - Files in `configs/` (Ansible variable files)
  - Task parameters or behavior
  - RBAC requirements in task definitions
- Run `make generate-pipelines` if changes affect:
  - Files in `templates-pipelines/` directory
  - Pipeline definitions or structure
- Commit any regenerated YAML files along with code changes

#### Code Quality Review

Verify code meets project standards:
- Run `make test` to verify unit tests still pass
- Address any test failures before proceeding
</checklist>

### After Implementation

<instructions>
Once implementation and the review checklist are complete, **ask the user** if they would like to proceed with committing and pushing the changes to GitHub (STEP 3). Do not automatically proceed to STEP 3 without user confirmation.
</instructions>

## STEP 3: Git Workflow

<instructions>
After successfully implementing the changes:

1. **Create or verify your feature branch**:
   - If you're already on a feature branch for this work, stay on it
   - Otherwise, create a new branch with a descriptive name:
   ```bash
   git checkout -b feature/descriptive-name
   # or
   git checkout -b fix/issue-description
   ```

2. **Commit changes** following the commit message guidelines:
   - **Format**: `type(scope): description`
   - **Types**: `feat`, `fix`, `docs`, `test`, `build`, `refactor`
   - **Subject line**: Imperative mood, max 72 chars, capitalized, no period
   - **Body**: Wrap at 72 chars, explain what and why
   - **DCO**: Always use `--signoff`

   ```bash
   git add <files>
   git commit --signoff -m "feat(execute-in-vm): Add support for custom SSH timeout

   Implements configurable SSH timeout parameter to handle slow VM boots.
   Default timeout remains 5 minutes for backwards compatibility."
   ```

3. **Push to origin**:
   ```bash
   git push origin feature/descriptive-name
   ```

4. **Inform the user** that the branch has been pushed and provide the branch name for PR creation.
</instructions>

## Jira Task Management

<context>
See [AGENTS.md - Jira Task Management](../AGENTS.md#jira-task-management) for instructions on fetching Jira task information.

When working on Jira-related changes, include the task ID and relevant information in `change.md` and commit messages.
</context>

## Why This Process?

<context>
- **Prevents mistakes**: Plan is reviewed before implementation
- **Saves time**: Avoids implementing wrong solutions
- **Builds trust**: User has full visibility and control
- **Improves quality**: Thoughtful planning leads to better code
</context>

## Example Workflow

<example>
```
User: "Add a new Tekton task for disk snapshot"

AI: [Creates change.md with detailed plan]
AI: "I have outlined the plan in change.md. Please review it.
     Reply with 'APPROVED' to proceed, or let me know what needs
     to be changed."

User: [Reviews change.md, suggests modifications]

AI: [Updates change.md based on feedback]

User: "APPROVED"

AI: [Implements the plan step by step]
AI: [Creates branch, commits with DCO sign-off, pushes to origin]
AI: "Changes have been pushed to branch 'feature/disk-snapshot-task'.
     You can now create a PR on GitHub."
```
</example>

---
<- Back to [AGENTS.md](../AGENTS.md) | [Documentation Index](../AGENTS.md#documentation)
