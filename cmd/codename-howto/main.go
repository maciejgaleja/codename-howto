package main

import (
	"os"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"

	"fmt"
)

type Filename string
type Interpreter string

const (
	CreateFile = iota
	ExecuteCode
)

type Step struct {
	Action      int
	Interpreter Interpreter
	Code        []byte
}

type HowTo struct {
	Environment []byte
	Steps       []Step
}

func getHeaderText(h *ast.Heading) string {
	return string(h.Children[0].AsLeaf().Literal)
}

func getCode(b *ast.CodeBlock) (i Interpreter, c []byte) {
	i = Interpreter(b.Info)
	c = b.Literal
	return
}

func ParseMd(mdf Filename) (h HowTo, err error) {
	md, err := os.ReadFile(string(mdf))
	if err != nil {
		return
	}
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	var hed string

	if _, ok := doc.(*ast.Document); ok {
		for _, c := range doc.GetChildren() {
			switch c.(type) {
			case *ast.Heading:
				hed = getHeaderText(c.(*ast.Heading))
			case *ast.CodeBlock:
				interpreter, code := getCode(c.(*ast.CodeBlock))
				if hed == "" {
					err = fmt.Errorf("found code not in a header block")
					return
				}
				if strings.EqualFold(string(interpreter), "dockerfile") {
					if len(h.Environment) == 0 {
						h.Environment = code
					}
				} else {
					if strings.HasPrefix(hed, "Create file:") {
						h.Steps = append(h.Steps, Step{Action: CreateFile, Interpreter: "", Code: code})
					} else {
						h.Steps = append(h.Steps, Step{Action: ExecuteCode, Interpreter: interpreter, Code: code})
					}
				}
			}
		}
	} else {
		err = fmt.Errorf("error while parsing md")
		return
	}

	return
}

func main() {
	md, err := ParseMd("./examples/how-to-print.md")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(md.Environment))

	for _, s := range md.Steps {
		fmt.Println(s)
	}
}
