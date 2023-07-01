package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func processLine(key string, t map[string]string, fileList *[]string, currentFileNum int) {
	filenames, ok := t[key]
	if !ok {
		filenames = strings.Repeat(",0", len(*fileList))
	}
	fb := []byte(filenames)
	fb[currentFileNum*2+1] = '1'
	t[key] = string(fb)
}

func processFile(fileName string, tm map[string]string, fileList *[]string, currentFileIndex int) {
	log.Println("started ", fileName)
	var i uint64
	fd, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		key := scanner.Bytes()
		processLine(string(key), tm, fileList, currentFileIndex)
		i++
		if i%1000000 == 0 {
			log.Println(fileName, i)
		}
	}
}

func dump(tm map[string]string, fileList []string) {
	out, err := os.Create("out.csv")
	if err != nil {
		log.Fatal(err)
	}
	out.WriteString("phone")
	for _, fileName := range fileList {
		out.WriteString("," + fileName)
	}
	out.WriteString("\n")
	for key, val := range tm {
		out.WriteString(key + val + "\n")
	}
}

func main() {
	tm := map[string]string{}
	fileList := []string{}

	abs, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	files, err := ioutil.ReadDir(abs)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.Name() == "out.csv" {
			continue
		}
		fileList = append(fileList, file.Name())
	}
	for i, fileName := range fileList {
		processFile(fileName, tm, &fileList, i)
	}
	dump(tm, fileList)
}
