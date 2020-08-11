package main

import (
	"encoding/json"
	"fmt"

	"golang.org/x/xerrors"
)

type Type string

const (
	String  Type = "string"
	Number       = "number"
	Object       = "object"
	Array        = "array"
	Boolean      = "boolean"
	Null         = "null"
)

type Schema struct {
	*Property

	Schema string `json:"$schema"`
}

type Property struct {
	Type Type `json:"type"`

	Description string `json:"description,omitempty"`

	Properties map[string]*Property `json:"properties,omitempty"`

	Contains *Property `json:"contains,omitempty"`

	Required []string `json:"required,omitempty"`
}

func nameFromTag(field Field, tag string) string {
	if tag == "" {
		return field.Name
	}

	name, ok := field.Tags[tag]
	if !ok {
		return field.Name
	}
	return name
}

func decodeFromName(target string, structs map[string]StructDoc, tagName string) (*Property, error) {
	switch target {
	case "boolean":
		return &Property{
			Type: Boolean,
		}, nil
	case "string":
		return &Property{
			Type: String,
		}, nil
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return &Property{
			Type: Number,
		}, nil
	case "int", "int8", "int16", "int32", "int64":
		return &Property{
			Type: Number,
		}, nil
	case "float32", "float64":
		return &Property{
			Type: Number,
		}, nil
	default:
		st, ok := structs[target]
		if !ok {
			return nil, fmt.Errorf("Unexpected target: %s", target)
		}
		p := Property{
			Type:        Object,
			Description: st.Document,
			Properties:  make(map[string]*Property),
			Required:    make([]string, 0),
		}
		for _, s := range st.Fields {
			property, err := decodeStructDoc(s.Type, structs, tagName)
			if err != nil {
				continue
			}
			property.Description = s.Document
			name := nameFromTag(s, tagName)
			p.Properties[name] = property
			if s.Type.Required {
				p.Required = append(p.Required, name)
			}
		}
		return &p, nil
	}
}

func decodeStructDoc(target *FieldType, structs map[string]StructDoc, tagName string) (*Property, error) {
	if target.IsArray {
		p := &Property{
			Type: Array,
		}
		var err error
		p.Contains, err = decodeFromName(target.Name, structs, tagName)
		if err != nil {
			return nil, err
		}
		return p, nil
	}
	return decodeFromName(target.Name, structs, tagName)
}

func Decode(target string, structs map[string]StructDoc, tagName string) ([]byte, error) {
	property, err := decodeFromName(target, structs, tagName)
	if err != nil {
		return []byte{}, xerrors.Errorf("failed to decodeStructDoc: %w", err)
	}

	schema := Schema{Property: property, Schema: "http://json-schema.org/schema#"}
	result, err := json.Marshal(schema)
	if err != nil {
		return []byte{}, xerrors.Errorf("failed to json.Marshal: %w", err)
	}

	return result, nil
}
