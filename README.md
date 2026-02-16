# plugin-morphe-morpherepo

Generates MorpheRepo data access contract definitions (`.repo` YAML files) from Morphe model definitions (`KA:MO1:YAML1`).

## What it generates

For each Morphe model, this plugin produces a `.repo` file describing the repository contract:

- **Repository name** -- `{ModelName}Repository`
- **Identifiers** -- primary and secondary, extracted from model identifiers
- **Filters** -- derived from `ForOne` and `ForOnePoly` relationships
- **CRUD operations** -- `list`, `get`, `create`, `update`, `delete`

### Example

Given a `Project` model with a `ForOne` relationship to `Organization`:

```yaml
# project.repo
name: ProjectRepository
model: Project

identifiers:
  primary:
    fields:
      - name: ID
        type: UUID
  code:
    fields:
      - name: Code
        type: String

filters:
  - name: organizationID
    type: UUID
    relation: Organization

operations:
  list: true
  get: true
  create: true
  update: true
  delete: true
```

### Relationship handling

| Relationship type | Effect                                  |
|-------------------|-----------------------------------------|
| `ForOne`          | Generates a filter (`{relation}ID`)     |
| `ForOnePoly`      | Generates a filter (uses `Through` field if present) |
| `HasOne`          | No filter generated                     |
| `HasMany`         | No filter generated                     |

## Input / output

| Direction | Format           | Store suggestion | Description                        |
|-----------|------------------|------------------|------------------------------------|
| Input     | `KA:MO1:YAML1`  | `KA_MO_YAML`    | Morphe registry (models, enums, structures, entities) |
| Output    | `KA:MR1:YAML1`  | `KA_RE_YAML`    | MorpheRepo `.repo` YAML files      |

Output files are named in snake_case: `Organization` becomes `organization.repo`, `UserProfile` becomes `user_profile.repo`.

## Configuration

| Key              | Type     | Default | Description                                  |
|------------------|----------|---------|----------------------------------------------|
| `excludeModels`  | string[] | `[]`    | List of model names to skip during generation |

## Pipeline context

This plugin is typically used in a **compile pipeline** that transforms Morphe schemas
into intermediate contract definitions, which are then consumed by language-specific
code generators:

```yaml
stores:
  KA_MO_YAML:
    format: "KA:MO1:YAML1"
    type: "localFileSystem"
    options:
      path: "./morphe"

  KA_RE_YAML:
    format: "KA:MR1:YAML1"
    type: "localFileSystem"
    options:
      path: "./morphe/repo"

plugins:
  "@kalo-build/plugin-morphe-morpherepo":
    version: "v1.0.0"
    inputs:
      morphe:
        format: "KA:MO1:YAML1"
        store: "KA_MO_YAML"
    output:
      format: "KA:MR1:YAML1"
      store: "KA_RE_YAML"

pipelines:
  generate:
    stages:
      - name: "morphe-repo"
        steps:
          - "plugin: @kalo-build/plugin-morphe-morpherepo"
```

## Project structure

```
plugin-morphe-morpherepo/
├── cmd/plugin/             # WASM entry point
├── pkg/compile/
│   ├── compile.go          # Main pipeline (MorpheToMorpheRepo)
│   ├── model_info.go       # Model analysis (identifiers, filters)
│   ├── generate_repo.go    # .repo YAML generation
│   └── cfg/                # CompileConfig definition
├── internal/testutils/     # Test helpers and ground-truth regeneration
├── testdata/
│   ├── registry/minimal/   # Sample Morphe registry input
│   └── ground-truth/       # Expected .repo output for integration tests
└── plugin.yaml             # Kalo plugin manifest
```

## Building

```bash
# Native binary
go build ./cmd/plugin

# WASM (for Kalo CLI)
GOOS=wasip1 GOARCH=wasm go build -o dist/plugin.wasm cmd/plugin/main.go
```

## Testing

```bash
go test ./...
```
