//go:build go1.18

package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/format"
	"go/parser"
	"go/token"
	"html"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	err := main2(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func main2(packageDir string) error {
	pkg, err := build.Default.ImportDir(packageDir, build.ImportComment)
	if err != nil {
		return err
	}
	pkg.ImportPath = "github.com/bradenaw/juniper/" + packageDir

	include := func(info fs.FileInfo) bool {
		for _, name := range pkg.GoFiles {
			if name == info.Name() {
				return true
			}
		}
		for _, name := range pkg.TestGoFiles {
			if name == info.Name() {
				return true
			}
		}
		for _, name := range pkg.XTestGoFiles {
			if name == info.Name() {
				return true
			}
		}
		for _, name := range pkg.CgoFiles {
			if name == info.Name() {
				return true
			}
		}
		return false
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkg.Dir, include, parser.ParseComments)
	if err != nil {
		return err
	}
	var astFiles []*ast.File
	for _, astPkg := range pkgs {
		for _, astFile := range astPkg.Files {
			astFiles = append(astFiles, astFile)
		}
	}
	docPkg, err := doc.NewFromFiles(fset, astFiles, pkg.ImportPath, doc.PreserveAST)
	if err != nil {
		return err
	}

	imports := importNames(docPkg.Imports)
	localSymbols := localSymbolLinks(docPkg)

	fmt.Println("# `package " + docPkg.Name + "`")
	fmt.Println()
	fmt.Println("```")
	fmt.Println("import \"" + docPkg.ImportPath + "\"")
	fmt.Println("```")
	fmt.Println()
	fmt.Println("## Overview")
	fmt.Println()
	fmt.Println(docPkg.Doc)

	funcLink := func(func_ *doc.Func) (string, error) {
		strippedDecl := &ast.FuncDecl{
			Recv: func_.Decl.Recv,
			Name: func_.Decl.Name,
			Type: func_.Decl.Type,
			// Leave off Doc and Body.
		}

		var sb strings.Builder
		sb.WriteString("<a href=\"")
		sb.WriteString(localSymbolLink(func_.Name))
		sb.WriteString("\">")

		var innerSB strings.Builder
		err := format.Node(&innerSB, fset, strippedDecl)
		if err != nil {
			return "", err
		}
		sb.WriteString(html.EscapeString(innerSB.String()))
		sb.WriteString("</a>")
		return sb.String(), nil
	}

	fmt.Println()
	fmt.Println("## Index")
	fmt.Println()
	for _, func_ := range docPkg.Funcs {
		if !token.IsExported(func_.Name) {
			continue
		}
		l, err := funcLink(func_)
		if err != nil {
			return err
		}
		fmt.Println("<samp>" + l + "</samp>")
		fmt.Println()
	}
	for _, type_ := range docPkg.Types {
		if !token.IsExported(type_.Name) {
			continue
		}
		fmt.Print("<samp><a href=\"")
		fmt.Print(localSymbolLink(type_.Name))
		fmt.Print("\">type ")
		fmt.Print(type_.Name)
		fmt.Println("</a></samp>")
		fmt.Println()

		for _, funcs := range [][]*doc.Func{type_.Funcs, type_.Methods} {
			for _, func_ := range funcs {
				l, err := funcLink(func_)
				if err != nil {
					return err
				}
				fmt.Println("<samp>" + indent(l, 4) + "</samp>")
				fmt.Println()
			}
		}
	}
	fmt.Println()

	fmt.Println("## Constants")
	fmt.Println()
	if len(docPkg.Consts) > 0 {
		fmt.Println("<pre>")
		for _, const_ := range docPkg.Consts {
			for _, name := range const_.Names {
				fmt.Print("<a id=\"")
				fmt.Print(name)
				fmt.Print("\"></a>")
			}
			fmt.Print(strWithLinks(fset, docPkg.ImportPath, imports, localSymbols, const_.Decl))
			fmt.Println()
		}
		fmt.Println("</pre>")
	} else {
		fmt.Println("This section is empty.")
	}
	fmt.Println()

	fmt.Println("## Variables")
	fmt.Println()
	if len(docPkg.Vars) > 0 {
		fmt.Println("<pre>")
		for _, var_ := range docPkg.Vars {
			for _, name := range var_.Names {
				fmt.Print("<a id=\"")
				fmt.Print(name)
				fmt.Print("\"></a>")
			}
			fmt.Print(strWithLinks(fset, docPkg.ImportPath, imports, localSymbols, var_.Decl))
			fmt.Println()
		}
		fmt.Println("</pre>")
	} else {
		fmt.Println("This section is empty.")
	}
	fmt.Println()

	fmt.Println("## Functions")
	fmt.Println()
	for _, func_ := range docPkg.Funcs {
		if !token.IsExported(func_.Name) {
			continue
		}
		printFunc(fset, docPkg.ImportPath, imports, localSymbols, func_)
	}

	fmt.Println("## Types")
	fmt.Println()
	for _, type_ := range docPkg.Types {
		if !token.IsExported(type_.Name) {
			continue
		}
		fmt.Print("<h3><a id=\"")
		fmt.Print(type_.Name)
		fmt.Print("\"></a><samp>type ")
		fmt.Print(type_.Name)
		fmt.Println("</samp></h3>")
		type_.Decl.Doc = nil
		fmt.Println("```go")
		err := format.Node(os.Stdout, fset, type_.Decl)
		if err != nil {
			return err
		}
		fmt.Println()
		fmt.Println("```")
		fmt.Println()
		fmt.Println(type_.Doc)
		fmt.Println()
		for _, funcs := range [][]*doc.Func{type_.Funcs, type_.Methods} {
			for _, func_ := range funcs {
				if !token.IsExported(func_.Name) {
					continue
				}
				printFunc(fset, docPkg.ImportPath, imports, localSymbols, func_)
			}
		}
	}

	return nil
}

func printFunc(
	fset *token.FileSet,
	importPath string,
	imports map[string]string,
	localSymbols map[string]string,
	func_ *doc.Func,
) {
	fmt.Print("<h3><a id=\"")
	fmt.Print(func_.Name)
	fmt.Print("\"></a><samp>")
	fmt.Print(strWithLinks(fset, importPath, imports, localSymbols, func_.Decl))
	fmt.Println("</samp></h3>")
	fmt.Println()
	fmt.Println(func_.Doc)
	fmt.Println()

	for _, example := range func_.Examples {
		printExample(fset, example)
	}
}

func printExample(
	fset *token.FileSet,
	example *doc.Example,
) {
	fmt.Println("#### Example " + example.Suffix)
	fmt.Println("```go")
	err := format.Node(os.Stdout, fset, example.Code)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println("```")
	fmt.Println()
	if !example.EmptyOutput {
		if example.Unordered {
			fmt.Println("Unordered output:")
		} else {
			fmt.Println("Output:")
		}
		fmt.Println("```text")
		fmt.Print(example.Output)
		fmt.Println("```")
	}
}

func receiver(fset *token.FileSet, func_ *doc.Func) (string, error) {
	var sb strings.Builder
	err := format.Node(&sb, fset, func_.Decl.Recv.List[0].Type)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func importNames(imports []string) map[string]string {
	out := make(map[string]string)
	for _, import_ := range imports {
		idx := strings.LastIndex(import_, "/")
		if idx == 0 {
			out[import_] = import_
		} else {
			out[import_[idx+1:]] = import_
		}
	}
	return out
}

func localSymbolLinks(pkg *doc.Package) map[string]string {
	out := make(map[string]string)
	for _, type_ := range pkg.Types {
		out[type_.Name] = localSymbolLink(type_.Name)
		// TODO methods
	}
	for _, func_ := range pkg.Funcs {
		out[func_.Name] = localSymbolLink(func_.Name)
	}
	// TODO consts, vars
	return out
}

func localSymbolLink(symbol string) string {
	return "#" + symbol
}

func nonLocalSymbolLink(packagePath string, importPath string, symbol string) string {
	if !strings.HasPrefix(importPath, "github.com/bradenaw/juniper") {
		return "https://pkg.go.dev/" + importPath + "#" + symbol
	}

	s, err := filepath.Rel(filepath.Dir(packagePath), filepath.Dir(importPath))
	if err != nil {
		panic("")
	}
	if !strings.HasSuffix(s, "/") {
		s += "/"
	}
	return s + filepath.Base(importPath) + ".md#" + symbol
}

func strWithLinks(
	fset *token.FileSet,
	packagePath string,
	imports map[string]string,
	localSymbols map[string]string,
	node ast.Node,
) string {
	var sb strings.Builder
	var visit func(node ast.Node)
	visit = func(node ast.Node) {
		switch node := node.(type) {
		case *ast.ArrayType:
			sb.WriteString("[")
			if node.Len != nil {
				visit(node.Len)
			}
			sb.WriteString("]")
			visit(node.Elt)
		case *ast.ChanType:
			switch node.Dir {
			case ast.SEND:
				sb.WriteString("chan&lt;-")
			case ast.RECV:
				sb.WriteString("&lt;-chan")
			default:
				sb.WriteString("chan")
			}
			sb.WriteString(" ")
			visit(node.Value)
		case *ast.Ellipsis:
			sb.WriteString("...")
		case *ast.FuncDecl:
			sb.WriteString("func ")
			if node.Recv != nil {
				sb.WriteString("(")
				visit(node.Recv)
				sb.WriteString(") ")
			}
			visit(node.Name)
			visit(node.Type)
		case *ast.FuncType:
			if node.TypeParams != nil {
				sb.WriteString("[")
				visit(node.TypeParams)
				sb.WriteString("]")
			}
			sb.WriteString("(")
			visit(node.Params)
			sb.WriteString(")")
			if node.Results != nil {
				sb.WriteString(" ")
				needParens := len(node.Results.List) > 1 ||
					(len(node.Results.List) == 1 && len(node.Results.List[0].Names) > 1)
				if needParens {
					sb.WriteString("(")
				}
				visit(node.Results)
				if needParens {
					sb.WriteString(")")
				}
			}
		case *ast.Field:
			for i, name := range node.Names {
				if i != 0 {
					sb.WriteString(", ")
				}
				visit(name)
			}
			if len(node.Names) > 0 {
				sb.WriteString(" ")
			}
			visit(node.Type)
		case *ast.FieldList:
			for i, field := range node.List {
				if i != 0 {
					sb.WriteString(", ")
				}
				visit(field)
			}
		case *ast.GenDecl:
			switch node.Tok {
			case token.CONST:
				sb.WriteString("const ")
			case token.VAR:
				sb.WriteString("var ")
			default:
				panic(fmt.Sprintf("unrecognized token %s", node.Tok))
			}
			if len(node.Specs) > 1 {
				sb.WriteString("(\n")
			}
			for _, spec := range node.Specs {
				if len(node.Specs) > 1 {
					sb.WriteString("    ")
				}
				visit(spec)
				if len(node.Specs) > 1 {
					sb.WriteString("\n")
				}
			}
			if len(node.Specs) > 1 {
				sb.WriteString(")")
			}
		case *ast.Ident:
			link, ok := localSymbols[node.Name]
			if ok {
				sb.WriteString("<a href=\"")
				sb.WriteString(link)
				sb.WriteString("\">")
				sb.WriteString(node.Name)
				sb.WriteString("</a>")
			} else {
				sb.WriteString(node.Name)
			}
		case *ast.IndexExpr:
			visit(node.X)
			sb.WriteString("[")
			visit(node.Index)
			sb.WriteString("]")
		case *ast.IndexListExpr:
			visit(node.X)
			sb.WriteString("[")
			for i, index := range node.Indices {
				if i != 0 {
					sb.WriteString(", ")
				}
				visit(index)
			}
			sb.WriteString("]")
		case *ast.MapType:
			sb.WriteString("map[")
			visit(node.Key)
			sb.WriteString("]")
			visit(node.Value)
		case *ast.SelectorExpr:
			if ident, ok := node.X.(*ast.Ident); ok {
				importPath, ok := imports[ident.Name]
				if ok {
					sb.WriteString("<a href=\"")
					sb.WriteString(nonLocalSymbolLink(packagePath, importPath, node.Sel.Name))
					sb.WriteString("\">")
					sb.WriteString(ident.Name)
					sb.WriteString(".")
					sb.WriteString(node.Sel.Name)
					sb.WriteString("</a>")
					return
				}
			}
			visit(node.X)
			sb.WriteString(".")
			visit(node.Sel)
		case *ast.StarExpr:
			sb.WriteString("*")
			visit(node.X)
		case *ast.ValueSpec:
			for i, name := range node.Names {
				if i != 0 {
					sb.WriteString(", ")
				}
				visit(name)
			}
			if node.Type != nil {
				sb.WriteString(" ")
				visit(node.Type)
			}
			sb.WriteString(" = ")
			for i, value := range node.Values {
				if i != 0 {
					sb.WriteString(", ")
				}
				err := format.Node(&sb, fset, value)
				if err != nil {
					panic(err)
				}
			}
		default:
			panic(fmt.Sprintf("unrecognized *ast.Node %T after genning %s", node, sb.String()))
		}
	}
	visit(node)
	return sb.String()
}

func indent(s string, by int) string {
	indentation := strings.Repeat("&nbsp;", by)
	return indentation + strings.ReplaceAll(s, "\n", "\n"+indentation)
}
