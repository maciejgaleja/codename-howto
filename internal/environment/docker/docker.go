package docker

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Image struct {
	Tag string
}

func (i Image) Build(dockerfile string) (err error) {
	cmd := exec.Command("docker", "build", "-t", i.Tag, "-")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, dockerfile)
	}()
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("error during image build: %s: %v", out, err)
		return
	}

	return
}

func (i Image) Run() (c Container, err error) {
	cidfile := "./codename-howto-cid"
	defer os.Remove(cidfile)
	os.Remove(cidfile)
	cmd := exec.Command("docker", "run", "-d", "--cidfile", cidfile, "--rm", i.Tag, "/bin/sleep", "10")
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("error during starting: %s: %v", out, err)
		return
	}

	cid, err := os.ReadFile(cidfile)

	c.ID = string(cid)
	return
}

type Container struct {
	ID string
}

func (c Container) Stop() (err error) {
	cmd := exec.Command("docker", "stop", "-t", "1", c.ID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("error during stopping container: %s: %v", out, err)
	}
	return
}

func (c Container) Exec(interpreter string, input []byte) (o []byte, err error) {
	cmd := exec.Command("docker", "exec", "-i", c.ID, "/usr/bin/env", interpreter)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	go func() {
		defer stdin.Close()
		stdin.Write(input)
	}()
	o, err = cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("error during image build: %s: %v", o, err)
		return
	}

	return
}
