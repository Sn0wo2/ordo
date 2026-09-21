//go:build !nohcl

package ordo

import (
	"encoding/json"
	"fmt"
	"slices"

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

	var flatten func(body *hclsyntax.Body) (map[string]any, error)
	flatten = func(body *hclsyntax.Body) (map[string]any, error) {
		out := make(map[string]any, len(body.Attributes)+len(body.Blocks))

		for name, attr := range body.Attributes {
			val, valDiags := attr.Expr.Value(nil)
			if valDiags.HasErrors() {
				return nil, valDiags
			}

			raw, err := ctyjson.Marshal(val, val.Type())
			if err != nil {
				return nil, err
			}

			var goVal any
			if err := json.Unmarshal(raw, &goVal); err != nil {
				return nil, err
			}

			out[name] = goVal
		}

		for _, block := range body.Blocks {
			nested, err := flatten(block.Body)
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

	m, err := flatten(file.Body.(*hclsyntax.Body))
	if err != nil {
		return err
	}

	return decodeFlat(m, v)
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

	var value func(val any) cty.Value
	value = func(val any) cty.Value {
		switch typed := val.(type) {
		case map[string]any:
			attrs := make(map[string]cty.Value, len(typed))
			for key, inner := range typed {
				attrs[key] = value(inner)
			}

			return cty.ObjectVal(attrs)
		case []any:
			items := make([]cty.Value, len(typed))
			for i, inner := range typed {
				items[i] = value(inner)
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

	var write func(body *hclwrite.Body, m map[string]any)
	write = func(body *hclwrite.Body, m map[string]any) {
		keys := make([]string, 0, len(m))
		for key, val := range m {
			if val != nil {
				keys = append(keys, key)
			}
		}

		slices.Sort(keys)

		for _, key := range keys {
			switch typed := m[key].(type) {
			case map[string]any:
				write(body.AppendNewBlock(key, nil).Body(), typed)
			case []any:
				composite := slices.ContainsFunc(typed, func(item any) bool {
					switch item.(type) {
					case map[string]any, []any, nil:
						return true
					}

					return false
				})

				if !composite {
					body.SetAttributeValue(key, value(typed))

					continue
				}

				for _, item := range typed {
					if inner, ok := item.(map[string]any); ok {
						write(body.AppendNewBlock(key, nil).Body(), inner)

						continue
					}

					body.SetAttributeValue(key, value(item))
				}
			default:
				body.SetAttributeValue(key, value(typed))
			}
		}
	}

	file := hclwrite.NewEmptyFile()
	write(file.Body(), root)

	return file.Bytes(), nil
}

func init() { RegisterFormat(hclFormat{}) }
