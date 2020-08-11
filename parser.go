package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/fatih/structtag"
	"golang.org/x/xerrors"
)

type StructDoc struct {
	Name     string
	Document string
	Fields   []Field
}

type Field struct {
	Name     string
	Type     *FieldType
	Document string
	Tags     map[string]string
}

type FieldType struct {
	Name     string
	IsArray  bool
	Required bool
}

func getFieldType(expr ast.Expr) (*FieldType, error) {
	switch t := expr.(type) {
	case *ast.Ident:
		return &FieldType{
			Name:     t.Name,
			Required: true,
		}, nil
	case *ast.StarExpr:
		ty, err := getFieldType(t.X)
		if err != nil {
			return nil, err
		}
		ty.Required = false
		return ty, nil
	case *ast.ArrayType:
		ty, err := getFieldType(t.Elt)
		if err != nil {
			return nil, err
		}
		ty.IsArray = true
		return ty, nil
	default:
		return nil, fmt.Errorf("Unimplemented expr type: %s", reflect.TypeOf(expr))
	}
}

func parseField(l *ast.Field) (Field, error) {
	field := Field{}

	if len(l.Names) == 1 {
		field.Name = l.Names[0].String()
	}
	field.Document = l.Doc.Text()

	var err error
	field.Type, err = getFieldType(l.Type)
	if err != nil {
		return field, xerrors.Errorf("failed to get FieldType: %w", err)
	}

	if l.Tag == nil || len(l.Tag.Value) == 0 {
		return field, nil
	}

	tags, err := structtag.Parse(strings.Trim(l.Tag.Value, "`"))
	if err != nil {
		return field, xerrors.Errorf("failed to parse struct tag: %w", err)
	}

	field.Tags = make(map[string]string)
	for _, t := range tags.Tags() {
		field.Tags[t.Key] = t.Value()
	}

	return field, nil
}

func Parse(pkgPath, srcPath string) (map[string]StructDoc, error) {
	pkg, err := build.Import(pkgPath, srcPath, 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Couldn't find package: %v", err))
	}

	fset := token.NewFileSet()
	result := make(map[string]StructDoc)
	for _, file := range pkg.GoFiles {
		fmt.Println(file)
		f, err := parser.ParseFile(fset, filepath.Join(pkg.Dir, file), nil, parser.ParseComments)
		if err != nil {
			fmt.Printf("build Import error: %v", err)
			continue
		}

		for _, decl := range f.Decls {
			// ast.Print(fset, decl)
			switch d := decl.(type) {
			case *ast.GenDecl:
				tok := d.Tok.String()
				if tok != "type" {
					continue
				}
				for _, spec := range d.Specs {
					s, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					st, ok := s.Type.(*ast.StructType)
					if !ok {
						continue
					}
					doc := StructDoc{}
					doc.Name = s.Name.String()
					doc.Document = d.Doc.Text()
					for _, l := range st.Fields.List {
						field, err := parseField(l)
						if err != nil {
							return nil, xerrors.Errorf("failed to parseField: %w", err)
						}
						doc.Fields = append(doc.Fields, field)
					}
					result[s.Name.String()] = doc
				}
			}
		}
	}
	return result, nil
}
