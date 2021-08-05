// +build trace
package functrace

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var mu sync.Mutex
var m = make(map[uint64]int)

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func printTrace(id uint64, file, name, typ string, indent int, line int) {
	indents := ""
	for i := 0; i < indent; i++ {
		indents += "\t"
	}
	fmt.Printf("g[%02d]:%s %s %s [%s:%d]\n", id, indents, typ, name, file, line)
}

func Trace() func() {
	//取二层的文件路径和行号
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		panic("not found caller")
	}

	//取一层的函数名
	pc, _, _, _ = runtime.Caller(1)
	if !ok {
		panic("not found caller")
	}
	id := getGID()
	fn := runtime.FuncForPC(pc)
	name := fn.Name()

	mu.Lock()
	v := m[id]
	m[id] = v + 1
	mu.Unlock()

	//parse last file
	subPaths := strings.Split(file, "/")
	subPathLength := len(subPaths)
	lastFileIndex := 0
	if subPathLength > 0 {
		lastFileIndex = subPathLength - 1
	}
	printTrace(id, subPaths[lastFileIndex], name, "->", v+1, line)
	return func() {
		mu.Lock()
		v := m[id]
		m[id] = v - 1
		mu.Unlock()
		printTrace(id, subPaths[lastFileIndex], name, "<-", v, line)
	}
}
