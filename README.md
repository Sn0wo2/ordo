# ordo

A generic, format-pluggable configuration loader for Go.

## Features

- **Formats**: JSON (stdlib, always available), YAML (`gopkg.in/yaml.v3`), TOML (`github.com/pelletier/go-toml/v2`), HCL (`github.com/hashicorp/hcl`), INI (`gopkg.in/ini.v1`), ENV (`github.com/joho/godotenv`, read-only), EDN (`olympos.io/encoding/edn`)
- **Opt-out build tags**: exclude a format — and its third-party dependency — from the build
- **Generic loader** with optional injected hooks: defaults, validation, on-loaded callback
- **Format probing**: files with unregistered extensions are probed against all registered formats in priority order
- **Safe persistence**: `Save` marshals by extension and writes with `0o750` dirs / `0o600` files

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

| Tag      | Effect                     |
|----------|----------------------------|
| (none)   | all formats                |
| `noyaml` | YAML support excluded      |
| `notoml` | TOML support excluded      |
| `nohcl`  | HCL support excluded       |
| `noini`  | INI support excluded       |
| `noenv`  | ENV support excluded       |
| `noedn`  | EDN support excluded       |

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

## License

Apache-2.0
