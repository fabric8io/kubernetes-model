package main

import (
	"fmt"
	"reflect"
	"strings"
)

type schemaGenerator struct {
	types map[reflect.Type]*JSONObjectDescriptor
}

func GenerateSchema(t reflect.Type) (*JSONSchema, error) {
	g := schemaGenerator{types: make(map[reflect.Type]*JSONObjectDescriptor)}
	return g.generate(t)
}

func getFieldName(f reflect.StructField) string {
	json := f.Tag.Get("json")
	if len(json) > 0 {
		parts := strings.Split(json, ",")
		return parts[0]
	}
	return f.Name
}

func qualifiedName(t reflect.Type) string {
	prefix := strings.Replace(t.PkgPath(), "/", "_", -1)
	prefix = strings.Replace(prefix, ".", "_", -1)
	prefix = strings.Replace(prefix, "-", "_", -1)
	return prefix + "_" + t.Name()
}

func generateReference(t reflect.Type) string {
	return "#/definitions/" + qualifiedName(t)
}

func (g *schemaGenerator) generate(t reflect.Type) (*JSONSchema, error) {
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Only struct types can be converted.")
	}

	s := JSONSchema{
		ID:     "http://openshift.com/origin/v3/" + t.Name() + "#",
		Schema: "http://json-schema.org/schema#",
		JSONDescriptor: JSONDescriptor{
			Type: "object",
		},
	}
	s.JSONObjectDescriptor = g.generateObjectDescriptor(t)
	if len(g.types) > 0 {
		s.Definitions = make(map[string]JSONPropertyDescriptor)
		for k, v := range g.types {
			name := qualifiedName(k)
			value := JSONPropertyDescriptor{
				JSONDescriptor: &JSONDescriptor{
					Type: "object",
				},
				JSONObjectDescriptor: v,
			}
			s.Definitions[name] = value
		}
	}
	return &s, nil
}

func (g *schemaGenerator) getPropertyDescriptor(t reflect.Type) JSONPropertyDescriptor {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Bool:
		return JSONPropertyDescriptor{
			JSONDescriptor: &JSONDescriptor{
				Type: "boolean",
			},
		}
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		return JSONPropertyDescriptor{
			JSONDescriptor: &JSONDescriptor{
				Type: "integer",
			},
		}
	case reflect.Float32, reflect.Float64, reflect.Complex64,
		reflect.Complex128:
		return JSONPropertyDescriptor{
			JSONDescriptor: &JSONDescriptor{
				Type: "number",
			},
		}
	case reflect.String:
		return JSONPropertyDescriptor{
			JSONDescriptor: &JSONDescriptor{
				Type: "string",
			},
		}
	case reflect.Array:
	case reflect.Slice:
		return JSONPropertyDescriptor{
			JSONDescriptor: &JSONDescriptor{
				Type: "array",
			},
			JSONArrayDescriptor: &JSONArrayDescriptor{
				Items: g.getPropertyDescriptor(t.Elem()),
			},
		}
	case reflect.Map:
		return JSONPropertyDescriptor{
			JSONDescriptor: &JSONDescriptor{
				Type: "object",
			},
			JSONMapDescriptor: &JSONMapDescriptor{
				MapValueType: g.getPropertyDescriptor(t.Elem()),
			},
		}
	case reflect.Struct:
		definedType, ok := g.types[t]
		if !ok {
			g.types[t] = &JSONObjectDescriptor{}
			definedType = g.generateObjectDescriptor(t)
			g.types[t] = definedType
		}
		return JSONPropertyDescriptor{
			JSONReferenceDescriptor: &JSONReferenceDescriptor{
				Reference: generateReference(t),
			},
		}
	}
	return JSONPropertyDescriptor{}
}

func (g *schemaGenerator) getStructProperties(t reflect.Type) map[string]JSONPropertyDescriptor {
	props := map[string]JSONPropertyDescriptor{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if len(field.PkgPath) > 0 { // Skip private fields
			continue
		}
		name := getFieldName(field)
		prop := g.getPropertyDescriptor(field.Type)
		if field.Anonymous {
			var newProps map[string]JSONPropertyDescriptor
			if prop.JSONReferenceDescriptor != nil {
				pType := field.Type
				if pType.Kind() == reflect.Ptr {
					pType = pType.Elem()
				}
				newProps = g.types[pType].Properties
			} else {
				newProps = prop.Properties
			}
			for k, v := range newProps {
				props[k] = v
			}
		} else {
			props[name] = prop
		}
	}
	return props
}
func (g *schemaGenerator) generateObjectDescriptor(t reflect.Type) *JSONObjectDescriptor {
	desc := JSONObjectDescriptor{}
	desc.Properties = g.getStructProperties(t)
	return &desc
}
