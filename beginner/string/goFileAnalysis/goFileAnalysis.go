package goFileAnalysis

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type Method struct {
	Name    string
	Params  string
	Results string
}

// 输出某go文件下所有的方法&参数
func GetGrpcFuncByFile(path string) []Method {
	fset := token.NewFileSet()
	// 这里的参数基本是原封不动的传给了scanner的Init函数
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	var attr *ast.GenDecl
	var ts *ast.TypeSpec
	var ok bool
	var nameSlice []Method
	for _, decl := range node.Decls {
		attr, ok = decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range attr.Specs {
			ts, ok = spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			var ty *ast.InterfaceType
			ty, ok = ts.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}
			if strings.HasPrefix(ts.Name.Name, "Unsafe") || strings.HasSuffix(ts.Name.Name, "Client") {
				continue
			}
			nameSlice = make([]Method, 0)
			var startExpr *ast.Ident
			for _, item := range ty.Methods.List {
				if strings.HasPrefix(item.Names[0].Name, "must") {
					continue
				}
				itemTy, ok := item.Type.(*ast.FuncType)
				if !ok {
					continue
				}
				method := Method{}
				method.Name = item.Names[0].Name
				startExpr, ok = itemTy.Params.List[1].Type.(*ast.StarExpr).X.(*ast.Ident)
				if !ok {
					continue
				}
				method.Params = startExpr.Name
				startExpr, ok = itemTy.Results.List[0].Type.(*ast.StarExpr).X.(*ast.Ident)
				if !ok {
					continue
				}
				method.Results = startExpr.Name
				nameSlice = append(nameSlice, method)
			}
		}
	}
	return nameSlice
}

// 输出某go文件下所有结构体名称
func GetStructNameByFile(file string) []string {
	fset := token.NewFileSet()
	// 这里的参数基本是原封不动的传给了scanner的Init函数
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	var attr *ast.GenDecl
	var ts *ast.TypeSpec
	var ok bool
	var nameMap = make(map[string]struct{})
	for _, decl := range node.Decls {
		attr, ok = decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range attr.Specs {
			ts, ok = spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if _, ok = ts.Type.(*ast.StructType); !ok {
				continue
			}
			nameMap[ts.Name.Name] = struct{}{}
		}
	}
	var result = make([]string, len(nameMap))
	var index int
	for k, _ := range nameMap {
		result[index] = k
		index++
	}
	return result
}
