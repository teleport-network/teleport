package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

func findDiffModule(falseDir string, rightDir string) ([]string, []string, []string) {
	var module []string
	var rightFiles []string
	var falseFiles []string

	falseFileList, err := ioutil.ReadDir(falseDir)
	if err != nil {
		panic(err)
	}
	fileMap := make(map[string]int)
	for _, value := range falseFileList {
		fileMap[value.Name()] = 1
	}
	rightFileList, err := ioutil.ReadDir(rightDir)
	if err != nil {
		panic(err)
	}
	for _, value := range rightFileList {
		if v, ok := fileMap[value.Name()]; ok {
			_ = v
			delete(fileMap, value.Name())
		}
	}
	for key, value := range fileMap {
		prefix := strings.Split(key, "_")[0]
		module = append(module, prefix)
		_ = value
	}

	for _, v := range module {
		rightFiles = append(rightFiles, findFileByModule(rightFileList, v))
		falseFiles = append(falseFiles, findFileByModule(falseFileList, v))
	}

	return module, rightFiles, falseFiles
}

func findFileByModule(lists []fs.FileInfo, module string) string {
	for _, v := range lists {
		if strings.Contains(v.Name(), module) {
			return v.Name()
		}
	}
	return ""
}

func findDiffKeyInModule(rightFileName string, falseFileName string) ([]string, []string, []string) {
	var keys, rightValues, falseValues []string

	fmt.Println(rightFileName)
	fmt.Println(falseFileName)

	rightLines, err := ReadLine(rightFileName)
	if err != nil {
		panic(err)
	}
	falseLines, err := ReadLine(falseFileName)
	if err != nil {
		panic(err)
	}

	if len(rightLines) != len(falseLines) {
		panic("not equal")
	}

	rightMap := make(map[string]string)
	for _, v := range rightLines {
		data := strings.Split(v, ",")
		rightMap[data[0]] = data[1]
	}

	for _, v := range falseLines {
		data := strings.Split(v, ",")
		if strings.Compare(data[1], rightMap[data[0]]) == 0 {
			delete(rightMap, data[0])
		} else {
			keys = append(keys, data[0])
			rightValues = append(rightValues, rightMap[data[0]])
			falseValues = append(falseValues, data[1])
		}
	}

	return keys, rightValues, falseValues
}

func ReadLine(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	var result []string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				return result, nil
			}
			return nil, err
		}
		result = append(result, line)
	}
}
