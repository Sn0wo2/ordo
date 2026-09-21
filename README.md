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
if err != nil {
    if errors.Is(err, os.ErrNotExist) { /* handle missing file */ }
}

// Persist back to the same file (format chosen by extension).
err = ordo.Save(cfg, path)
```

Primitives, without orchestration:

```go
cfg, err := ordo.Load[Config](path) // decode only
err = ordo.Save(cfg, "./config.toml")
```

## Build tags

| Tag      | Effect                |
| -------- | --------------------- |
| (none)   | all formats           |
| `noyaml` | YAML support excluded |
| `notoml` | TOML support excluded |
| `nohcl`  | HCL support excluded  |
| `noini`  | INI support excluded  |
| `noenv`  | ENV support excluded  |
| `noedn`  | EDN support excluded  |

The env format is flat `KEY=value` and therefore read-only: `Save` refuses
it. Read-only formats simply don't implement the optional `Marshaler`
interface.

## Custom formats

Implement the `Format` interface (decoding) and — optionally — `Marshaler`
(writing) and register it (typically from an `init`):

```go
type iniFormat struct{}

func (iniFormat) Name() string                     { return "ini" }
func (iniFormat) Extensions() []string             { return []string{".ini"} }
func (iniFormat) Priority() int                    { return 40 }
func (iniFormat) Unmarshal(b []byte, v any) error  { /* ... */ }

func (iniFormat) Marshal(v any) ([]byte, error) { /* optional */ }

func init() { ordo.RegisterFormat(iniFormat{}) }
```
