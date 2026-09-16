# CV Generator

Generate professional CVs as PDF from a JSON file.

## Installation

Download a prebuilt binary for your platform:

```bash
curl -sSL https://raw.githubusercontent.com/othmaneBakkass/cv_gen/main/scripts/install.sh | bash -s -- latest
```

Pass a specific release tag instead of `latest` to pin a version (e.g. `bash -s -- v1.0.0`).
Set `BIN_DIR` to change the install location.

Or build from source (Go 1.27+):

```bash
go install github.com/othmaneBakkass/cv_gen@latest
# or, from a clone:
make build   # produces ./_build/cv_gen
```

## Usage

```bash
cv_gen generate -i examples/cv.json -o ./output
```

One PDF is written per entry in the `data` array, named after that entry's
`fileName` (a `.pdf` extension is added automatically).

### Flags

- `-i, --input`: path to the JSON data file (required)
- `-o, --output`: output directory (default: current directory, created if missing)

## Input format

See [`schemas/v1.json`](schemas/v1.json) for the full JSON Schema and
[`examples/cv.json`](examples/cv.json) for a complete example.

```json
{
  "data": [
    {
      "template": "t1",
      "fileName": "jane-doe",
      "head": { "fullName": "Jane Doe", "address": "…", "phone": "…", "email": "jane@example.com" },
      "education": [ { "school": "…", "location": "…", "startedAt": "…", "endedAt": "…", "degree": "…", "description": "…" } ],
      "jobs": [ { "company": "…", "location": "…", "position": "…", "startedAt": "…", "endedAt": "…", "tools": ["…"], "highlights": ["…"] } ],
      "languages": [ { "language": "english", "level": "fluent" } ]
    }
  ]
}
```

All fields are required and every list must contain at least one item.
`email` must be a valid email address.

### Templates

| Name | Description |
|------|-------------|
| `t1` | Centered contact header, then Education / Experience / Languages |

## Development

```bash
make fmt    # gofmt
make vet    # go vet ./...
make build  # build the CLI
make tidy   # go mod tidy
```
