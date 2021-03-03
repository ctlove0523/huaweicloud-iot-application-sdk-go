package iot

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/golang/glog"
	"strconv"
	"time"
)

type ApplicationClient interface {
	// 资源空间管理
	ListApplications() *Applications
	ShowApplication(appId string) *Application
	DeleteApplication(appId string) bool
	CreateApplication(request ApplicationCreateRequest) *Application

	// 设备管理
	ListDevices(queryParas map[string]string) *ListDeviceResponse
	CreateDevice(request CreateDeviceRequest) *CreateDeviceResponse
	ShowDevice(deviceId string) *DeviceDetailResponse
	UpdateDevice(deviceId string, request UpdateDeviceRequest) *DeviceDetailResponse
	DeleteDevice(deviceId string) bool
	FreezeDevice(deviceId string) bool
	UnFreezeDevice(deviceId string) bool
	ResetDeviceSecret(deviceId, secret string, forceDisconnect bool) *ResetDeviceSecretResponse

	// 设备消息
	ListDeviceMessages(deviceId string) *DeviceMessages
	ShowDeviceMessage(deviceId, messageId string) *DeviceMessage
	SendDeviceMessage(deviceId string, msg SendDeviceMessageRequest) *SendDeviceMessageResponse

	// 设备命令
	SendDeviceSyncCommand(deviceId string, request DeviceSyncCommandRequest) *DeviceSyncCommandResponse

	// 设备属性
	QueryDeviceProperties(deviceId, serviceId string) string
	UpdateDeviceProperties(deviceId string, services interface{}) bool

	// AMQP队列管理
	ListAmqpQueues(req ListAmqpQueuesRequest) *ListAmqpQueuesResponse
	CreateAmqpQueue(queueName string) *CreateAmqpQueueResponse
	ShowAmqpQueue(queueId string) (*ShowAmqpQueueResponse, error)
	DeleteAmqpQueue(queueId string) bool
}

type iotApplicationClient struct {
	client  *resty.Client
	options ApplicationOptions
}

func (a *iotApplicationClient) DeleteAmqpQueue(queueId string) bool {
	glog.Infof("begin to delete amqp queue with id %s", queueId)
	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParam("queue_id", queueId).
		Delete("v5/iot/{project_id}/amqp-queues/{queue_id}")
	if err != nil {
		return false
	}

	if response.StatusCode() != 204 {
		glog.Warningf("delete amqp queue response code is %d", response.StatusCode())
		return false
	}

	return true

}

func (a *iotApplicationClient) ShowAmqpQueue(queueId string) (*ShowAmqpQueueResponse, error) {
	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParam("queue_id", queueId).
		Get("v5/iot/{project_id}/amqp-queues/{queue_id}")
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != 200 {
		err = convertResponseToApplicationError(response)
		return nil, err
	}

	resp := &ShowAmqpQueueResponse{}

	err = json.Unmarshal(response.Body(), resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *iotApplicationClient) CreateAmqpQueue(queueName string) *CreateAmqpQueueResponse {
	createAmqpRequest := struct {
		QueueName string `json:"queue_name,omitempty"`
	}{QueueName: queueName}

	requestBytes, err := json.Marshal(createAmqpRequest)
	if err != nil {
		return nil
	}
	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBytes).
		Post("/v5/iot/{project_id}/amqp-queues")
	if err != nil {
		return nil
	}

	if response.StatusCode() != 201 {
		fmt.Println(response.StatusCode())
		fmt.Println(string(response.Body()))
		return nil
	}

	resp := &CreateAmqpQueueResponse{}

	err = json.Unmarshal(response.Body(), resp)
	if err != nil {
		return nil
	}

	return resp
}

func (a *iotApplicationClient) ListAmqpQueues(req ListAmqpQueuesRequest) *ListAmqpQueuesResponse {
	queryParas := map[string]string{}
	if len(req.QueueName) != 0 {
		queryParas["queue_name"] = req.QueueName
	}
	if req.Limit == 0 {
		queryParas["limit"] = "10"
	} else {
		queryParas["limit"] = strconv.Itoa(req.Limit)
	}
	if len(req.Marker) != 0 {
		queryParas["marker"] = req.Marker
	}
	if len(req.Offset) != 0 {
		queryParas["offset"] = req.Offset
	}

	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(queryParas).
		Get("/v5/iot/{project_id}/amqp-queues")

	if err != nil {
		return nil
	}

	if response.StatusCode() != 200 {
		fmt.Println(response.StatusCode())
		fmt.Println(string(response.Body()))
		return nil
	}

	resp := &ListAmqpQueuesResponse{}

	err = json.Unmarshal(response.Body(), resp)
	if err != nil {
		return nil
	}

	return resp
}

func (a *iotApplicationClient) ResetDeviceSecret(deviceId, secret string, forceDisconnect bool) *ResetDeviceSecretResponse {
	resetSecret := struct {
		Secret          string `json:"secret,omitempty"`
		ForceDisconnect bool   `json:"force_disconnect,omitempty"`
	}{Secret: secret, ForceDisconnect: forceDisconnect}

	body, err := json.Marshal(resetSecret)
	if err != nil {
		fmt.Println("marshal failed")
		return nil
	}
	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		SetQueryParams(map[string]string{
			"action_id": "resetSecret",
		}).
		Post("/v5/iot/{project_id}/devices/{device_id}/action")
	if err != nil {
		fmt.Println("reset device secret failed failed")
		return nil
	}

	resp := &ResetDeviceSecretResponse{}

	fmt.Println(string(response.Body()))
	err = json.Unmarshal(response.Body(), resp)
	if err != nil {
		return nil
	}
	return resp
}

func (a *iotApplicationClient) FreezeDevice(deviceId string) bool {
	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Post("/v5/iot/{project_id}/devices/{device_id}/freeze")
	if err != nil {
		fmt.Println("freeze device failed")
		return false
	}

	return successResponse(response)
}

func (a *iotApplicationClient) UnFreezeDevice(deviceId string) bool {
	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Post("/v5/iot/{project_id}/devices/{device_id}/unfreeze")
	if err != nil {
		fmt.Println("unfreeze device failed")
		return false
	}

	return successResponse(response)
}

func (a *iotApplicationClient) DeleteDevice(deviceId string) bool {
	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Delete("/v5/iot/{project_id}/devices/{device_id}")
	if err != nil {
		fmt.Println("list devices failed")
		return false
	}

	return successResponse(response)
}

func (a *iotApplicationClient) UpdateDevice(deviceId string, request UpdateDeviceRequest) *DeviceDetailResponse {
	body, err := json.Marshal(request)
	if err != nil {
		fmt.Println("marshal failed")
		return nil
	}

	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Put("/v5/iot/{project_id}/devices/{device_id}")
	if err != nil {
		fmt.Println("list devices failed")
		return nil
	}

	device := &DeviceDetailResponse{}
	err = json.Unmarshal(response.Body(), device)
	if err != nil {
		fmt.Println("unmarshal response failed")
		return nil
	}

	return device

}

func (a *iotApplicationClient) ShowDevice(deviceId string) *DeviceDetailResponse {
	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Get("/v5/iot/{project_id}/devices/{device_id}")
	if err != nil {
		fmt.Println("list devices failed")
		return nil
	}

	deviceDetail := &DeviceDetailResponse{}
	err = json.Unmarshal(response.Body(), deviceDetail)
	if err != nil {
		fmt.Println("unmarshal failed")
		return nil
	}

	return deviceDetail
}

func (a *iotApplicationClient) CreateDevice(request CreateDeviceRequest) *CreateDeviceResponse {
	bytesBody, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
	}
	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bytesBody).
		Post("/v5/iot/{project_id}/devices")

	if err != nil {
		fmt.Println("create device failed")
		fmt.Println(err)
	}

	resp := &CreateDeviceResponse{}
	err = json.Unmarshal(response.Body(), resp)
	fmt.Println(string(response.Body()))
	return resp
}

func (a *iotApplicationClient) ListDevices(queryParas map[string]string) *ListDeviceResponse {
	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(queryParas).
		Get("/v5/iot/{project_id}/devices")
	if err != nil {
		fmt.Println("list devices failed")
		return nil
	}

	if !successResponse(response) {
		fmt.Println("response failed")
		return nil
	}

	devices := &ListDeviceResponse{}

	err = json.Unmarshal(response.Body(), devices)
	if err != nil {
		fmt.Println("un marshal failed")
		return nil
	}

	return devices
}

func (a *iotApplicationClient) UpdateDeviceProperties(deviceId string, services interface{}) bool {
	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		SetBody(services).
		Put("/v5/iot/{project_id}/devices/{device_id}/properties")

	if err != nil {
		fmt.Printf("query device properties failed %s", err)
		return false
	}

	return response.StatusCode() == 200
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

	go logFlush()

	return c
}

func logFlush() {
	ticker := time.Tick(5 * time.Second)
	for {
		select {
		case <-ticker:
			glog.Flush()
		}
	}
}

func successResponse(response *resty.Response) bool {
	if response.StatusCode() >= 200 && response.StatusCode() < 300 {
		return true
	}

	return false
}
