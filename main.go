package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/StevenZack/tools/fileToolkit"
)

var (
	file = flag.String("f", "", "specific your file")
)

func main() {
	flag.Parse()
	if *file != "" {
		handleFile(*file)
		return
	}
	wd, e := os.Getwd()
	if e != nil {
		log.Fatal(e)
	}
	e = filepath.Walk(wd, func(path string, f os.FileInfo, e error) error {
		if f == nil || f.IsDir() {
			return nil
		}
		switch filepath.Ext(f.Name()) {
		case ".go":
		default:
			return nil
		}
		handleFile(path)
		return nil
	})
	if e != nil {
		log.Fatal(e)
	}
}

func handleFile(path string) {
	content, e := fileToolkit.ReadFileAll(path)
	if e != nil {
		fmt.Println("read file error :", e)
		return
	}
	oldStr := flag.Arg(0)
	newStr := flag.Arg(1)
	newContent := strings.ReplaceAll(content, oldStr, newStr)
	e = fileToolkit.WriteFile(path, []byte(newContent))
	if e != nil {
		fmt.Println(" write file error :", e)
		return
	}
}
