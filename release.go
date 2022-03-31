package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

const (
	defaultPath = "./RELEASE"
)

// ReleaseFile の tag を書き換える
func writeReleaseFile(filePath, tag string) {
	if filePath == "" {
		filePath = defaultPath
	}
	f, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	parseFile, err := parser.ParseBytes(f, parser.ParseComments)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	yamlPath, err := yaml.PathString("$.tag")
	if err != nil {
		log.Fatalf("failed to file path: %v", err)
	}

	oldNode, err := yamlPath.FilterFile(parseFile)
	if err != nil {
		log.Fatalf("failed to read old node")
	}
	newNode := &ast.StringNode{
		BaseNode: &ast.BaseNode{},
		Token:    oldNode.GetToken(),
	}
	oldTag := oldNode.String()

	var newTag, prefix string
	if tag == "" {
		newTag, prefix = increment(oldTag)
		newNode.Value = prefix + newTag
	} else {
		if oldTag == "" {
			log.Fatalf("tag is not exists")
		}
		newTag = replacement(tag)
		newNode.Value = newTag
	}

	if c := oldNode.GetComment(); c != nil {
		if err := newNode.SetComment(c); err != nil {
			log.Fatalf("failed to set comment: %v", err)
		}
	}

	if err := yamlPath.ReplaceWithNode(parseFile, newNode); err != nil {
		log.Fatalf("failed to replace file: %v", err)
	}

	if err := os.WriteFile(filePath, []byte(parseFile.String()), os.ModeExclusive); err != nil {
		log.Fatalf("failed to write file: %v", err)
	}
}

// version の increment実装
// If the tag is not added, it will be patch release.
func increment(oldTag string) (string, string) {
	if !validate(oldTag) {
		log.Fatalf("Input value is not a tag: %v", oldTag)
	}

	ary := strings.Split(oldTag, ".")

	oldVersion, err := strconv.Atoi(ary[len(ary)-1])
	if err != nil {
		log.Fatalf("failed to conv version: %v", err)
	}
	oldVersion++
	ary[len(ary)-1] = strconv.Itoa(oldVersion)

	return strings.Join(ary, "."), getPrefix(oldTag)
}

// workflow dispatch 入力値を反映させる
// 既存タグとの比較 or
func replacement(tag string) string {
	if !validate(tag) {
		log.Fatalf("Input value is not a tag: %v", tag)
	}
	return tag
}

func validate(str string) bool {
	// v1.0.0, hogehoge/v1.0.0 のパターンに対応する
	reg := regexp.MustCompile(`v([0-9]+).([0-9]+).([0-9]+)`)
	return reg.MatchString(str)
}

func getPrefix(tag string) string {
	splitTag := strings.Split(tag, "/")
	// hoge/microservice/v1.0.0 のようなパターンに対応
	var pr string
	for _, v := range splitTag {
		if !validate(v) {
			pr += v + "/"
		}
	}
	return pr
}
