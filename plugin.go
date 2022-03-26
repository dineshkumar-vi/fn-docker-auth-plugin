package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/engine-api/types"
	"github.com/docker/go-plugins-helpers/authorization"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
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

	uri, err := url.QueryUnescape(req.RequestURI)

	if err != nil {
		return authorization.Response{Err: err.Error()}
	}

	// Remove query parameters
	i := strings.Index(uri, "?")
	if i > 0 {
		uri = uri[:i]
	}

	fmt.Println("checking " + req.RequestMethod + " request to '" + uri + "' from user : " + req.User)
	if req.RequestMethod == "DELETE" && deletePlugin.MatchString(uri) {
		return authorization.Response{Err: "Permission Denied !"}
	}

	if req.RequestMethod == "POST" && disablePlugin.MatchString(uri) {
		return authorization.Response{Err: "Permission Denied !"}
	}

	if req.RequestMethod == "GET" && strings.Contains(uri, "images") {

	}

	return authorization.Response{Allow: true}
}

func (p *plugin) AuthZRes(req authorization.Request) authorization.Response {

	return authorization.Response{Allow: true}
}

func inspectAndDrop(imageName string) (bool, error) {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
	}

	var response *http.Response
	var err error

	response, err = httpc.Get("http://unix" + "/images/" + imageName + "/json")

	body, err := ioutil.ReadAll(response.Body)

	var data types.ImageInspect
	err = json.Unmarshal(body, &data)

	c1 := isRepoValid(data.RepoTags)
	c2 := strings.Contains(data.Config.User, "nsnbdy")
	c3 := strings.Contains(data.Config.Labels["maintainer"], "fanniemae.com")
	return c1 || c2 || c3, err
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
