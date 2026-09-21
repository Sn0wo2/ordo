//go:build !nohcl

package ordo

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type hclFormat struct{}

func (hclFormat) Name() string         { return "hcl" }
func (hclFormat) Extensions() []string { return []string{".hcl"} }
func (hclFormat) Priority() int        { return 60 }

func (hclFormat) Unmarshal(b []byte, v any) error {
	file, diags := hclsyntax.ParseConfig(b, "config.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return diags
	}

	m, err := hclBodyToMap(file.Body.(*hclsyntax.Body))
	if err != nil {
		return err
	}

	return decodeFlatMap(m, v)
}

// hclBodyToMap flattens an HCL body into a map: attributes keep their names,
// blocks become nested maps, and repeated blocks of the same type become a
// list of maps. Attribute values are converted to plain Go values.
func hclBodyToMap(body *hclsyntax.Body) (map[string]any, error) {
	out := make(map[string]any, len(body.Attributes)+len(body.Blocks))

	for name, attr := range body.Attributes {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, diags
		}

		goVal, err := hclToGo(val)
		if err != nil {
			return nil, err
		}

		out[name] = goVal
	}

	for _, block := range body.Blocks {
		nested, err := hclBodyToMap(block.Body)
		if err != nil {
			return nil, err
		}

		switch existing := out[block.Type].(type) {
		case nil:
			out[block.Type] = nested
		case map[string]any:
			out[block.Type] = []map[string]any{existing, nested}
		case []map[string]any:
			out[block.Type] = append(existing, nested)
		default:
			return nil, fmt.Errorf("hcl: %q is both an attribute and a block", block.Type)
		}
	}

	return out, nil
}

// hclToGo converts a cty value into plain Go values (maps, slices, strings,
// bools, float64) via its JSON representation.
func hclToGo(val cty.Value) (any, error) {
	raw, err := ctyjson.Marshal(val, val.Type())
	if err != nil {
		return nil, err
	}

	var out any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}

	return out, nil
}

func (hclFormat) Marshal(v any) ([]byte, error) {
	normalized, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var root map[string]any
	if err := json.Unmarshal(normalized, &root); err != nil {
		return nil, err
	}

	file := hclwrite.NewEmptyFile()
	writeHCLBody(file.Body(), root)

	return file.Bytes(), nil
}

// writeHCLBody emits a map as an HCL body: nested maps become blocks, lists
// of maps become one block per item, and everything else becomes an
// attribute. String escaping is handled by hclwrite.
func writeHCLBody(body *hclwrite.Body, m map[string]any) {
	for _, key := range hclKeys(m) {
		switch typed := m[key].(type) {
		case map[string]any:
			writeHCLBody(body.AppendNewBlock(hclAttrName(key), nil).Body(), typed)
		case []any:
			if !isHCLBlockList(typed) {
				body.SetAttributeValue(hclAttrName(key), hclValue(typed))

				continue
			}

			for _, item := range typed {
				if inner, ok := item.(map[string]any); ok {
					writeHCLBody(body.AppendNewBlock(hclAttrName(key), nil).Body(), inner)

					continue
				}

				body.SetAttributeValue(hclAttrName(key), hclValue(item))
			}
		default:
			body.SetAttributeValue(hclAttrName(key), hclValue(typed))
		}
	}
}

// hclKeys returns the map's keys sorted, omitting nil values.
func hclKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for key, val := range m {
		if val != nil {
			keys = append(keys, key)
		}
	}

	slices.Sort(keys)

	return keys
}

// isHCLBlockList reports whether a list contains composite items that must be
// written as repeated blocks instead of a single attribute.
func isHCLBlockList(items []any) bool {
	return slices.ContainsFunc(items, func(item any) bool {
		switch item.(type) {
		case map[string]any, []any, nil:
			return true
		}

		return false
	})
}

func hclValue(val any) cty.Value {
	switch typed := val.(type) {
	case map[string]any:
		attrs := make(map[string]cty.Value, len(typed))
		for key, inner := range typed {
			attrs[key] = hclValue(inner)
		}

		return cty.ObjectVal(attrs)
	case []any:
		items := make([]cty.Value, len(typed))
		for i, inner := range typed {
			items[i] = hclValue(inner)
		}

		return cty.TupleVal(items)
	case string:
		return cty.StringVal(typed)
	case bool:
		return cty.BoolVal(typed)
	case float64:
		return cty.NumberFloatVal(typed)
	default:
		return cty.NullVal(cty.DynamicPseudoType)
	}
}

// hclAttrName quotes keys that are not valid HCL identifiers.
func hclAttrName(key string) string {
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
