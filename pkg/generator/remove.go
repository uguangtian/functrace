package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

func Remove(filename string) ([]byte, error) {
	fset := token.NewFileSet()
	oldAST, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", filename, err)
	}
	//fmt.Printf("%#v\n", *oldAST)

	//没有函数
	if !hasFuncDecl(oldAST) {
		return nil, nil
	}

	// add import declaration
	astutil.DeleteImport(fset, oldAST, "github.com/uguangtian/functrace")
	//fmt.Printf("added=%#v\n", added)

	// inject code into each function declaration
	removeDeferTraceIntoFuncDecls(oldAST)

	buf := &bytes.Buffer{}
	err = format.Node(buf, fset, oldAST)
	if err != nil {
		return nil, fmt.Errorf("error formatting new code: %w", err)
	}
	return buf.Bytes(), nil
}
func removeDeferTraceIntoFuncDecls(f *ast.File) {
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if ok {
			// inject code to fd
			removeDeferStmt(fd)
		}
	}
}
func removeDeferStmt(fd *ast.FuncDecl) (removed bool) {
	stmts := fd.Body.List

	// check whether "defer functrace.Trace()()" has already exists
	for _, stmt := range stmts {
		ds, ok := stmt.(*ast.DeferStmt)
		if !ok {
			continue
		}
		// it is a defer stmt
		ce, ok := ds.Call.Fun.(*ast.CallExpr)
		if !ok {
			continue
		}

		se, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		x, ok := se.X.(*ast.Ident)
		if !ok {
			continue
		}
		if (x.Name == "functrace") && (se.Sel.Name == "Trace") {
			// already exist , remove
			newList := make([]ast.Stmt, len(stmts)-1)
			//第一个元素stmt?
			copy(newList, stmts[1:])
			fd.Body.List = newList
			removed = true
			return
		}
	}

	// not found "defer functrace.Trace()()"
	return
}
