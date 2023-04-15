package docker_test

import (
	"testing"

	"github.com/maciejgaleja/codename-howto/internal/environment/docker"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	i := docker.Image{"test-image"}

	df := `FROM ubuntu:latest

	WORKDIR /home
	
	ENTRYPOINT ["/bin/sh"]
	`

	err := i.Build(df)
	assert.NoError(t, err)
}
