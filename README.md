# ordo

> A generic, format-pluggable configuration loader for Go.

## Usage

Formats are explicit: pick the built-in presets you need and hand them to the
loader.

```go
loader := &ordo.Loader[Config]{
    Formats: []ordo.Format{ordo.JSON, ordo.XML},
    Defaults: func(cfg *Config) {
        if cfg.Address == "" {
            cfg.Address = ":3000"
        }
    },
    Validate: func(cfg *Config) error {
        if cfg.Address == "" {
            return errors.New("address is required")
        }
        return nil
    },
    OnLoaded: func(path string) { log.Printf("config loaded from %s", path) },
}

cfg, path, err := loader.Load("./config.json")

err = loader.Save(cfg, "./config.json")
```

If the path has no matching extension, the given formats are probed in
priority order until one decodes; a bare `./config` also resolves by trying
each format's extensions.

Primitives, without orchestration:

```go
cfg, err := ordo.Load[Config](path, ordo.JSON)
err = ordo.Save(cfg, "./config.json", ordo.JSON)
```

## Built-in formats

| Preset    | Format | Extensions | Priority | Build tag |
| --------- | ------ | ---------- | -------- | --------- |
| `ordo.JSON` | json | `.json`    | 20       | —         |
| `ordo.XML`  | xml  | `.xml`     | 80       | `noxml`   |

## Custom formats

Implement the `Format` interface, or use `NewSimpleFormat`, and pass it like
any built-in preset:

```go
var ini = ordo.NewSimpleFormat(
    "ini", []string{".ini"}, 40,
    func(b []byte, v any) error { /* decode */ },
    func(v any) ([]byte, error) { /* encode; omit or return unsupported */ },
)

cfg, err := ordo.Load[Config](path, ini)
```
