# AGENTS.md

## What this repo is

This is the **API contract repo** for `ecommerce-category-query-service` — the CQRS read side
for categories. It contains protobuf definitions, generated clients, and a thin hand-written Fx
client helper. The service that implements these RPCs and the query service that projects category
events into a read model live in separate repos.

The `Category`/`CategoryAttribute` messages are **denormalized** — attributes carry their full
details (name, slug, options, role, flags) copied from catalog-service, not just IDs. This is
intentional for a read model: consumers get everything in one call without joining back to the
write side.

## Source of truth vs. generated code

- **Edit only `proto/category_query/v1/category_query.proto`.** Everything under `gen/` is
  generated and committed to git (Go clients + TypeScript clients). Never hand-edit `gen/` —
  change the `.proto` and run `make generate`.
- Generated TS **source** is committed; only `gen/typescript/{node_modules,dist,package-lock.json}`
  are gitignored (see `.gitignore`).

## Common commands

```bash
make generate            # lint proto + generate TS client + generate Go client (the main command)
make lint                # buf lint only
make format              # buf format -w
make clean               # remove all generated files
make connect-breaking    # check for breaking proto changes against the 'master' branch
make tidy                # go mod tidy
make update-proto-deps   # buf dep update (refresh buf.lock)
make help                # list all targets grouped by category
```

First-time setup installs the proto toolchain (buf, protoc-gen-go, protoc-gen-connect-go,
protoc-gen-go-grpc) at pinned versions:

```bash
make connect-install-tools
```

TS generation uses **buf remote plugins** (no local install needed); Go generation uses the
**locally installed** plugins above.

## Code generation pipeline

`make generate` (defined in this repo's root `Makefile`, with logic split into `makefiles/`):

- `makefiles/protobuf-connect.mk` — Go: `buf generate --template buf.gen.yaml` emits
  `protoc-gen-go` (messages), `protoc-gen-connect-go` (Connect handlers/clients), and
  `protoc-gen-go-grpc` (native gRPC) into `gen/go/`. The Go module path is this repo itself, so
  consumers import `gen/go/category_query/v1`.
- `makefiles/connect-ts.mk` — TypeScript: `buf generate --template buf.gen.ts.yaml` emits
  `protoc-gen-es` into `gen/typescript/`, then **synthesizes** `package.json`, `tsconfig.json`,
  and `index.ts` (barrel of `export * from './*_pb.js'`). Package name is
  `@sokol111/ecommerce-category-query-service-api`, version taken from the `VERSION` file,
  published to GitHub Packages npm registry.
- The `_PROTO_SUBDIR` in `protobuf-connect.mk` is **auto-detected** as the first dir under
  `proto/` (here `category_query`). Directory layout must stay consistent for paths to resolve.

## Releasing (version → tag → consumers bump)

Releases are driven by the **`VERSION` file**, not by manual tags. Pushing a change to `VERSION`
on the `master` branch triggers `.github/workflows/release.yml`, which delegates to the reusable
`Sokol111/ecommerce-infrastructure/.github/workflows/api-release.yml`. That workflow tags the Go
module and publishes the TS package. Downstream Go services then bump their `go.mod` to the new
version (**release-then-bump** — the api must be released before a consumer can pin it).

## Consuming the generated Go client (`pkg/fxconfig`)

`pkg/fxconfig/grpc.go` is the one piece of hand-written Go here:
`fxconfig.NewGrpcClientsModule()` wires a native gRPC client for `CategoryQueryService`. It reads
config from koanf under key
`category-query.grpc`, builds a connection from `ecommerce-commons/pkg/http/grpc/client`, and
provides `categoryv1.NewCategoryQueryServiceClient`. A consuming service composes this module in
its `main.go` rather than dialing the connection itself.
