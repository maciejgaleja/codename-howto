package main

import (
	"github.com/maciejgaleja/codename-howto/internal/environment/docker"
	"github.com/maciejgaleja/codename-howto/pkg/howto"

	"fmt"
)

func main() {
	md, err := howto.ParseMd("./examples/how-to-print.md")
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
