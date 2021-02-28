package iot

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

type ApplicationClient interface {
	ListApplications(instanceId, projectId string, defaultApp bool) *Applications
}

type iotApplicationClient struct {
	client  *resty.Client
	options ApplicationOptions
}

func CreateIotApplicationClient(options ApplicationOptions) *iotApplicationClient {
	c := &iotApplicationClient{

	}
	c.options = options
	c.client = resty.New().SetHostURL("https://iotda.cn-north-4.myhuaweicloud.com")
	c.client.SetRetryCount(3)
	c.client.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		if len(request.Header.Get("Content-Type")) == 0 {
			fmt.Println("content type not exist,begin to set")
			request.SetHeader("Content-Type", "application/json")
		}

		xSdkDate := time.Now().UTC().Format("20060102T150405Z")
		request.SetHeader("X-Sdk-Date", xSdkDate)

		if options.Credential.UseAkSk {
			signedMsg := SignMessage(request, options.Credential.Sk, options.Credential.Ak)
			fmt.Println(signedMsg)
			request.SetHeader("Authorization", " "+signedMsg)
		}

		return nil
	})

	return c
}

func (a *iotApplicationClient) ListApplications(instanceId, projectId string) *Applications {
	req := a.client.R().
		SetPathParams(map[string]string{
			"project_id": projectId,
		})
	if len(instanceId) > 0 {
		fmt.Println("begin to set instance id")
		req.SetHeader("Instance-Id", instanceId)
	}

	response, err := req.Get("/v5/iot/{project_id}/apps")
	if err != nil {
		fmt.Println("get apps failed")
		return &Applications{}
	}

	app := &Applications{}
	err = json.Unmarshal(response.Body(), app)
	if err != nil {
		fmt.Println("deserialize applications failed")
	}

	return app
}
