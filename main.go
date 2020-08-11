package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func parsePkg(item string) (string, string, error) {
	ind := strings.LastIndex(item, ".")
	if ind == -1 {
		return "", "", fmt.Errorf(`item should format "github.com/hoge/fuga.HogeFuga"`)
	}
	return item[:ind], item[ind+1:], nil
}

func main() {
	tag := flag.String("tag", "json", "tag name")
	srcPath := flag.String("src-path", "", "go build src path")
	output := flag.String("output", "", "output file name")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "args should be set 1")
		os.Exit(1)
	}

	item := args[0] // github.com/kitagry/gostjs/test.Test
	pkgPath, target, err := parsePkg(item)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse Pkg: %v", err)
		os.Exit(1)
	}

	if *srcPath == "" {
		*srcPath, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get pwd: %v\n", err)
			os.Exit(1)
		}
	}
	structs, err := Parse(pkgPath, *srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse `%s`: %v\n", pkgPath, err)
		os.Exit(1)
	}
	b, err := Decode(target, structs, *tag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode: %v\n", err)
		os.Exit(1)
	}

	if output != nil && *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create `%s`: %v\n", *output, err)
			os.Exit(1)
		}
		defer f.Close()

		_, err = f.Write(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write: %v\n", err)
			os.Exit(1)
		}
		return
	}
	_, err = os.Stdout.Write(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write to stdout: %v\n", err)
		os.Exit(1)
	}
}
