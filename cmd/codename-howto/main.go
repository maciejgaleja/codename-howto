package main

import (
	"os"
	"path"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/maciejgaleja/codename-howto/internal/environment/docker"
	"github.com/maciejgaleja/codename-howto/pkg/howto"

	"fmt"
)

var CLI struct {
	OutputDir  string   `type:"existingdir" default:"." help:"directory to which put resulting files"`
	InputFiles []string `arg:"" type:"existingfile"`
}

func getDelimiters(id string) (begin, end string) {
	begin = fmt.Sprintf("begin %s", id)
	end = fmt.Sprintf("end %s", id)
	return
}

func findFragment(begin, end string, text string) (ret string, err error) {
	tokens := strings.Split(text, begin)
	if len(tokens) != 2 {
		err = fmt.Errorf("problem while finding fragment")
	}
	ret = strings.Split(tokens[1], end)[0]
	ret = strings.TrimSpace(ret)
	return
}

func handleSingleFile(f string) {
	md, err := howto.ParseMd(howto.Filename(f))
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

	fmt.Printf("Created container: %s\n", container.ID)

	script := []byte{}
	for _, step := range md.Steps {
		script = append(script, []byte("\n")...)
		preScript := ""
		postScript := ""
		if step.OutputNode != nil {
			dB, dE := getDelimiters(fmt.Sprintf("%p", step.OutputNode))
			preScript = fmt.Sprintf("echo; echo %s", dB)
			postScript = fmt.Sprintf("echo; echo %s", dE)
		}
		script = append(script, []byte("\n"+preScript+"\n")...)
		script = append(script, step.Code...)
		script = append(script, []byte("\n"+postScript+"\n")...)
	}

	o, err := container.Exec("bash", script)
	fmt.Printf("%s\n", strings.TrimSpace(string(script)))
	fmt.Printf("%s\n", strings.TrimSpace(string(o)))
	if err != nil {
		panic(err)
	}

	for _, step := range md.Steps {
		if step.OutputNode != nil {
			dB, dE := getDelimiters(fmt.Sprintf("%p", step.OutputNode))
			frag, err := findFragment(dB, dE, string(o))
			if err != nil {
				panic(err)
			}
			step.OutputNode.Literal = []byte(frag)
		}
	}

	filename := path.Base(f)
	if err := os.WriteFile(path.Join(CLI.OutputDir, filename), md.AsMarkdown(), 0644); err != nil {
		panic(err)
	}
}

func main() {
	_ = kong.Parse(&CLI)
	for _, filepath := range CLI.InputFiles {
		handleSingleFile(filepath)
	}
}
