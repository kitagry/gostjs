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

	nameToDoc, err := Parse("github.com/kitagry/gojs/test", pwd)
	if err != nil {
		t.Fatalf("failed to Parse: %v", err)
	}

	expected := map[string]StructDoc{
		"Test": {
			Name:     "Test",
			Document: "Test's docs\n",
			Fields: []Field{
				{
					Type: "lowercase",
				},
				{
					Name:     "Name",
					Type:     "string",
					Document: "Name is name\n",
					Tags: map[string]string{
						"json": "name",
					},
				},
				{
					Name: "ID",
					Type: "int",
					Tags: map[string]string{
						"json": "id",
						"yaml": "id",
					},
				},
				{
					Name: "Child",
					Type: "Child",
					Tags: map[string]string{
						"json": "child",
					},
				},
				{
					Name: "Child2",
					Type: "*Child",
				},
			},
		},
		"Child": {
			Name:     "Child",
			Document: "Child is Test's child\n",
			Fields: []Field{
				{
					Name: "Name",
					Type: "string",
				},
			},
		},
		"lowercase": {
			Name: "lowercase",
			Fields: []Field{
				{
					Name: "Name",
					Type: "string",
				},
			},
		},
	}

	if !reflect.DeepEqual(nameToDoc, expected) {
		t.Errorf("expected %+v, got %+v", expected, nameToDoc)
	}
}
