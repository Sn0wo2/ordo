# ordo

> A generic, format-pluggable configuration loader for Go.

## Usage

```go
loader := &ordo.Loader[Config]{
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

cfg, path, err := loader.Load("./config.yml")
```

If the path has no known extension, ordo tries every registered format in
priority order and uses the first that decodes; a bare `./config` also resolves
by trying each registered extension.

## Options

```go
cfg, err := ordo.Load[Config](path, ordo.WithFormat("json"))

err := ordo.Save(cfg, "./config.json")
```

- `WithFormat(name)` decode as the given format instead of by extension
- `Save` picks the format by extension, falling back to the highest-priority
  registered format

## Built-in formats

| Format | Extensions | Priority | Build tag |
| ------ | ---------- | -------- | --------- |
| json   | `.json`    | 20       | —         |
| xml    | `.xml`     | 80       | `noxml`   |

## Custom formats

Implement the `Format` interface, or use `NewSimpleFormat`:

```go
func init() {
    ordo.RegisterFormat(ordo.NewSimpleFormat(
        "ini", []string{".ini"}, 40,
        func(b []byte, v any) error { /* decode */ },
        func(v any) ([]byte, error) { /* encode; omit or return unsupported */ },
    ))
}
```

Formats without working `Marshal` should not be saved to; formats that cannot
write at all can simply return an error.
