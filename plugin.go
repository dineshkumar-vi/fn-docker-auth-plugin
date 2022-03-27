package main

import (
	"context"
	"github.com/docker/engine-api/client"
	"github.com/docker/go-plugins-helpers/authorization"
	"os"
	"regexp"
	"strings"
)

func newPlugin() (*plugin, error) {
	return &plugin{}, nil
}

type plugin struct {
}

var (
	create       = regexp.MustCompile(`/containers/create`)
	inageInspect = regexp.MustCompile(`/images/.+/json`)

	deletePlugin  = regexp.MustCompile(`/v.*/plugins/dineshviveck5/docker-authz-plugin:dev`)
	disablePlugin = regexp.MustCompile(`/v.*/plugins/dineshviveck5/docker-authz-plugin:dev/disable`)
)

func (p *plugin) AuthZReq(req authorization.Request) authorization.Response {

	if strings.Contains(req.RequestURI, "/images/load") {

		tag, e := inspectAndDrop("demo")

		if e != nil {
			return authorization.Response{Err: e.Error()}
		} else {
			return authorization.Response{Err: tag}
		}
	}
	return authorization.Response{Allow: true}
}

func (p *plugin) AuthZRes(req authorization.Request) authorization.Response {

	return authorization.Response{Allow: true}
}

func inspectAndDrop(imageName string) (string, error) {

	var version string
	if version = os.Getenv("DOCKER_API_VERSION"); version == "" {
		version = "v1.41"
	}

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli:1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", version, nil, defaultHeaders)
	if err != nil {
		return "", err
	}

	imagedetails, _, err := cli.ImageInspectWithRaw(context.Background(), "demo", false)
	return imagedetails.ID, err

}

func isRepoValid(repo []string) bool {
	isRepoValid := false
	for i := 0; i < len(repo); i++ {
		if strings.Contains(repo[i], "fanniemae.com") {
			isRepoValid = true
		}
	}
	return isRepoValid
}
