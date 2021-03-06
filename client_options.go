package iot

import (
	"fmt"
)

type Credentials struct {
	Ak      string
	Sk      string
	Token   string
	UseAkSk bool
}

type ApplicationOptions struct {
	ServerAddress string
	ServerPort    int
	InstanceId    string
	ProjectId     string
	Credential    *Credentials
}

func NewApplicationOptions() *ApplicationOptions {
	o := &ApplicationOptions{
		ServerAddress: "",
		ServerPort:    443,
		InstanceId:    "",
		ProjectId:     "",
		Credential:    nil,
	}
	return o
}

func (o *ApplicationOptions) AddServer(server string) *ApplicationOptions {
	if len(server) == 0 {
		fmt.Println("server is empty")
		o.ServerAddress = "https://iotda.cn-north-4.myhuaweicloud.com:443"
	} else {
		o.ServerAddress = server
	}

	return o
}

func (o *ApplicationOptions) AddServerPort(port int) *ApplicationOptions {
	o.ServerPort = port
	return o
}

func (o *ApplicationOptions) AddAk(ak string) *ApplicationOptions {
	if len(ak) == 0 {
		fmt.Println("ak is empty")
	} else {
		o.Credential.Ak = ak
	}

	return o
}

func (o *ApplicationOptions) AddSk(sk string) *ApplicationOptions {
	if len(sk) == 0 {
		fmt.Println("ak is empty")
	} else {
		o.Credential.Sk = sk
	}

	return o
}

func (o *ApplicationOptions) SetToken(token string) *ApplicationOptions {
	if len(token) == 0 {
		fmt.Println("Token is empty")
	} else {
		o.Credential.Token = token
	}

	return o
}

func (o *ApplicationOptions) AddInstanceId(instanceId string) *ApplicationOptions {
	if len(instanceId) == 0 {
		fmt.Println("Instance is is empty")
	} else {
		o.InstanceId = instanceId
	}

	return o
}

func (o *ApplicationOptions) IsUseAkSk(useAkSk bool) *ApplicationOptions {
	o.Credential.UseAkSk = useAkSk
	return o
}

func (o *ApplicationOptions) SetProjectId(projectId string) *ApplicationOptions {
	o.ProjectId = projectId
	return o
}
