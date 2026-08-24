---
inclusion: fileMatch
fileMatchPattern: "*.md"
---

# Markdown style rules

This project enforces markdownlint via pre-commit. Config is in `.markdownlint.yaml`.

## Line length (MD013)

- Maximum line length: **120 characters** for prose.
- Code blocks and tables are **exempt** — do not wrap content inside fenced code blocks.
- Wrap long prose lines by breaking at a natural word boundary before column 120.

## Fenced code blocks (MD040)

- Every fenced code block **must** have a language identifier.
- Use `text` for generic output or directory listings when no specific language applies.
- Common identifiers used in this project: `sh`, `go`, `yaml`, `text`, `hcl`.

```sh
# correct
make build
```

```text
# also correct — generic/no language
redshift/   All provider Go source
```

## Other active rules

- **MD033** — inline HTML is allowed only for `<a>`, `<p>`, `<img>` tags.
- Standard markdownlint defaults apply for everything else (headings, blank lines, etc.).
