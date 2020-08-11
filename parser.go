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
	"unicode"

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
	Type     FieldType
	Document string
	Tags     map[string]string
}

type FieldType interface {
}

type basicFieldType struct {
	Name string
}

type starFieldType struct {
	Value FieldType
}

type arrayFieldType struct {
	Value FieldType
}

type mapFieldType struct {
	Key   FieldType
	Value FieldType
}

type unknownFieldType struct {
}

type interfaceFieldType struct {
}

func getFieldType(expr ast.Expr) (FieldType, error) {
	switch t := expr.(type) {
	case *ast.Ident:
		return &basicFieldType{
			Name: t.Name,
		}, nil
	case *ast.StarExpr:
		ty, err := getFieldType(t.X)
		if err != nil {
			return nil, err
		}
		return &starFieldType{
			Value: ty,
		}, nil
	case *ast.ArrayType:
		ty, err := getFieldType(t.Elt)
		if err != nil {
			return nil, err
		}
		return &arrayFieldType{
			Value: ty,
		}, nil
	case *ast.MapType:
		tk, err := getFieldType(t.Key)
		if err != nil {
			return nil, err
		}
		tv, err := getFieldType(t.Value)
		if err != nil {
			return nil, err
		}
		return &mapFieldType{
			Key:   tk,
			Value: tv,
		}, nil
	case *ast.InterfaceType:
		return &interfaceFieldType{}, nil
	case *ast.SelectorExpr:
		return &unknownFieldType{}, nil
	default:
		return nil, fmt.Errorf("Unimplemented expr type: %s", reflect.TypeOf(expr))
	}
}

func parseField(l *ast.Field) (Field, bool, error) {
	field := Field{}

	if len(l.Names) == 1 {
		field.Name = l.Names[0].String()
		if unicode.IsLower(rune(field.Name[0])) {
			return field, true, nil
		}
	}
	field.Document = l.Doc.Text()

	var err error
	field.Type, err = getFieldType(l.Type)
	if err != nil {
		return field, false, xerrors.Errorf("failed to get FieldType of %s: %w", field.Name, err)
	}

	if l.Tag == nil || len(l.Tag.Value) == 0 {
		return field, false, nil
	}

	tags, err := structtag.Parse(strings.Trim(l.Tag.Value, "`"))
	if err != nil {
		return field, false, xerrors.Errorf("failed to parse struct tag: %w", err)
	}

	field.Tags = make(map[string]string)
	for _, t := range tags.Tags() {
		field.Tags[t.Key] = t.Value()
	}

	return field, false, nil
}

func Parse(pkgPath, srcPath string) (map[string]StructDoc, error) {
	pkg, err := build.Import(pkgPath, srcPath, 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Couldn't find package: %v", err))
	}

	fset := token.NewFileSet()
	result := make(map[string]StructDoc)
	for _, file := range pkg.GoFiles {
		f, err := parser.ParseFile(fset, filepath.Join(pkg.Dir, file), nil, parser.ParseComments)
		if err != nil {
			fmt.Printf("build Import error: %v", err)
			continue
		}

		for _, decl := range f.Decls {
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
						field, skip, err := parseField(l)
						if err != nil {
							return nil, xerrors.Errorf("failed to parseField: %w", err)
						}
						if skip {
							continue
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
