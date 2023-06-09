package howto

import (
	"fmt"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/md"
	"github.com/gomarkdown/markdown/parser"
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

func isOutputPlaceholder(c *ast.CodeBlock) bool {
	return strings.TrimSpace(string(c.Literal)) == "<output placeholder>"
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
						if isOutputPlaceholder(c.(*ast.CodeBlock)) {
							if len(h.Steps) == 0 {
								err = fmt.Errorf("otput placeholder found without previous code")
								return
							}
							h.Steps[len(h.Steps)-1].OutputNode = c.(*ast.CodeBlock)
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
