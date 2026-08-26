# Terraform Provider for AWS Redshift

[![CI](https://github.com/gr8-toolkit/terraform-provider-redshift/actions/workflows/ci.yml/badge.svg)](https://github.com/gr8-toolkit/terraform-provider-redshift/actions/workflows/ci.yml)
[![prek](https://img.shields.io/badge/prek-enabled-brightgreen)](https://github.com/j178/prek)
[![Terraform Registry](https://img.shields.io/badge/terraform-registry-623CE4?logo=terraform)](https://registry.terraform.io/providers/gr8-toolkit/redshift/latest/docs)
[![OpenTofu Registry](https://img.shields.io/badge/opentofu-registry-FFDA18?logo=opentofu&logoColor=000)](https://registry.opentofu.org/providers/gr8-toolkit/redshift)

A Terraform provider for managing [AWS Redshift](https://aws.amazon.com/redshift/) objects:
users, groups, schemas, databases, grants, default privileges, and datashares.

Published on the
[Terraform Registry](https://registry.terraform.io/providers/gr8-toolkit/redshift/latest/docs)
and the
[OpenTofu Registry](https://registry.opentofu.org/providers/gr8-toolkit/redshift).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.8 or
  [OpenTofu](https://opentofu.org/docs/intro/install/) >= 1.6
- [Go](https://golang.org/doc/install) >= 1.26 (to build the provider from source)

## Usage

```hcl
terraform {
  required_providers {
    redshift = {
      source  = "gr8-toolkit/redshift"
      version = "~> 1.0"
    }
  }
}

provider "redshift" {
  host     = "my-cluster.us-east-1.redshift.amazonaws.com"
  port     = 5439
  username = "admin"
  password = var.redshift_password
  database = "dev"
}
```

See the [registry documentation](https://registry.terraform.io/providers/gr8-toolkit/redshift/latest/docs)
for the full list of resources and data sources.

## Building from Source

```sh
git clone git@github.com:gr8-toolkit/terraform-provider-redshift.git
cd terraform-provider-redshift
make build
```

`make build` runs `gofmt`, then `go install`, placing the binary in `$GOPATH/bin`.

## Documentation

Docs are generated from schema `Description` fields and the
[examples/](./examples) directory using
[tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs).
Files under `docs/` are auto-generated — do not edit them manually.

Regenerate after changing the schema or examples:

```sh
make doc
```

## Releasing

Releases are fully automated via GitHub Actions and
[GoReleaser](https://github.com/goreleaser/goreleaser/).

1. Update the changelog:

   ```sh
   RELEASE_VERSION=v1.x.y \
   CHANGELOG_GITHUB_TOKEN=<token> \
   make changelog
   ```

2. Push the changelog commit:

   ```sh
   git push origin master
   ```

3. Tag and push — this triggers the release workflow automatically:

   ```sh
   RELEASE_VERSION=v1.x.y make release
   ```

The [Release workflow](https://github.com/gr8-toolkit/terraform-provider-redshift/actions/workflows/release.yml)
builds cross-platform binaries, signs the checksums with GPG, and publishes
the GitHub release.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md).
