//go:build go1.18

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/format"
	"go/parser"
	"go/token"
	"html"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bradenaw/juniper/xerrors"
	"github.com/bradenaw/juniper/slices"
)

func main() {
	outDir := "docs/"
	flag.StringVar(&outDir, "out_dir", outDir, "The directory to write output to.")
	flag.Parse()

	err := func() error {
		modBasePath, err := modPath()
		if err != nil {
			return err
		}
		packages, err := modPackages()
		if err != nil {
			return err
		}

		packages = slices.Filter(packages, func(packagePath string) bool  {
			return !strings.Contains(packagePath, "/internal/")
		})

		err = os.MkdirAll(outDir, os.ModeDir|0755)
		if err != nil {
			return xerrors.WithStack(err)
		}
		indexOut, err := os.Create(filepath.Join(outDir, "index.md"))
		if err != nil {
			return xerrors.WithStack(err)
		}

		for _, pkg := range packages {
			pkgRel := strings.TrimPrefix(pkg, modBasePath +"/")
			err := os.MkdirAll(
				filepath.Join(outDir, filepath.Dir(pkgRel)),
				os.ModeDir|0755,
			)
			if err != nil {
				return xerrors.WithStack(err)
			}
			out, err := os.Create(filepath.Join(outDir, pkgRel + ".md"))
			if err != nil {
				return xerrors.WithStack(err)
			}
			err = writePackageDoc(out, modBasePath, pkg)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(
				indexOut,
				"<samp><a href=\"%s.html\">%s</a></samp>\n\n",
				pkgRel,
				pkgRel,
			)
			if err != nil {
				return err
			}
		}
		return nil
	}()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func goBin() string {
	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		return "go"
	}
	return goroot + "/bin/go"
}

func modPath() (string, error) {
	type output struct {
		Module struct {
			Path string
		}
	}

	b, err := exec.Command(goBin(), "mod", "edit", "--json").Output()
	if err != nil {
		return "", xerrors.WithStack(err)
	}
	var o output
	err = json.Unmarshal(b, &o)
	if err != nil {
		return "", err
	}
	return o.Module.Path, nil
}

func modPackages() ([]string, error) {
	b, err := exec.Command(goBin(), "list", "./...").Output()
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	return strings.Split(strings.TrimSpace(string(b)), "\n"), nil
}

type captureErrWriter struct {
	inner io.Writer
	err   error
}

func (w *captureErrWriter) Write(b []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	n, err := w.inner.Write(b)
	if err != nil {
		w.err = err
	}
	return n, err
}

func writePackageDoc(out io.Writer, modulePath string, packageImportPath string) (retErr error) {
	cOut := captureErrWriter{inner: out}
	out = &cOut
	defer func() {
		if retErr == nil {
			retErr = cOut.err
		}
	}()

	pkg, err := build.Default.ImportDir(
		strings.TrimPrefix(packageImportPath, modulePath + "/"),
		build.ImportComment,
	)
	if err != nil {
		return xerrors.WithStack(err)
	}
	pkg.ImportPath = packageImportPath

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
		return xerrors.WithStack(err)
	}
	var astFiles []*ast.File
	for _, astPkg := range pkgs {
		for _, astFile := range astPkg.Files {
			astFiles = append(astFiles, astFile)
		}
	}
	docPkg, err := doc.NewFromFiles(fset, astFiles, pkg.ImportPath, doc.PreserveAST)
	if err != nil {
		return xerrors.WithStack(err)
	}

	imports := importNames(docPkg.Imports)
	localSymbols := localSymbolLinks(docPkg)

	fmt.Fprintln(out, "# `package "+docPkg.Name+"`")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "```")
	fmt.Fprintln(out, "import \""+docPkg.ImportPath+"\"")
	fmt.Fprintln(out, "```")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "## Overview")
	fmt.Fprintln(out)
	fmt.Fprintln(out, docPkg.Doc)

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

	fmt.Fprintln(out)
	fmt.Fprintln(out, "## Index")
	fmt.Fprintln(out)
	for _, func_ := range docPkg.Funcs {
		if !token.IsExported(func_.Name) {
			continue
		}
		l, err := funcLink(func_)
		if err != nil {
			return err
		}
		fmt.Fprintln(out, "<samp>"+l+"</samp>")
		fmt.Fprintln(out)
	}
	for _, type_ := range docPkg.Types {
		if !token.IsExported(type_.Name) {
			continue
		}
		fmt.Fprint(out, "<samp><a href=\"")
		fmt.Fprint(out, localSymbolLink(type_.Name))
		fmt.Fprint(out, "\">type ")
		fmt.Fprint(out, type_.Name)
		fmt.Fprintln(out, "</a></samp>")
		fmt.Fprintln(out)

		for _, funcs := range [][]*doc.Func{type_.Funcs, type_.Methods} {
			for _, func_ := range funcs {
				l, err := funcLink(func_)
				if err != nil {
					return err
				}
				fmt.Fprintln(out, "<samp>"+indent(l, 4)+"</samp>")
				fmt.Fprintln(out)
			}
		}
	}
	fmt.Fprintln(out)

	fmt.Fprintln(out, "## Constants")
	fmt.Fprintln(out)
	if len(docPkg.Consts) > 0 {
		fmt.Fprintln(out, "<pre>")
		for _, const_ := range docPkg.Consts {
			for _, name := range const_.Names {
				fmt.Fprint(out, "<a id=\"")
				fmt.Fprint(out, name)
				fmt.Fprint(out, "\"></a>")
			}
			fmt.Fprint(out, strWithLinks(fset, modulePath, docPkg.ImportPath, imports, localSymbols, const_.Decl))
			fmt.Fprintln(out)
		}
		fmt.Fprintln(out, "</pre>")
	} else {
		fmt.Fprintln(out, "This section is empty.")
	}
	fmt.Fprintln(out)

	fmt.Fprintln(out, "## Variables")
	fmt.Fprintln(out)
	if len(docPkg.Vars) > 0 {
		fmt.Fprintln(out, "<pre>")
		for _, var_ := range docPkg.Vars {
			for _, name := range var_.Names {
				fmt.Fprint(out, "<a id=\"")
				fmt.Fprint(out, name)
				fmt.Fprint(out, "\"></a>")
			}
			fmt.Fprint(out, strWithLinks(fset, modulePath, docPkg.ImportPath, imports, localSymbols, var_.Decl))
			fmt.Fprintln(out)
		}
		fmt.Fprintln(out, "</pre>")
	} else {
		fmt.Fprintln(out, "This section is empty.")
	}
	fmt.Fprintln(out)

	fmt.Fprintln(out, "## Functions")
	fmt.Fprintln(out)
	for _, func_ := range docPkg.Funcs {
		if !token.IsExported(func_.Name) {
			continue
		}
		printFunc(out, fset, modulePath, docPkg.ImportPath, imports, localSymbols, func_)
	}

	fmt.Fprintln(out, "## Types")
	fmt.Fprintln(out)
	for _, type_ := range docPkg.Types {
		if !token.IsExported(type_.Name) {
			continue
		}
		fmt.Fprint(out, "<h3><a id=\"")
		fmt.Fprint(out, type_.Name)
		fmt.Fprint(out, "\"></a><samp>type ")
		fmt.Fprint(out, type_.Name)
		fmt.Fprintln(out, "</samp></h3>")
		type_.Decl.Doc = nil
		fmt.Fprintln(out, "```go")
		err := format.Node(out, fset, type_.Decl)
		if err != nil {
			return err
		}
		fmt.Fprintln(out)
		fmt.Fprintln(out, "```")
		fmt.Fprintln(out)
		fmt.Fprintln(out, type_.Doc)
		fmt.Fprintln(out)
		for _, funcs := range [][]*doc.Func{type_.Funcs, type_.Methods} {
			for _, func_ := range funcs {
				if !token.IsExported(func_.Name) {
					continue
				}
				printFunc(out, fset, modulePath, docPkg.ImportPath, imports, localSymbols, func_)
			}
		}
	}

	return nil
}

func printFunc(
	out io.Writer,
	fset *token.FileSet,
	modulePath string,
	packagePath string,
	imports map[string]string,
	localSymbols map[string]string,
	func_ *doc.Func,
) {
	fmt.Fprint(out, "<h3><a id=\"")
	fmt.Fprint(out, func_.Name)
	fmt.Fprint(out, "\"></a><samp>")
	fmt.Fprint(out, strWithLinks(fset, modulePath, packagePath, imports, localSymbols, func_.Decl))
	fmt.Fprintln(out, "</samp></h3>")
	fmt.Fprintln(out)
	fmt.Fprintln(out, func_.Doc)
	fmt.Fprintln(out)

	for _, example := range func_.Examples {
		printExample(out, fset, example)
	}
}

func printExample(
	out io.Writer,
	fset *token.FileSet,
	example *doc.Example,
) {
	fmt.Fprintln(out, "#### Example "+example.Suffix)
	fmt.Fprintln(out, "```go")
	err := format.Node(out, fset, example.Code)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, "```")
	fmt.Fprintln(out)
	if !example.EmptyOutput {
		if example.Unordered {
			fmt.Fprintln(out, "Unordered output:")
		} else {
			fmt.Fprintln(out, "Output:")
		}
		fmt.Fprintln(out, "```text")
		fmt.Fprint(out, example.Output)
		fmt.Fprintln(out, "```")
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

func nonLocalSymbolLink(
	modulePath string,
	packagePath string,
	importPath string,
	symbol string,
) string {
	if !strings.HasPrefix(importPath, modulePath+"/") {
		return "https://pkg.go.dev/" + importPath + "#" + symbol
	}

	s, err := filepath.Rel(filepath.Dir(packagePath), filepath.Dir(importPath))
	if err != nil {
		panic("")
	}
	if !strings.HasSuffix(s, "/") {
		s += "/"
	}
	return s + filepath.Base(importPath) + ".html#" + symbol
}

func strWithLinks(
	fset *token.FileSet,
	modulePath string,
	packagePath string,
	imports map[string]string,
	localSymbols map[string]string,
	node ast.Node,
) string {
	var sb strings.Builder
	var visit func(node ast.Node)

	visitFuncParams := func(node *ast.FuncType) {
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
	}
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
			visitFuncParams(node.Type)
		case *ast.FuncType:
			sb.WriteString("func")
			visitFuncParams(node)
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
					sb.WriteString(nonLocalSymbolLink(
						modulePath,
						packagePath,
						importPath,
						node.Sel.Name,
					))
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
