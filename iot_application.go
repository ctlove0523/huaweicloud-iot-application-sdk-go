package iot

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

type ApplicationClient interface {
	ListApplications() *Applications
	ShowApplication(appId string) *Application
	DeleteApplication(appId string) bool
	CreateApplication(request ApplicationCreateRequest) *Application

	// 设备消息
	ListDeviceMessages(deviceId string) *DeviceMessages
	ShowDeviceMessage(deviceId, messageId string) *DeviceMessage
	SendDeviceMessage(deviceId string, msg SendDeviceMessageRequest) *SendDeviceMessageResponse

	// 设备命令
	SendDeviceSyncCommand(deviceId string, request DeviceSyncCommandRequest) *DeviceSyncCommandResponse

	// 设备属性
	QueryDeviceProperties(deviceId, serviceId string) string
}

type iotApplicationClient struct {
	client  *resty.Client
	options ApplicationOptions
}

func (a *iotApplicationClient) QueryDeviceProperties(deviceId, serviceId string) string {
	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("service_id", serviceId).
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Get("/v5/iot/{project_id}/devices/{device_id}/properties")

	if err != nil {
		fmt.Printf("query device properties failed %s", err)
		return ""
	}

	return string(response.Body())
}

func (a *iotApplicationClient) SendDeviceSyncCommand(deviceId string, request DeviceSyncCommandRequest) *DeviceSyncCommandResponse {
	reqBody, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("marshal device sync command request failed %s", err)
		return &DeviceSyncCommandResponse{}
	}
	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Post("/v5/iot/{project_id}/devices/{device_id}/commands")
	if err != nil {
		fmt.Printf("send device command failed %s", err)
		return &DeviceSyncCommandResponse{}
	}

	fmt.Printf(response.Status())

	resp := &DeviceSyncCommandResponse{}

	err = json.Unmarshal(response.Body(), resp)
	if err != nil {
		fmt.Println(err)
		return &DeviceSyncCommandResponse{}
	}

	return resp
}

func (a *iotApplicationClient) SendDeviceMessage(deviceId string, msg SendDeviceMessageRequest) *SendDeviceMessageResponse {
	reqBody, err := json.Marshal(msg)
	if err != nil {
		return &SendDeviceMessageResponse{}
	}
	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Post("/v5/iot/{project_id}/devices/{device_id}/messages")
	if err != nil {
		return &SendDeviceMessageResponse{}
	}

	resp := &SendDeviceMessageResponse{}

	err = json.Unmarshal(response.Body(), resp)
	if err != nil {
		return &SendDeviceMessageResponse{}
	}

	return resp
}

func (a *iotApplicationClient) ListDeviceMessages(deviceId string) *DeviceMessages {
	response, err := a.client.R().
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Get("/v5/iot/{project_id}/devices/{device_id}/messages")
	if err != nil {
		fmt.Println("list device messages error")
		return &DeviceMessages{}
	}

	messages := &DeviceMessages{}
	err = json.Unmarshal(response.Body(), messages)
	if err != nil {
		fmt.Println("deserialize device message failed")
		fmt.Println(err)
	}

	return messages
}

func (a *iotApplicationClient) ShowDeviceMessage(deviceId, messageId string) *DeviceMessage {
	response, err := a.client.R().
		SetPathParams(map[string]string{
			"device_id":  deviceId,
			"message_id": messageId,
		}).
		Get("/v5/iot/{project_id}/devices/{device_id}/messages/{message_id}")
	if err != nil {
		fmt.Println("list device messages error")
		return &DeviceMessage{}
	}

	messages := &DeviceMessage{}
	err = json.Unmarshal(response.Body(), messages)
	if err != nil {
		fmt.Println("deserialize device message failed")
		fmt.Println(err)
	}

	return messages
}

func (a *iotApplicationClient) ListApplications() *Applications {
	response, err := a.client.R().Get("/v5/iot/{project_id}/apps")
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

func (a *iotApplicationClient) ShowApplication(appId string) *Application {
	response, err := a.client.R().
		SetPathParams(map[string]string{
			"app_id": appId,
		}).
		Get("/v5/iot/{project_id}/apps/{app_id}")
	if err != nil {
		fmt.Println("get apps failed")
		return &Application{}
	}

	app := &Application{}
	err = json.Unmarshal(response.Body(), app)
	if err != nil {
		fmt.Println("deserialize applications failed")
	}

	return app
}

func (a *iotApplicationClient) DeleteApplication(appId string) bool {
	response, err := a.client.R().
		SetPathParams(map[string]string{
			"app_id": appId,
		}).
		Delete("/v5/iot/{project_id}/apps/{app_id}")
	if err != nil {
		fmt.Printf("delete apps %s failed", appId)
		return false
	}

	if response.StatusCode() != 204 {
		fmt.Printf("delete app %s failed,response code is %d", appId, response.StatusCode())
		return false
	}

	return true
}

func (a *iotApplicationClient) CreateApplication(request ApplicationCreateRequest) *Application {
	body, err := json.Marshal(request)
	if err != nil {
		fmt.Println("marshal application create request failed")
		return &Application{}
	}

	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("/v5/iot/{project_id}/apps")
	if err != nil {
		fmt.Println("create app failed")
		return &Application{}
	}

	fmt.Println(response.Status())
	fmt.Println(string(response.Body()))

	app := &Application{}
	err = json.Unmarshal(response.Body(), app)
	if err != nil {
		fmt.Println("deserialize applications failed")
	}

	return app
}

func CreateIotApplicationClient(options ApplicationOptions) *iotApplicationClient {
	c := &iotApplicationClient{

	}
	c.options = options
	c.client = resty.New()
	if len(options.ServerAddress) > 0 {
		c.client.SetHostURL("https://" + options.ServerAddress + ":" + strconv.Itoa(options.ServerPort))
	} else {
		c.client.SetHostURL("https://iotda.cn-north-4.myhuaweicloud.com")
	}

	c.client.SetPathParams(map[string]string{
		"project_id": options.ProjectId,
	})

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
			request.SetHeader("Authorization", " "+signedMsg)
		} else {
			request.SetHeader("X-Auth-Token", options.Credential.Token)
		}

		if len(options.InstanceId) != 0 {
			request.SetHeader("Instance-Id", options.InstanceId)
		}

		return nil
	})

	return c
}
