# Copilot instructions

## Build and test

- The project uses Go 1.23 (`go.mod`).
- Build the CLI with `go build -o image-tool .`.
- Run package tests with `go test ./...`. To run one test, use `go test ./ocrequest -run '^TestName$'`.
- The repository has no configured linter. Use `gofmt` on changed Go files.
- Runtime cluster configuration is read from `clusterconfig.json` next to the executable, before cluster flag validation. Generate it with `build/scripts/generate_config.sh <domain>` using authenticated cluster tools, or provide it separately. It contains credentials and must not be committed.

## Architecture

- `main.go` calls `ocrequest.Init()` before dispatching to command-line mode (`CmdlineMode`) or the HTTP server (`StartServer`).
- `ocrequest` owns the whole application flow: `userinterface.go` parses flags into the shared `CmdParams`; `init.go` loads family/namespace/cluster configuration, initializes credentials and family image-stream names; collection functions query OpenShift resources and use worker pools for concurrent namespace/cluster work.
- `CmdlineMode` coordinates collection of existing image streams/tags/images and tags referenced by running workloads, then joins and filters those results before writing JSON, YAML, CSV, table, or HTML output. The same package also implements deletion-filtering and shell-command output; check `README.md` before changing the CLI contract.
- Configuration is either generated dynamically from the config-tool data (`readConfigBitbucket.go`) or supplied by the static `FamilyNamespacesStat` map (`config.go`, selected with `-statcfg`). Runtime cluster credentials are separate from that family/namespace mapping.
- `webservice.go` exposes HTTP query and tag-use endpoints and reuses the same `CmdParams` and collection pipeline as the CLI. Keep changes to result semantics consistent across both entry points.

## Codebase conventions

- Domain types use package-level `T_`-prefixed names and typed string keys (`T_familyName`, `T_clName`, `T_nsName`, and image/tag types). Result data is represented by nested maps keyed by family, cluster, namespace, image stream, and tag; preserve these key relationships when changing collection or filtering.
- Collection code is split into resource-specific functions and worker-pool files (`worker*.go`). Keep OpenShift retrieval separate from aggregation/filtering and output formatting.
- Text/CSV/HTML renderers derive columns from result types' `Names()` and `Values()` methods in `types.go`; update those methods and the relevant renderer when adding or reordering displayed fields.
- The `-tagname`, `-isname`, `-istagname`, and namespace filters are regular expressions, not literal string matches. Preserve regex semantics at CLI and HTTP entry points.
- Generated API documentation lives under `docs/`; the build script invokes `swag i` before release builds.
