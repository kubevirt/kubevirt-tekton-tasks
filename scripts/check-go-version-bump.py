#!/usr/bin/env python3
"""Report Go directive changes between a PR merge base and head commit."""

from __future__ import annotations

import html
import re
import subprocess
import sys
from pathlib import PurePosixPath


DIRECTIVES = ("go", "toolchain")
DIRECTIVE_PATTERN = re.compile(r"^(go|toolchain)\s+(\S+)")


def git_output(*args: str) -> str:
    return subprocess.run(
        ["git", *args], check=True, capture_output=True, text=True
    ).stdout


def read_directives(revision: str, path: str) -> dict[str, str]:
    result = subprocess.run(
        ["git", "show", f"{revision}:{path}"],
        capture_output=True,
        text=True,
    )
    if result.returncode != 0:
        return {}

    directives: dict[str, str] = {}
    for raw_line in result.stdout.splitlines():
        match = DIRECTIVE_PATTERN.match(raw_line.strip())
        if match:
            name, version = match.groups()
            if name == "toolchain" and version.startswith("go"):
                version = version.removeprefix("go")
            directives[name] = version
    return directives


def main() -> int:
    if len(sys.argv) != 3:
        print(f"Usage: {sys.argv[0]} BASE_SHA HEAD_SHA", file=sys.stderr)
        return 2

    base_sha, head_sha = sys.argv[1:]
    for revision in (base_sha, head_sha):
        git_output("rev-parse", "--verify", f"{revision}^{{commit}}")

    merge_base = git_output("merge-base", base_sha, head_sha).strip()
    changed_paths = git_output(
        "diff", "--name-only", "--no-renames", merge_base, head_sha
    ).splitlines()
    go_mod_paths = sorted(
        path
        for path in changed_paths
        if PurePosixPath(path).name == "go.mod"
        and "vendor" not in PurePosixPath(path).parts
    )

    changes: list[tuple[str, str, str, str]] = []
    for path in go_mod_paths:
        base_directives = read_directives(merge_base, path)
        head_directives = read_directives(head_sha, path)
        for name in DIRECTIVES:
            before = base_directives.get(name)
            after = head_directives.get(name)
            if before != after:
                changes.append(
                    (path, name, before or "not set", after or "not set")
                )

    if not changes:
        return 0

    print("Renovate changed Go version directives in this PR:")
    print()
    print("| Module | Directive | Base | PR |")
    print("| --- | --- | --- | --- |")
    for path, name, before, after in changes:
        safe_path = html.escape(path).replace("|", "&#124;")
        print(
            f"| <code>{safe_path}</code> | `{name}` | `{before}` | `{after}` |"
        )
    print()
    print("Go version was bumped! Verify you can merge this PR!")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
