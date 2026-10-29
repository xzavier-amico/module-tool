# module-tool

A converter and generator for network client modules.

Supported formats:
- surge
- loon
- qx (Quantumult X)
- egern
- stash

## Features
- Convert a module between any of the supported formats (CLI / Web UI / HTTP API)
- Batch generation from a YAML config file (`src/` for sources, `dist/` for outputs)
- Metadata overrides (name, desc, author, icon)

## Build
```bash
go build -o moduletool ./cmd/moduletool
```

## Usage
```bash
# Convert
./moduletool convert --from surge --to loon --input in.sgmodule --output out.plugin

# Batch generate
./moduletool generate --config examples/config.yaml --base .

# Web UI
./moduletool serve --port 8080
```

## Config example
```yaml
src_dir: src
dist_dir: dist
jobs:
  - name: my-module
    source: https://example.com/module.sgmodule
    from: surge
    targets: [loon, egern, stash, qx]
    icon: shield.fill
    desc: Demo module
    author: me
```
EOF