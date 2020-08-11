package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

var parsedObj = map[string]StructDoc{
	"Test": {
		Name:     "Test",
		Document: "Test's docs\n",
		Fields: []Field{
			{
				Name:     "Name",
				Document: "Name's document\n",
				Type: &FieldType{
					Name:     "string",
					Required: true,
				},
				Tags: map[string]string{
					"json": "name",
				},
			},
			{
				Name: "ID",
				Type: &FieldType{
					Name: "int",
				},
			},
			{
				Name: "Score",
				Type: &FieldType{
					Name: "float64",
				},
				Tags: map[string]string{
					"json": "score",
				},
			},
			{
				Name: "IsBool",
				Type: &FieldType{
					Name: "boolean",
				},
			},
			{
				Name: "Strings",
				Type: &FieldType{
					Name:    "string",
					IsArray: true,
				},
			},
			{
				Name: "Child",
				Type: &FieldType{
					Name: "Child",
				},
				Tags: map[string]string{
					"json": "child",
				},
			},
		},
	},
	"Child": {
		Name: "Child",
		Fields: []Field{
			{
				Name: "Name",
				Type: &FieldType{
					Name:     "string",
					Required: true,
				},
				Tags: map[string]string{
					"json": "name",
				},
			},
		},
	},
}

func TestDecode(t *testing.T) {
	got, err := Decode("Test", parsedObj, "json")
	if err != nil {
		t.Fatalf("failed to Decode: %v", err)
	}

	var gotInterface map[string]interface{}
	err = json.Unmarshal(got, &gotInterface)
	if err != nil {
		t.Fatalf("failed to Unmarshal %v: %v", got, err)
	}

	expected := map[string]interface{}{
		"$schema":     "http://json-schema.org/schema#",
		"type":        "object",
		"description": "Test's docs\n",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Name's document\n",
			},
			"ID": map[string]interface{}{
				"type": "number",
			},
			"score": map[string]interface{}{
				"type": "number",
			},
			"IsBool": map[string]interface{}{
				"type": "boolean",
			},
			"Strings": map[string]interface{}{
				"type": "array",
				"contains": map[string]interface{}{
					"type": "string",
				},
			},
			"child": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []interface{}{"name"},
			},
		},
		"required": []interface{}{"name"},
	}

	if !reflect.DeepEqual(gotInterface, expected) {
		t.Errorf("Decode\nexpected:\n%+v\ngot:\n%+v", expected, gotInterface)
	}
}
