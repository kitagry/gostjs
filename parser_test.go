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
					Type: &basicFieldType{
						Name: "lowercase",
					},
				},
				{
					Name: "Name",
					Type: &basicFieldType{
						Name: "string",
					},
					Document: "Name is name\n",
					Tags: map[string]string{
						"json": "name",
					},
				},
				{
					Name: "ID",
					Type: &basicFieldType{
						Name: "int",
					},
					Tags: map[string]string{
						"json": "id",
						"yaml": "id",
					},
				},
				{
					Name: "Array",
					Type: &arrayFieldType{
						Value: &basicFieldType{
							Name: "string",
						},
					},
					Tags: map[string]string{
						"json": "array",
					},
				},
				{
					Name: "Map",
					Type: &mapFieldType{
						Key: &basicFieldType{
							Name: "string",
						},
						Value: &basicFieldType{
							Name: "string",
						},
					},
					Tags: map[string]string{
						"json": "map",
					},
				},
				{
					Name: "Selector",
					Type: &unknownFieldType{},
				},
				{
					Name: "Interface",
					Type: &interfaceFieldType{},
				},
				{
					Name: "Child",
					Type: &basicFieldType{
						Name: "Child",
					},
					Tags: map[string]string{
						"json": "child",
					},
				},
				{
					Name: "Child2",
					Type: &starFieldType{
						&basicFieldType{
							Name: "Child",
						},
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
					Type: &basicFieldType{
						Name: "string",
					},
				},
			},
		},
		"lowercase": {
			Name: "lowercase",
			Fields: []Field{
				{
					Name: "Name",
					Type: &basicFieldType{
						Name: "string",
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(nameToDoc, expected) {
		t.Errorf("expected:\n%+v\ngot:\n%+v", expected, nameToDoc)
	}
}
