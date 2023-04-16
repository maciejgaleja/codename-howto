package main

import (
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/md"
	"github.com/gomarkdown/markdown/parser"
	"github.com/maciejgaleja/codename-howto/internal/environment/docker"

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
	OutputNode  *ast.CodeBlock
}

type HowTo struct {
	Environment []byte
	Steps       []Step
	Doc         ast.Node
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
	h.Doc = p.Parse(md)

	var hed string

	if _, ok := h.Doc.(*ast.Document); ok {
		for _, c := range h.Doc.GetChildren() {
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
						h.Steps = append(h.Steps, Step{Action: CreateFile, Interpreter: "", Code: code, OutputNode: nil})
					} else {
						output, _ := c.(*ast.CodeBlock)
						if strings.TrimSpace(string(output.Literal)) == "<output placeholder>" {
							if len(h.Steps) == 0 {
								err = fmt.Errorf("otput placeholder found without previous code")
								return
							}
							h.Steps[len(h.Steps)-1].OutputNode = output
						} else {
							h.Steps = append(h.Steps, Step{Action: ExecuteCode, Interpreter: interpreter, Code: code, OutputNode: nil})
						}
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

func (h HowTo) AsMarkdown() (ret []byte) {
	renderer := md.NewRenderer()
	ret = markdown.Render(h.Doc, renderer)
	return
}

func main() {
	md, err := ParseMd("./examples/how-to-print.md")
	if err != nil {
		panic(err)
	}

	dockerfile := string(md.Environment)
	i := docker.Image{Tag: "codename-howto"}
	err = i.Build(dockerfile)
	if err != nil {
		panic(err)
	}

	container, err := i.Run()
	if err != nil {
		panic(err)
	}
	defer container.Stop()

	fmt.Println(container)

	for _, step := range md.Steps {
		fmt.Println(step)
		o, err := container.Exec(string(step.Interpreter), step.Code)
		if err != nil {
			panic(err)
		}
		fmt.Println("!!! ", string(o))
		if step.OutputNode != nil {
			step.OutputNode.Literal = o
		}
	}

	fmt.Println(string(md.AsMarkdown()))
}
