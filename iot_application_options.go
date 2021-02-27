package iot

import (
	"fmt"
)

type ApplicationOptions struct {
	Server     string
	Ak         string
	Sk         string
	User       string
	Password   string
	InstanceId string
	UseAkSk    bool
	ProjectId  string
}

func NewApplicationOptions() *ApplicationOptions {
	o := &ApplicationOptions{
		Server:     "",
		Ak:         "",
		Sk:         "",
		User:       "",
		Password:   "",
		InstanceId: "",
		UseAkSk:    false,
	}
	return o
}

func (o *ApplicationOptions) AddServer(server string) *ApplicationOptions {
	if len(server) == 0 {
		fmt.Println("server is empty")
		o.Server = "https://iotda.cn-north-4.myhuaweicloud.com:443"
	} else {
		o.Server = server
	}

	return o
}

func (o *ApplicationOptions) AddAk(ak string) *ApplicationOptions {
	if len(ak) == 0 {
		fmt.Println("ak is empty")
	} else {
		o.Ak = ak
	}

	return o
}

func (o *ApplicationOptions) AddSk(sk string) *ApplicationOptions {
	if len(sk) == 0 {
		fmt.Println("ak is empty")
	} else {
		o.Sk = sk
	}

	return o
}

func (o *ApplicationOptions) AddUser(user string) *ApplicationOptions {
	if len(user) == 0 {
		fmt.Println("ak is empty")
	} else {
		o.User = user
	}

	return o
}

func (o *ApplicationOptions) AddPassword(password string) *ApplicationOptions {
	if len(password) == 0 {
		fmt.Println("ak is empty")
	} else {
		o.Password = password
	}

	return o
}

func (o *ApplicationOptions) AddInstanceId(instanceId string) *ApplicationOptions {
	if len(instanceId) == 0 {
		fmt.Println("ak is empty")
	} else {
		o.InstanceId = instanceId
	}

	return o
}

func (o *ApplicationOptions) IsUseAkSk(useAkSk bool) *ApplicationOptions {
	o.UseAkSk = useAkSk

	return o
}

func (o *ApplicationOptions) SetProjectId(projectId string) *ApplicationOptions {
	o.ProjectId = projectId

	return o
}
