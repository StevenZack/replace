package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/StevenZack/tools/strToolkit"

	"github.com/StevenZack/tools/fileToolkit"
)

var (
	file        = flag.String("f", "", "specific your file")
	versionFile = flag.String("v", "", "from version.go file")
)

func main() {
	flag.Parse()
	content, e := fileToolkit.ReadFileAll(*file)
	if e != nil {
		fmt.Println("read file error :", e)
		return
	}
	// normal replacement
	if *versionFile == "" {
		replace(content)
		return
	}
	// version replacement
	replaceVersion(content)
}

func replace(content string) {
	if len(flag.Args()) < 2 {
		fmt.Println("not enough args")
		return
	}
	oldStr := flag.Arg(0)
	newStr := flag.Arg(1)
	newContent := strings.ReplaceAll(content, oldStr, newStr)
	f, e := fileToolkit.OpenFileForWrite(*file)
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

func replaceVersion(content string) {
	version := readVersion(*versionFile)
	match := flag.Arg(0)
	tpl := flag.Arg(1)

	newContent := bytes.NewBufferString("")
	var e error
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, match) {
			t := template.New("name")
			t, e = t.Parse(tpl)
			if e != nil {
				log.Fatal(e)
			}
			buf := bytes.NewBufferString("")
			t.Execute(buf, version)
			newContent.WriteString(buf.String() + "\n")
			continue
		}
		newContent.WriteString(line + "\n")
	}
	e = fileToolkit.WriteFile(*file, []byte(strings.TrimSuffix(newContent.String(), "\n")))
	if e != nil {
		log.Fatal(e)
	}
}
func readVersion(path string) string {
	content, e := fileToolkit.ReadFileAll(path)
	if e != nil {
		log.Fatal(e)
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		str := strings.TrimSuffix(string(line), "\r")
		if strings.Contains(str, "Version = ") {
			return strToolkit.TrimBoth(strToolkit.SubAfterLast(str, "Version = ", str), `"`)
		}
	}
	log.Fatal("const Version = \"1.1.1\" , doesn't exists")
	return ""
}
