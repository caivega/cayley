package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	packageName = flag.String("pck", "github.com/caivega/cayley/query/gizmo", "")
	out         = flag.String("o", "-", "output file")
	in          = flag.String("i", "", "input file")
)

const placeholder = `#AUTOGENERATED#`

func main() {
	flag.Parse()

	path := filepath.Join(os.Getenv("GOPATH"), "src", *packageName)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	p := pkgs[filepath.Base(*packageName)]

	dp := doc.New(p, *packageName, doc.AllDecls)

	var w io.Writer = os.Stdout
	if fname := *out; fname != "" && fname != "-" {
		f, err := os.Create(fname)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		w = f
	}
	var r io.Reader = strings.NewReader(placeholder)
	if fname := *in; fname != "" {
		f, err := os.Open(fname)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		r = f
	}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := bytes.TrimSpace(sc.Bytes())
		if bytes.Equal(line, []byte(placeholder)) {
			writeDocs(w, dp)
		} else {
			w.Write(line)
			w.Write([]byte("\n"))
		}
	}
}

func writeDocs(w io.Writer, dp *doc.Package) {
	type Type struct {
		Title string
		Name  string
	}

	names := map[string]Type{
		"graphObject": {
			Title: "The `graph` object",
			Name:  "graph",
		},
		"pathObject": {
			Title: "Path object",
			Name:  "path",
		},
	}
	for _, tp := range dp.Types {
		t, ok := names[tp.Name]
		if !ok {
			continue
		}
		s := tp.Doc
		if i := strings.IndexAny(s, "\n\r"); i >= 0 {
			s = s[i+1:]
		}
		s = strings.TrimSpace(s)
		fmt.Fprintf(w, "## %s\n\n", t.Title)
		fmt.Fprintf(w, "%s\n\n", funcDocs(s))
		for _, m := range tp.Methods {
			if !isExported(m.Name) {
				continue
			}
			m.Doc = strings.TrimSpace(m.Doc)
			sig := Signature(m)
			fmt.Fprintf(w, "### `%s.%s%s`\n\n%s\n\n", t.Name, m.Name, sig, funcDocs(m.Doc))
		}
	}
}

var reSignature = regexp.MustCompile(`Signature:\s+\((.+)\)`)

func Signature(m *doc.Func) string {
	if reSignature.MatchString(m.Doc) {
		sub := reSignature.FindStringSubmatch(m.Doc)
		m.Doc = strings.Replace(m.Doc, sub[0], "", 1)
		return "(" + sub[1] + ")"
	}
	tp := m.Decl.Type
	if isJsArgs(tp.Params) {
		return "(*)"
	}
	var names []string
	for _, a := range tp.Params.List {
		for _, name := range a.Names {
			names = append(names, name.Name)
		}
	}
	buf := bytes.NewBuffer(nil)
	buf.WriteRune('(')
	buf.WriteString(strings.Join(names, ", "))
	buf.WriteRune(')')
	return buf.String()
}

func isExported(s string) bool {
	return ast.IsExported(s)
}

func isJsArgs(f *ast.FieldList) bool {
	if len(f.List) != 1 {
		return false
	}
	p := f.List[0]
	if len(p.Names) != 1 {
		return false
	}
	sel, ok := p.Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	return sel.Sel.Name == "FunctionCall"
}

var reScript = regexp.MustCompile(`//\s*(\w+)`)

func funcDocs(s string) string {
	if s == "" {
		return "TODO: docs"
	}
	buf := bytes.NewBuffer(nil)
	buf.Grow(len(s))
	sc := bufio.NewScanner(strings.NewReader(s))
	const defaultLang = ""
	var (
		inCode bool
		lang   string
	)
	for sc.Scan() {
		line := sc.Text()
		if code := strings.HasPrefix(line, "\t"); code {
			if !inCode {
				inCode = true
				lang = defaultLang
				skip := false
				if reScript.MatchString(line) {
					skip = true
					lang = reScript.FindStringSubmatch(line)[1]
				}
				buf.WriteString("```")
				buf.WriteString(lang)
				buf.WriteString("\n")
				if skip {
					continue
				}
			}
			line = strings.TrimPrefix(line, "\t")
		} else if inCode && !code {
			inCode = false
			buf.WriteString("```\n")
		}
		buf.WriteString(line)
		buf.WriteRune('\n')
	}
	if inCode {
		buf.WriteString("```\n")
	}
	return buf.String()
}
