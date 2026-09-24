# ordo

> A generic, format-pluggable configuration loader for Go.

## Usage

```go
loader := &ordo.Loader[Config]{
    Formats: []format.Format{format.JSON, format.XML},
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

Without orchestration, a load with default/merge support:

```go
cfg, err := ordo.Load[Config](path, format.JSON)
```

`Loader` additionally supports seeding from a default config and a custom
merge hook:

```go
loader := &ordo.Loader[Config]{
    Formats: []format.Format{format.JSON},
    Default: &Config{Address: ":3000"},
    Merge: func(base *Config, f format.Format, data []byte) (*Config, error) {
        // decode over the base as you see fit
        return base, f.Unmarshal(data, base)
    },
}
```

## Built-in formats

Presets live in the `format` package:

| Preset        | Format | Extensions | Priority |
| ------------- | ------ | ---------- | -------- |
| `format.JSON` | json   | `.json`    | 10       |
| `format.XML`  | xml    | `.xml`     | 10       |

## Custom formats

Implement the `format.Format` interface, or use `format.NewFormatter`, and
pass it like any built-in preset:

```go
var ini = format.NewFormatter(
    "ini", []string{".ini"}, 40,
    func(b []byte, v any) error { /* decode */ },
    func(v any) ([]byte, error) { /* encode; omit or return unsupported */ },
)

cfg, err := ordo.Load[Config](path, ini)
```
