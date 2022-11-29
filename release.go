package main

import (
	"fmt"
	"io"
	"os"

	"github.com/go-git/go-billy/v5"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

func readReleaseFile(fs billy.Filesystem, filePath string) (string, *ast.StringNode, *yaml.Path, *ast.File, error) {
	file, err := fs.Open(filePath)
	if err != nil {
		fmt.Printf("failed to open file: %s\n", filePath)
	}

	f, _ := io.ReadAll(file)
	parseFile, err := parser.ParseBytes(f, parser.ParseComments)
	if err != nil {
		fmt.Println("failed to read file")
		return "", nil, nil, nil, err
	}

	tagNode, err := yaml.PathString("$.tag")
	if err != nil {
		fmt.Println("failed to parse tag field")
		return "", nil, nil, nil, err
	}

	oldNode, err := tagNode.FilterFile(parseFile)
	if err != nil {
		fmt.Println("failed to get value of tag field")
		return "", nil, nil, nil, err
	}
	newNode := &ast.StringNode{
		BaseNode: &ast.BaseNode{},
		Token:    oldNode.GetToken(),
	}

	if c := oldNode.GetComment(); c != nil {
		if err := newNode.SetComment(c); err != nil {
			fmt.Println("failed to set comment for new tag field")
			return "", nil, nil, nil, err
		}
	}

	return oldNode.String(), newNode, tagNode, parseFile, nil
}

func writeFile(yamlPath *yaml.Path, fs billy.Filesystem, parseFile *ast.File, newNode *ast.StringNode, filePath string) error {
	if err := yamlPath.ReplaceWithNode(parseFile, newNode); err != nil {
		fmt.Println("failed to replace yaml file node")
		return err
	}
	fileBytes, err := io.ReadAll(parseFile)
	if err != nil {
		fmt.Println("failed to read file bytes")
		return err
	}
	fileBytesWithNewLine := append(fileBytes, "\n"...)
	f, err := fs.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("failed to open file %s\n", filePath)
		return err
	}
	if _, err := f.Write(fileBytesWithNewLine); err != nil {
		fmt.Printf("failed to write file %s\n", filePath)
		return err
	}
	return nil
}
