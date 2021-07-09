package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/StevenZack/tools/fileToolkit"
	"github.com/StevenZack/tools/strToolkit"
)

var (
	file        = flag.String("f", "", "specific your file")
	versionFile = flag.String("v", "", "from version.go file")
)

func main() {
	flag.Parse()
	// normal replacement
	if *versionFile != "" {
		replaceVersion()
		return
	}
	if *file != "" {
		singleFile(*file)
		return
	}
	replace()
	// version replacement
}

func replace() {
	if len(flag.Args()) < 2 {
		fmt.Println("not enough args")
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
		singleFile(path)
		return nil
	})
	if e != nil {
		log.Fatal(e)
	}
}

func singleFile(path string) {
	b, e := ioutil.ReadFile(path)
	if e != nil {
		log.Fatal(e)
	}

	oldStr := flag.Arg(0)
	newStr := flag.Arg(1)

	out := new(bytes.Buffer)
	changed := false
	ss := strings.Split(string(b), "\n")
	for i, line := range ss {
		if !changed && strings.Contains(line, oldStr) {
			out.WriteString(strings.Replace(line, oldStr, newStr, -1))
			changed = true
		} else {
			out.WriteString(line)
		}

		if i < len(ss)-1 {
			out.WriteString("\n")
		}
	}

	e = ioutil.WriteFile(path, out.Bytes(), 0644)
	if e != nil {
		log.Fatal(e)
	}

}

func replaceVersion() {
	content, e := fileToolkit.ReadFileAll(*file)
	if e != nil {
		log.Fatal(e)
	}
	version := readVersion(*versionFile)
	match := flag.Arg(0)
	tpl := flag.Arg(1)

	newContent := bytes.NewBufferString("")
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
