# Contributing to Keecon Proto.Actor Go

This repository is an independently maintained derivative of
[asynkron/protoactor-go](https://github.com/asynkron/protoactor-go), not the upstream Proto.Actor repository. All
contributions are licensed under Apache-2.0 and must preserve existing copyright and attribution notices.

## Choose the correct branch

- Target `main` for normal development, dependency updates, runtime updates, and security fixes.
- Do not target `origin/dev`. It is the first-stage intake branch for `upstream/dev` and is not a release line.
- Target `v0.2.x` or `v0.3.x` only when a maintainer has explicitly requested a backport.
- Keep a legacy backport limited to the requested fix. Adapt and test it against the target branch rather than merging
  `main` or upgrading unrelated dependencies.

New release tags start at `v0.5.0`. Until that tag exists, `main` remains an active development line.

The `main` branch supports the two most recent major Go releases and uses the older release as its minimum Go version.
When Go publishes a new major release, update the root module, independent example modules, and CI matrix together.
Legacy branches retain their own historical Go baselines.

## Development workflow

1. Start from the current target branch and create a focused topic branch.
2. Keep structural refactoring separate from behavior changes.
3. Add or update tests in proportion to the behavior and risk being changed.
4. Run the relevant local checks before requesting review.
5. Describe behavior, compatibility, security, and dependency impacts in the change request.

For changes integrated from upstream, use the repository's controlled flow:

```text
upstream/dev -> origin/dev -> staged conflict review -> main
```

Do not overwrite local module-path changes or local fixes merely to match upstream. Resolve and verify each conflict in
the integration stage.

## Required checks

For a typical root-module change, run:

```bash
go build ./...
make test-short
make vet
make lint
make vuln
```

Run targeted tests for every changed package. Changes to cluster, scheduler, remote, or persistence behavior may also
need the services in `docker-compose.yml`, including Consul, followed by the relevant full test targets. Example
directories are independent Go modules and must be built separately when they are affected.

Generated Protocol Buffer files must stay synchronized with their source `.proto` files. Avoid modifying generated
files by hand.

## Compatibility and dependencies

Changes to public APIs, message schemas, serialization, persistence formats, or remote protocols must document the
migration and compatibility impact. A passing compile is not sufficient evidence for protocol compatibility.

Keep new dependencies to a minimum. A dependency change must state why it is needed, its maintenance status, and its
license. Prefer permissive licenses compatible with Apache-2.0. Legacy backports must not contain unrelated dependency
refreshes.

## Commit and review scope

Use focused commits with an imperative summary that explains the change. Do not mix generated output, formatting,
unrelated cleanup, or broad renaming into a behavioral fix. Report unrelated dead code or lint findings separately.

A review should verify the target branch, tests, race and lifecycle implications, security impact, module path, and any
deviation from upstream behavior before integration.
