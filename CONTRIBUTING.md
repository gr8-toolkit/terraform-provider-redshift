# Contributing

Contributions are welcome. Please read this guide before opening a pull request.

## Requirements

- [Go](https://golang.org/doc/install) >= 1.26
- [prek](https://github.com/j178/prek) — hook runner (`cargo install prek` or download from
  releases)
- [golangci-lint](https://golangci-lint.run/welcome/install/) >= v2.13.1

## Getting started

```sh
git clone git@github.com:gr8-toolkit/terraform-provider-redshift.git
cd terraform-provider-redshift
prek install
```

`prek install` sets up the Git hooks defined in `.pre-commit-config.yaml`. They run automatically
on every `git commit` and enforce formatting, linting, and markdownlint rules.

## Development workflow

| Command | What it does |
|---|---|
| `make build` | Format, then `go install` the provider binary |
| `make test` | Format, vet, then run unit tests (no cluster needed) |
| `make testacc` | Run acceptance tests against a real Redshift cluster |
| `make fmt` | Run `gofmt -w` on all `.go` files |
| `make vet` | Run `go vet` |
| `make doc` | Regenerate `docs/` from schema and `examples/` |
| `prek run --all-files` | Run all hooks against every file |

## Running acceptance tests

Acceptance tests require a real AWS Redshift cluster (provisioned or Serverless).

```sh
export REDSHIFT_HOST=<cluster-endpoint>
export REDSHIFT_USER=<superuser>
export REDSHIFT_PASSWORD=<password>
export REDSHIFT_DATABASE=dev   # optional, defaults to "redshift"
export REDSHIFT_PORT=5439       # optional, defaults to 5439

make testacc
```

If the cluster is only reachable from inside a VPC, connect via a SOCKS proxy:

```sh
export ALL_PROXY=socks5h://<host>:<port>
export NO_PROXY=127.0.0.1,localhost
```

### Optional feature flags

Some tests are skipped unless these variables are set:

| Variable | Feature |
|---|---|
| `REDSHIFT_DATASHARE_SUPPORTED` | Datashare resources (requires RA3 instance) |
| `REDSHIFT_DATASHARE_CONSUMER_NAMESPACE` | Datashare privilege — same-account namespace |
| `REDSHIFT_DATASHARE_CONSUMER_ACCOUNT` | Datashare privilege — cross-account |
| `REDSHIFT_TEMPORARY_CREDENTIALS_CLUSTER_IDENTIFIER` | Temporary credentials via `GetClusterCredentials` |
| `REDSHIFT_TEMPORARY_CREDENTIALS_ASSUME_ROLE_ARN` | Cross-account assume-role credentials |
| `REDSHIFT_EXTERNAL_SCHEMA_*` | External schema variants (Data Catalog, Hive, RDS, Redshift) |

## Code conventions

- All resource CRUD operations use `CreateContext` / `ReadContext` / `UpdateContext` /
  `DeleteContext` — not the deprecated non-context variants.
- Every `Read` function must handle `sql.ErrNoRows` by calling `d.SetId("")` and returning `nil`.
- Error strings are lowercase and do not end with punctuation:
  `fmt.Errorf("could not connect: %w", err)`.
- `d.Get(key).(type)` assertions on Terraform schema values are intentional — the schema
  guarantees the type. The `forcetypeassert` linter is suppressed for the `redshift/` package.
- Use `startTransaction(db.client)` — the database parameter was removed because it was always
  passed as `""`.

## Pull request checklist

- [ ] `prek run --all-files` passes with no errors
- [ ] New resources or attributes have a `Description` field in the schema
- [ ] Acceptance tests are updated or added for the changed behaviour
- [ ] `make doc` has been run if the schema changed (commit the updated `docs/`)
- [ ] The PR title is concise and describes the change (not "fix" or "update")

## Project layout

```text
redshift/          Provider source — resources, data sources, helpers, tests
docs/              Auto-generated documentation (do not edit manually)
examples/          HCL examples used by tfplugindocs
tools/             go:generate entry point for tfplugindocs
.goreleaser.yml    GoReleaser config — cross-platform builds + GPG signing
```
