package main

import (
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get pwd: %v", err)
	}

	nameToDoc, err := Parse("github.com/kitagry/gostjs/test", pwd)
	if err != nil {
		t.Fatalf("failed to Parse: %v", err)
	}

	expected := map[string]StructDoc{
		"Test": {
			Name:     "Test",
			Document: "Test's docs\n",
			Fields: []Field{
				{
					Type: &FieldType{
						Name:     "lowercase",
						Required: true,
					},
				},
				{
					Name: "Name",
					Type: &FieldType{
						Name:     "string",
						Required: true,
					},
					Document: "Name is name\n",
					Tags: map[string]string{
						"json": "name",
					},
				},
				{
					Name: "ID",
					Type: &FieldType{
						Name:     "int",
						Required: true,
					},
					Tags: map[string]string{
						"json": "id",
						"yaml": "id",
					},
				},
				{
					Name: "Array",
					Type: &FieldType{
						Name:     "string",
						Required: true,
						IsArray:  true,
					},
					Tags: map[string]string{
						"json": "array",
					},
				},
				{
					Name: "Child",
					Type: &FieldType{
						Name:     "Child",
						Required: true,
					},
					Tags: map[string]string{
						"json": "child",
					},
				},
				{
					Name: "Child2",
					Type: &FieldType{
						Name:     "Child",
						Required: false,
					},
				},
			},
		},
		"Child": {
			Name:     "Child",
			Document: "Child is Test's child\n",
			Fields: []Field{
				{
					Name: "Name",
					Type: &FieldType{
						Name:     "string",
						Required: true,
					},
				},
			},
		},
		"lowercase": {
			Name: "lowercase",
			Fields: []Field{
				{
					Name: "Name",
					Type: &FieldType{
						Name:     "string",
						Required: true,
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(nameToDoc, expected) {
		t.Errorf("expected:\n%+v\ngot:\n%+v", expected, nameToDoc)
	}
}
