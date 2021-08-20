package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/uguangtian/functrace/pkg/generator"
)

var (
	wrote   bool
	unset   bool
	version bool
)

func init() {

	// -unset -w file.go
	// flag.BoolVar(&unset, "unset ", false, "reset functrace flag")

	// -w file.go
	// flag.BoolVar(&wrote, "w", false, "write result to (source) file instead of stdout")

}

func usage() {

	fmt.Println("gen [-w] xxx.go")
	fmt.Println("gen [-w] [-unset] xxx.go")
	// flag.PrintDefaults()
}

func Init() {

	// 当为false时,设-v则为true,当为true时，设-v仍为true

	flag.BoolVar(&unset, "unset", false, "reset functrace")
	flag.BoolVar(&wrote, "w", false, "write result to (source) file instead of stdout")
	flag.Parse()

}

func main() {

	fmt.Println(os.Args)
	Init()
	flag.Parse()

	flag.Usage = usage

	if len(os.Args) < 2 {
		usage()
		return
	}

	var file string
	if len(os.Args) == 3 {
		file = os.Args[2]
	}

	if len(os.Args) == 4 {
		file = os.Args[3]
	}

	if len(os.Args) == 2 {
		file = os.Args[1]
	}

	if filepath.Ext(file) != ".go" {
		usage()
		return
	}
	var (

		//新文件缓存
		newSrc []byte
		err    error
	)
	if unset {
		newSrc, err = generator.Remove(file)
		if err != nil {
			panic(err)
		}
	} else {
		newSrc, err = generator.Rewrite(file)
		if err != nil {
			panic(err)
		}
	}

	if newSrc == nil {
		// add nothing to the source file. no change
		fmt.Printf("no trace added for %s\n", file)
		return
	}

	if !wrote {
		fmt.Println(string(newSrc))
		return
	}

	// write to the source file
	if err = ioutil.WriteFile(file, newSrc, 0666); err != nil {
		fmt.Printf("write %s error: %v\n", file, err)
		return
	}
	fmt.Printf("add trace for %s ok\n", file)
}
