// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package build

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Lang string

func (l Lang) String() string {
	return string(l)
}

// 未识别语言类型
var NO Lang = "no"

// Python
var Python Lang = "python"

// Php
var PHP Lang = "php"

// Java
var Java Lang = "java"

// Node
var Node Lang = "node"

// Go
var Go Lang = "go"

type langTypeFunc func(homepath string) Lang

var checkFuncList []langTypeFunc

func init() {
	checkFuncList = append(checkFuncList, python)
	checkFuncList = append(checkFuncList, php)
	checkFuncList = append(checkFuncList, java)
	checkFuncList = append(checkFuncList, node)
	checkFuncList = append(checkFuncList, golang)
}

func GetLangType(homepath string) Lang {
	if ok, _ := fileExists(homepath); !ok {
		return NO
	}
	//判断是否有代码
	if ok := isHaveFile(homepath); !ok {
		return NO
	}
	//获取确定的语言
	for _, check := range checkFuncList {
		if lang := check(homepath); lang != NO {
			return lang
		}
	}
	//无法识别
	return NO
}

func python(homepath string) Lang {
	if ok, _ := fileExists(path.Join(homepath, "requirements.txt")); ok {
		return Python
	}
	if ok, _ := fileExists(path.Join(homepath, "setup.py")); ok {
		return Python
	}
	if ok, _ := fileExists(path.Join(homepath, "Pipfile")); ok {
		return Python
	}
	return NO
}

func php(homepath string) Lang {
	if ok, _ := fileExists(path.Join(homepath, "composer.json")); ok {
		return PHP
	}
	if ok := searchFile(homepath, "index.php", 2); ok {
		return PHP
	}
	return NO
}

func java(homepath string) Lang {
	if ok, _ := fileExists(path.Join(homepath, "pom.xml")); ok {
		return Java
	}
	if ok, _ := fileExists(path.Join(homepath, "pom.atom")); ok {
		return Java
	}
	if ok, _ := fileExists(path.Join(homepath, "pom.clj")); ok {
		return Java
	}
	if ok, _ := fileExists(path.Join(homepath, "pom.groovy")); ok {
		return Java
	}
	if ok, _ := fileExists(path.Join(homepath, "pom.rb")); ok {
		return Java
	}
	if ok, _ := fileExists(path.Join(homepath, "pom.scala")); ok {
		return Java
	}
	if ok, _ := fileExists(path.Join(homepath, "pom.yaml")); ok {
		return Java
	}
	if ok, _ := fileExists(path.Join(homepath, "pom.yml")); ok {
		return Java
	}
	if ok := fileExistsWithSuffix(homepath, ".war"); ok {
		return Java
	}
	if ok := fileExistsWithSuffix(homepath, ".jar"); ok {
		return Java
	}
	return NO
}

func node(homepath string) Lang {
	if ok, _ := fileExists(path.Join(homepath, "package.json")); ok {
		return Node
	}
	return NO
}

func golang(homepath string) Lang {
	if ok, _ := fileExists(path.Join(homepath, "go.mod")); ok {
		return Go
	}
	if ok, _ := fileExists(path.Join(homepath, "Gopkg.lock")); ok {
		return Go
	}
	if ok, _ := fileExists(path.Join(homepath, "Godeps", "Godeps.json")); ok {
		return Go
	}
	if ok, _ := fileExists(path.Join(homepath, "vendor", "vendor.json")); ok {
		return Go
	}
	if ok, _ := fileExists(path.Join(homepath, "glide.yaml")); ok {
		return Go
	}
	if ok := fileExistsWithSuffix(path.Join(homepath, "src"), ".go"); ok {
		return Go
	}
	return NO
}

func fileExists(filename string) (bool, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

//searchFileBody 搜索文件中是否含有指定字符串
// func searchFileBody(filename, searchStr string) bool {
// 	body, _ := ioutil.ReadFile(filename)
// 	return strings.Contains(string(body), searchStr)
// }

//searchFile 搜索指定目录是否有指定文件，指定搜索目录层数，-1为全目录搜索
func searchFile(pathDir, name string, level int) bool {
	if level == 0 {
		return false
	}
	files, _ := ioutil.ReadDir(pathDir)
	var dirs []os.FileInfo
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file)
			continue
		}
		if file.Name() == name {
			return true
		}
	}
	if level == 1 {
		return false
	}
	for _, dir := range dirs {
		ok := searchFile(path.Join(pathDir, dir.Name()), name, level-1)
		if ok {
			return ok
		}
	}
	return false
}

//fileExistsWithSuffix 指定目录是否含有指定后缀的文件
func fileExistsWithSuffix(pathDir, suffix string) bool {
	files, _ := ioutil.ReadDir(pathDir)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), suffix) {
			return true
		}
	}
	return false
}

//isHaveFile 指定目录是否含有文件
//.开头文件除外
func isHaveFile(path string) bool {
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			return true
		}
	}
	return false
}
