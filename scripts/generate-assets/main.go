package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/shurcooL/vfsgen"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("usage: go run scripts/generate-assets/main.go SRC DST")
	}

	src := os.Args[1]
	dst := os.Args[2]
	pkg := filepath.Base(filepath.Dir(dst))
	err := vfsgen.Generate(http.Dir(src), vfsgen.Options{
		Filename:     dst,
		PackageName:  pkg,
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
