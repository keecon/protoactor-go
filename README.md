# Proto.Actor Go — Keecon maintained

[English](README.md) | [한국어](README.ko.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/keecon/protoactor-go.svg)](https://pkg.go.dev/github.com/keecon/protoactor-go)
[![checks](https://github.com/keecon/protoactor-go/actions/workflows/checks.yml/badge.svg?branch=main)](https://github.com/keecon/protoactor-go/actions/workflows/checks.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/keecon/protoactor-go)](https://goreportcard.com/report/github.com/keecon/protoactor-go)

This repository is an independently maintained derivative of
[asynkron/protoactor-go](https://github.com/asynkron/protoactor-go). It preserves the original Git history and
Apache-2.0 license while maintaining legacy compatibility branches and an updated `main` development line. This
repository is not the upstream Proto.Actor repository.

The module provides an actor runtime for Go together with routing, scheduling, streams, persistence, remoting, and
clustering packages. Serialization and cross-process protocols remain explicit concerns, using Protocol Buffers and
gRPC where applicable.

## Support and release policy

| Version line | Branch | Status | Maintenance policy |
| --- | --- | --- | --- |
| v0.5 and later | `main` | Active development | Current Go runtime, dependencies, security fixes, and upstream integration |
| v0.3 | `v0.3.x` | Legacy compatibility | Backports only after an explicit request |
| v0.2 | `v0.2.x` | Legacy compatibility | Backports only after an explicit request |
| No releases | `origin/dev` | Upstream intake | Internal staging only; do not build releases from this branch |

`main` is the only actively supported line. New release tags start at `v0.5.0`. Until a `v0.5.0` tag exists, `main`
is a development line and consumers should pin an audited commit instead of tracking the branch tip.

The `v0.2.x` and `v0.3.x` branches exist for projects that cannot yet migrate. They do not receive routine runtime,
dependency, feature, or security updates. A requested backport is reviewed against the target branch, adapted when
necessary, and verified on that branch before it is accepted. Compatibility is evaluated against that branch's own
API, dependencies, and Go toolchain—not against `main`.

## Upstream integration

Upstream changes follow this controlled path:

```text
asynkron/protoactor-go dev
          ↓
      origin/dev
          ↓
staged integration and conflict review
          ↓
         main
```

`origin/dev` is the first reception point for `upstream/dev`. Changes are integrated into `main` only after conflicts,
local behavior, tests, security impact, and module-path differences have been reviewed. Local changes are not pushed
back to upstream automatically.

## Installation

The module path is `github.com/keecon/protoactor-go`.

For an active-development build, select and pin a reviewed `main` commit:

```bash
go get github.com/keecon/protoactor-go@<commit-sha>
```

For a legacy release, use the required release tag or an explicitly reviewed branch revision. For example:

```bash
go get github.com/keecon/protoactor-go@v0.3.0
```

Do not replace a legacy version with `main` without testing the application's actor lifecycle, message contracts,
persistence, and any remote or cluster protocols it uses.

## Go versions

The active line is validated with the currently supported Go release lines configured in
[the checks workflow](.github/workflows/checks.yml). The `go` directive in `go.mod` is the module's language/toolchain
floor; it is not a promise that an otherwise unsupported Go release receives maintenance here.

Legacy branches retain their historical Go and dependency baselines. Their own workflow files are the source of truth
for backport validation.

## Hello actor

```go
package main

import (
	"fmt"

	"github.com/keecon/protoactor-go/actor"
)

type hello struct{ who string }

type helloActor struct{ done chan<- struct{} }

func (a *helloActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *hello:
		fmt.Printf("Hello %s\n", msg.who)
		a.done <- struct{}{}
	}
}

func main() {
	system := actor.NewActorSystem()
	done := make(chan struct{})
	props := actor.PropsFromProducer(func() actor.Actor {
		return &helloActor{done: done}
	})

	pid := system.Root.Spawn(props)
	system.Root.Send(pid, &hello{who: "Proto.Actor"})
	<-done
	system.Root.Stop(pid)
}
```

More runnable examples are available under [`examples`](examples).

## Build and test

Build the root module:

```bash
go build ./...
```

Run the short test suite and static checks:

```bash
make test-short
make vet
make lint
```

Some cluster, scheduler, remote, and persistence integration tests require the services declared in
[`docker-compose.yml`](docker-compose.yml), including a local Consul agent. Run the required services before using the
full test targets:

```bash
docker compose up -d
make test2
```

Example directories contain independent Go modules. The CI workflow builds them separately.

## Contributing

Development changes target `main`. Do not open a legacy-branch change unless a backport was explicitly requested.
Read [CONTRIBUTING.md](CONTRIBUTING.md) for the branch, test, dependency, and upstream-sync rules.

## License and attribution

This project is licensed under the [Apache License 2.0](LICENSE). The original Git history, copyright notices, and
license attribution from `asynkron/protoactor-go` are retained. Contributions made in this repository remain subject
to the same license.
