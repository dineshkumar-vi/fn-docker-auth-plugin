package main

import (
	"github.com/docker/go-plugins-helpers/authorization"
)

func newPlugin() (*plugin, error) {
	return &plugin{}, nil
}

type plugin struct {
}

func (p *plugin) AuthZReq(req authorization.Request) authorization.Response {
	return authorization.Response{Allow: true}
}

func (p *plugin) AuthZRes(req authorization.Request) authorization.Response {

	return authorization.Response{Allow: true}
}

//Diagnostics Bundle: /var/folders/yv/2n3czkrs0cxc__8hmhntm4m40000gn/T/B338FD20-D8D1-41BA-BCD9-16EC602B9C51/20220324050038.zip
//Diagnostics ID:     B338FD20-D8D1-41BA-BCD9-16EC602B9C51/20220324050038 (uploaded)
