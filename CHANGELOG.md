# Changelog

This file records user-visible changes to the actively maintained release line. Legacy `v0.2.x` and `v0.3.x`
branches keep their own historical compatibility baselines and receive changes only through explicitly requested
backports.

## [Unreleased]

### Added

- Added Go vulnerability scanning for both supported Go release lines on pushes, pull requests, and a weekly schedule.
- Documented the controlled `upstream/dev` to `origin/dev` to `main` integration path and the legacy branch policy.

### Changed

- Established `github.com/keecon/protoactor-go` as the maintained module identity for the `main` release line.
- Raised the minimum Go version to Go 1.26 and validated the module with Go 1.26 and Go 1.27.
- Pinned local test, lint, and vulnerability tools to reproducible versions.
- Updated gRPC, OpenTelemetry, Consul, etcd, Kubernetes, and supporting dependencies to current compatible releases.
- Synchronized independent example modules with the root module's Go and dependency baselines.

### Fixed

- Fixed actor, router, remote endpoint, cluster topology, and cluster provider lifecycle races found during integration
  and race testing.
- Made the virtual-actor cluster identity test wait for its probe to start before using the probe context.

### Compatibility

- `main` is the active development line for v0.5 and later. Consumers should pin a reviewed commit until `v0.5.0`
  is released.
- `v0.2.x` and `v0.3.x` remain available for legacy consumers but do not receive routine runtime, dependency, or
  security updates.

[Unreleased]: https://github.com/keecon/protoactor-go/compare/v0.3.0...HEAD
