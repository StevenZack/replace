package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/StevenZack/tools/fileToolkit"
)

var (
	file = flag.String("f", "", "specific your file")
)

func main() {
	flag.Parse()
	content, e := fileToolkit.ReadFileAll(*file)
	if e != nil {
		fmt.Println("read file error :", e)
		return
	}
	if len(flag.Args()) < 2 {
		fmt.Println("not enough args")
		return
	}
	oldStr := flag.Arg(0)
	newStr := flag.Arg(1)
	newContent := strings.ReplaceAll(content, oldStr, newStr)
	f, e := fileToolkit.WriteFile(*file)
	if e != nil {
		fmt.Println("open file to write error :", e)
		return
	}
	_, e = f.WriteString(newContent)
	if e != nil {
		fmt.Println("write error :", e)
		return
	}
}
