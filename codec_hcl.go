//go:build !nohcl

package ordo

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl"
)

type hclFormat struct{}

func (hclFormat) Name() string                    { return "hcl" }
func (hclFormat) Extensions() []string            { return []string{".hcl"} }
func (hclFormat) Priority() int                   { return 60 }
func (hclFormat) Unmarshal(b []byte, v any) error { return hcl.Unmarshal(b, v) }

// Marshal encodes v as HCL by normalizing it through JSON first (so struct
// tags decide the keys) and emitting assignments and blocks. Values that are
// nil are skipped: HCL v1 has no null literal.
func (hclFormat) Marshal(v any) ([]byte, error) {
	normalized, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var root map[string]any
	if err := json.Unmarshal(normalized, &root); err != nil {
		return nil, err
	}

	var sb strings.Builder

	emitHCLMap(&sb, root, 0)

	return []byte(sb.String()), nil
}

func emitHCLMap(sb *strings.Builder, m map[string]any, depth int) {
	keys := make([]string, 0, len(m))
	for k := range m {
		if m[k] == nil {
			continue
		}

		keys = append(keys, k)
	}

	sort.Strings(keys)

	indent := strings.Repeat("  ", depth)

	for _, key := range keys {
		emitHCLValue(sb, hclKey(key), m[key], indent, depth)
	}
}

func emitHCLValue(sb *strings.Builder, key string, val any, indent string, depth int) {
	switch typed := val.(type) {
	case map[string]any:
		fmt.Fprintf(sb, "%s%s {\n", indent, key)
		emitHCLMap(sb, typed, depth+1)
		fmt.Fprintf(sb, "%s}\n\n", indent)
	case []any:
		if isScalarList(typed) {
			fmt.Fprintf(sb, "%s%s = %s\n", indent, key, emitHCLList(typed))

			return
		}

		for _, item := range typed {
			emitHCLValue(sb, key, item, indent, depth)
		}
	default:
		fmt.Fprintf(sb, "%s%s = %s\n", indent, key, emitHCLScalar(val))
	}
}

func isScalarList(items []any) bool {
	for _, item := range items {
		switch item.(type) {
		case map[string]any, []any, nil:
			return false
		}
	}

	return true
}

func emitHCLList(items []any) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		parts = append(parts, emitHCLScalar(item))
	}

	return "[" + strings.Join(parts, ", ") + "]"
}

func emitHCLScalar(val any) string {
	switch typed := val.(type) {
	case string:
		// Escape HCL template interpolation so literals survive a round-trip.
		return strings.ReplaceAll(strconv.Quote(typed), "${", "$${}")
	case bool:
		return strconv.FormatBool(typed)
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", typed)
	}
}

func hclKey(key string) string {
	for _, r := range key {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' || r == '-' {
			continue
		}

		return strconv.Quote(key)
	}

	if key == "" {
		return strconv.Quote(key)
	}

	return key
}

func init() { RegisterFormat(hclFormat{}) }
