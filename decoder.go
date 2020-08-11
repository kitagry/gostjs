package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

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

	Interface = ""
)

type Schema struct {
	*Property

	Schema string `json:"$schema"`
}

type Property struct {
	Type Type `json:"type"`

	Description string `json:"description,omitempty"`

	Properties map[string]*Property `json:"properties,omitempty"`

	AdditionalProperties *Property `json:"additionalProperties,omitempty"`

	Contains *Property `json:"contains,omitempty"`

	Required []string `json:"required,omitempty"`
}

func nameFromTag(field Field, tag string) (string, bool) {
	if tag == "" {
		return field.Name, false
	}

	name, ok := field.Tags[tag]
	if !ok {
		return field.Name, false
	}

	if name == "-" {
		return "", true
	}

	if strings.Index(name, ",") != -1 {
		name = strings.Split(name, ",")[0]
	}
	return name, false
}

func decodeFromName(target string, structs map[string]StructDoc, tagName string) (*Property, error) {
	switch target {
	case "bool":
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
	case "interface":
		return &Property{
			Type: Interface,
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
			name, skip := nameFromTag(s, tagName)
			if skip {
				continue
			}
			property, required, err := decodeStructDoc(s.Type, structs, tagName)
			if err != nil {
				return nil, fmt.Errorf("failed to decodeStructDoc: %v", err)
			}
			property.Description = s.Document
			p.Properties[name] = property
			if required {
				p.Required = append(p.Required, name)
			}
		}
		return &p, nil
	}
}

func decodeStructDoc(target FieldType, structs map[string]StructDoc, tagName string) (*Property, bool, error) {
	switch t := target.(type) {
	case *basicFieldType:
		p, err := decodeFromName(t.Name, structs, tagName)
		if err != nil {
			return nil, false, err
		}
		return p, false, nil
	case *starFieldType:
		p, _, err := decodeStructDoc(t.Value, structs, tagName)
		if err != nil {
			return nil, false, err
		}
		return p, false, nil
	case *arrayFieldType:
		p := &Property{
			Type: Array,
		}
		var err error
		p.Contains, _, err = decodeStructDoc(t.Value, structs, tagName)
		if err != nil {
			return nil, false, err
		}
		return p, false, nil
	case *mapFieldType:
		p := &Property{
			Type: Object,
		}
		var err error
		p.AdditionalProperties, _, err = decodeStructDoc(t.Value, structs, tagName)
		if err != nil {
			return nil, false, err
		}
		return p, false, nil
	case *interfaceFieldType:
		return &Property{
			Type: Interface,
		}, false, nil
	case *unknownFieldType:
		// TODO: Should parse different File's struct
		return &Property{
			Type: Interface,
		}, false, nil
	default:
		return nil, false, fmt.Errorf("Unexpected decode fieldType: %s", reflect.TypeOf(target))
	}
}

func Decode(target string, structs map[string]StructDoc, tagName string) ([]byte, error) {
	property, err := decodeFromName(target, structs, tagName)
	if err != nil {
		return []byte{}, xerrors.Errorf("failed to decodeFromName: %w", err)
	}

	schema := Schema{Property: property, Schema: "http://json-schema.org/schema#"}
	result, err := json.Marshal(schema)
	if err != nil {
		return []byte{}, xerrors.Errorf("failed to json.Marshal: %w", err)
	}

	return result, nil
}
