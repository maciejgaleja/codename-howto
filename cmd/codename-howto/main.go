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

func handleSingleFile(f string) {
	md, err := howto.ParseMd(howto.Filename(CLI.InputFiles[0]))
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

	for _, step := range md.Steps {
		fmt.Printf("%s >> %s\n", step.Interpreter, strings.TrimSpace(string(step.Code)))
		o, err := container.Exec(string(step.Interpreter), step.Code)
		fmt.Printf("%s\n", strings.TrimSpace(string(o)))
		if err != nil {
			panic(err)
		}
		if step.OutputNode != nil {
			step.OutputNode.Literal = o
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
