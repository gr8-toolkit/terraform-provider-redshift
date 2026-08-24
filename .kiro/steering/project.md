# terraform-provider-redshift

A Terraform provider for Amazon Redshift, built with the HashiCorp Terraform Plugin SDK v2.

## Language and runtime

- Go 1.26 (minimum). The project targets `go 1.26` in `go.mod` — golangci-lint v2.13.1 is built
  with Go 1.26 and cannot lint `go 1.27`+ modules yet.
- Module path: `github.com/gr8-toolkit/terraform-provider-redshift`

## Key dependencies

- `github.com/hashicorp/terraform-plugin-sdk/v2` — Terraform Plugin SDK v2 (context-aware CRUD APIs)
- `github.com/aws/aws-sdk-go-v2` — AWS SDK v2 for Redshift GetClusterCredentials
- `github.com/lib/pq` — PostgreSQL driver (Redshift speaks the Postgres wire protocol)

## Project layout

```text
redshift/          All provider Go source (resources, data sources, helpers)
docs/              Auto-generated Terraform documentation (do not edit manually)
examples/          HCL usage examples for docs generation
tools/             go:generate tooling (tfplugindocs)
```

## Build, test, and lint

```sh
make build         # go install — builds and installs the provider binary
make test          # go test ./redshift — unit tests (no live cluster needed)
make testacc       # TF_ACC=1 go test ./... — acceptance tests (requires real Redshift)
make fmt           # gofmt -w on all .go files
make vet           # go vet
make doc           # go generate — regenerates docs/ from schema annotations
```

Acceptance tests require these environment variables pointing at a real Redshift cluster:
`REDSHIFT_HOST`, `REDSHIFT_USER`, `REDSHIFT_PASSWORD`, `REDSHIFT_DATABASE`
(optional, defaults to `redshift`).

## Code conventions

### Terraform SDK patterns

- All CRUD operations use the context-aware fields:
  `CreateContext`, `ReadContext`, `UpdateContext`, `DeleteContext`.
- Importers use `StateContext: schema.ImportStatePassthroughContext`.
- The provider itself uses `ConfigureContextFunc`.
- The `RedshiftResourceFuncContext` wrapper in `helpers.go` bridges the domain function signature
  `func(*DBConnection, *schema.ResourceData) error` to the SDK's `diag.Diagnostics` return.
- `RedshiftResourceFunc` and `RedshiftResourceExistsFunc` are legacy wrappers kept for
  compatibility; prefer `RedshiftResourceFuncContext` for new code.

### Read function contract

Every `Read` / `ReadContext` function must handle `sql.ErrNoRows` by calling `d.SetId("")` and
returning `nil` — this signals Terraform that the resource no longer exists and needs to be
recreated.

### Error strings

Follow Go conventions: error strings are lowercase and do not end with punctuation.
Example: `fmt.Errorf("could not connect to cluster: %w", err)`.

### Type assertions

`d.Get(key).(type)` assertions against Terraform schema values are safe and intentional — the
schema guarantees the type. Do not add checked-assertion boilerplate for these. The
`forcetypeassert` linter is suppressed for the `redshift/` package for this reason.

### startTransaction

Use `startTransaction(db.client)` — no database parameter (it was always passed as `""`).

## Pre-commit hooks

Configured in `.pre-commit-config.yaml`. Run manually with:

```sh
pre-commit run --all-files
```

Hooks: `check-merge-conflict`, `check-added-large-files`, `detect-private-key`,
`check-case-conflict`, `mixed-line-ending`, `trailing-whitespace`, `end-of-file-fixer`,
`prettier` (YAML/Avro), `markdownlint` (Markdown outside `docs/` and `CHANGELOG.md`),
`golangci-lint-full`.

The golangci-lint hook uses `default_language_version: golang: "1.26.0"` to match the version
golangci-lint was built with.

## CI

GitHub Actions (`.github/workflows/ci.yml`) runs on every PR and push to `main`:

- **prek** — runs all pre-commit hooks via `j178/prek-action`
- **build** — `go build` + `golangci-lint-action` v2.13.1
- **generate** — `go generate` and checks for uncommitted diff
- **test** — acceptance tests against Terraform 1.8

## Linting

`golangci-lint` v2.13.1 via pre-commit and CI (`golangci-lint-action`). Config in `.golangci.yml`.

Active linters: `errcheck`, `staticcheck`, `unused`, `forcetypeassert`, `predeclared`,
`unparam`, `usetesting`.

The `forcetypeassert` linter is excluded for the `redshift/` path because Terraform SDK
`d.Get()` calls and `StateFunc` closures use type assertions that are guaranteed safe by the
schema definition.
