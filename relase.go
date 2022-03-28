package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

// ReleaseFile の tag を書き換える
// func writeReleaseFile(filePath, prefix, tag string) {
// 	if prefix == "" {
//
// 	}
//
// 	os.OpenFile("")
// }

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

	return strings.Join(ary, "."), prefix(oldTag)
}

// workflow dispatch 入力値を反映させる
// 既存タグとの比較 or
func replacement(tag string) (string, string) {
	if !validate(tag) {
		log.Fatalf("Input value is not a tag: %v", tag)
	}
	return tag, prefix(tag)
}

func validate(str string) bool {
	// v1.0.0
	// hogehoge/v1.0.0 のパターンに対応する
	reg := regexp.MustCompile(`v([0-9]+).([0-9]+).([0-9]+)`)
	return reg.MatchString(str)
}

func prefix(tag string) string {
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
