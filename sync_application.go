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

	// 接入凭证管理
	CreateAccessCode(accessType string) (*CreateAccessCodeResponse, error)

	// 数据流转规则管理

	// 设备影子
	ShowDeviceShadow(deviceId string) (*ShowDeviceShadowResponse, error)
	UpdateDeviceShadow(deviceId string, request UpdateDeviceShadowRequest) (*ShowDeviceShadowResponse, error)

	// 设备组管理
	ListDeviceGroups(request ListDeviceGroupRequest) (*ListDeviceGroupResponse, error)
	CreateDeviceGroup(request CreateDeviceGroupRequest) (*CreateDeviceGroupResponse, error)
	ShowDeviceGroup(deviceGroupId string) (*ShowDeviceGroupResponse, error)
	UpdateDeviceGroup(deviceGroupId string, request UpdateDeviceGroupRequest) (*UpdateDeviceGroupResponse, error)
	DeleteDeviceGroup(deviceGroupId string) (bool, error)

	AddDeviceToDeviceGroup(deviceGroupId, deviceId string) (bool, error)
	RemoveDeviceFromDeviceGroup(deviceGroupId, deviceId string) (bool, error)
	ListDeviceInDeviceGroup(deviceGroupId string, request ListDeviceInDeviceGroupRequest) (*ListDeviceInDeviceGroupRequest, error)
	// 标签管理
	// 批量任务
	// 设备CA证书管理
}

type iotSyncApplicationClient struct {
	client  *resty.Client
	options ApplicationOptions
}

func (a *iotSyncApplicationClient) ListDeviceInDeviceGroup(deviceGroupId string, request ListDeviceInDeviceGroupRequest) (*ListDeviceInDeviceGroupRequest, error) {
	rawRequest := a.client.R().
		SetHeader("Content-Type", "application/json")
	if request.Limit >= 1 && request.Limit <= 50 {
		rawRequest.SetQueryParam("limit", strconv.Itoa(request.Limit))
	} else {
		rawRequest.SetQueryParam("limit", strconv.Itoa(10))
	}

	if len(request.Marker) != 0 {
		rawRequest.SetQueryParam("marker", request.Marker)
	}

	if request.Offset >= 0 && request.Offset <= 500 {
		rawRequest.SetQueryParam("offset", strconv.Itoa(request.Offset))
	} else {
		rawRequest.SetQueryParam("offset", strconv.Itoa(0))
	}

	httpResponse, err := rawRequest.
		SetPathParam("group_id", deviceGroupId).
		Get("/v5/iot/{project_id}/device-group/{group_id}/devices")
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode() != 200 {
		return nil, convertResponseToApplicationError(httpResponse)
	}

	response := &ListDeviceInDeviceGroupRequest{}

	err = json.Unmarshal(httpResponse.Body(), response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *iotSyncApplicationClient) AddDeviceToDeviceGroup(deviceGroupId, deviceId string) (bool, error) {
	return a.manageDeviceGroupDevices(deviceGroupId, "addDevice", deviceId)

}

func (a *iotSyncApplicationClient) RemoveDeviceFromDeviceGroup(deviceGroupId, deviceId string) (bool, error) {
	return a.manageDeviceGroupDevices(deviceGroupId, "removeDevice", deviceId)
}

func (a *iotSyncApplicationClient) manageDeviceGroupDevices(deviceGroupId, actionId, deviceId string) (bool, error) {
	httpResponse, err := a.client.R().
		SetPathParam("group_id", deviceGroupId).
		SetQueryParam("action_id", actionId).
		SetQueryParam("device_id", deviceId).
		Post("/v5/iot/{project_id}/device-group/{group_id}/action")
	if err != nil {
		return false, err
	}

	if httpResponse.StatusCode() != 200 {
		return false, convertResponseToApplicationError(httpResponse)
	}

	return true, nil
}
func (a *iotSyncApplicationClient) ListDeviceGroups(request ListDeviceGroupRequest) (*ListDeviceGroupResponse, error) {
	rawRequest := a.client.R().
		SetHeader("Content-Type", "application/json")
	if request.Limit >= 1 && request.Limit <= 50 {
		rawRequest.SetQueryParam("limit", strconv.Itoa(request.Limit))
	} else {
		rawRequest.SetQueryParam("limit", strconv.Itoa(10))
	}

	if len(request.Marker) != 0 {
		rawRequest.SetQueryParam("marker", request.Marker)
	}

	if request.Offset >= 0 && request.Offset <= 500 {
		rawRequest.SetQueryParam("offset", strconv.Itoa(request.Offset))
	} else {
		rawRequest.SetQueryParam("offset", strconv.Itoa(0))
	}

	if len(request.LastModifiedTime) != 0 {
		rawRequest.SetQueryParam("last_modified_time", request.LastModifiedTime)
	}

	if len(request.AppId) != 0 {
		rawRequest.SetQueryParam("app_id", request.AppId)
	}

	httpResponse, err := rawRequest.
		Get("/v5/iot/{project_id}/device-group")
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode() != 200 {
		return nil, convertResponseToApplicationError(httpResponse)
	}

	response := &ListDeviceGroupResponse{}

	err = json.Unmarshal(httpResponse.Body(), response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *iotSyncApplicationClient) DeleteDeviceGroup(deviceGroupId string) (bool, error) {
	httpResponse, err := a.client.R().
		SetPathParam("group_id", deviceGroupId).
		Delete("/v5/iot/{project_id}/device-group/{group_id}")
	if err != nil {
		return false, err
	}

	if httpResponse.StatusCode() != 200 {
		return false, convertResponseToApplicationError(httpResponse)
	}

	return true, nil
}

func (a *iotSyncApplicationClient) UpdateDeviceGroup(deviceGroupId string, request UpdateDeviceGroupRequest) (*UpdateDeviceGroupResponse, error) {
	binaryRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpResponse, err := a.client.R().
		SetPathParam("group_id", deviceGroupId).
		SetBody(binaryRequest).
		Get("/v5/iot/{project_id}/device-group/{group_id}")
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode() != 200 {
		return nil, convertResponseToApplicationError(httpResponse)
	}

	response := &UpdateDeviceGroupResponse{}

	err = json.Unmarshal(httpResponse.Body(), response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *iotSyncApplicationClient) ShowDeviceGroup(deviceGroupId string) (*ShowDeviceGroupResponse, error) {
	httpResponse, err := a.client.R().
		SetPathParam("group_id", deviceGroupId).
		Get("/v5/iot/{project_id}/device-group/{group_id}")
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode() != 200 {
		return nil, convertResponseToApplicationError(httpResponse)
	}

	response := &ShowDeviceGroupResponse{}

	err = json.Unmarshal(httpResponse.Body(), response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *iotSyncApplicationClient) CreateDeviceGroup(request CreateDeviceGroupRequest) (*CreateDeviceGroupResponse, error) {
	binaryRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpResponse, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(binaryRequest).
		Post("/v5/iot/{project_id}/device-group")
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode() != 201 {
		return nil, convertResponseToApplicationError(httpResponse)
	}

	response := &CreateDeviceGroupResponse{}

	err = json.Unmarshal(httpResponse.Body(), response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *iotSyncApplicationClient) UpdateDeviceShadow(deviceId string, request UpdateDeviceShadowRequest) (*ShowDeviceShadowResponse, error) {
	binaryRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpResponse, err := a.client.R().
		SetPathParam("device_id", deviceId).
		SetBody(binaryRequest).
		Put("/v5/iot/{project_id}/devices/{device_id}/shadow")
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode() != 200 {
		return nil, convertResponseToApplicationError(httpResponse)
	}

	response := &ShowDeviceShadowResponse{}

	err = json.Unmarshal(httpResponse.Body(), response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
func (a *iotSyncApplicationClient) ShowDeviceShadow(deviceId string) (*ShowDeviceShadowResponse, error) {
	response, err := a.client.R().
		SetPathParam("device_id", deviceId).
		Get("/v5/iot/{project_id}/devices/{device_id}/shadow")
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != 200 {
		return nil, convertResponseToApplicationError(response)
	}

	result := &ShowDeviceShadowResponse{}

	err = json.Unmarshal(response.Body(), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *iotSyncApplicationClient) CreateAccessCode(accessType string) (*CreateAccessCodeResponse, error) {
	glog.Infof("begin to create access code for type %s", accessType)
	req := struct {
		Type string `json:"type"`
	}{
		Type: "AMQP",
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBytes).
		Post("/v5/iot/{project_id}/auth/accesscode")
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != 201 {
		return nil, convertResponseToApplicationError(response)
	}

	resp := &CreateAccessCodeResponse{}

	err = json.Unmarshal(response.Body(), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *iotSyncApplicationClient) DeleteAmqpQueue(queueId string) bool {
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

func (a *iotSyncApplicationClient) ShowAmqpQueue(queueId string) (*ShowAmqpQueueResponse, error) {
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

func (a *iotSyncApplicationClient) CreateAmqpQueue(queueName string) *CreateAmqpQueueResponse {
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

func (a *iotSyncApplicationClient) ListAmqpQueues(req ListAmqpQueuesRequest) *ListAmqpQueuesResponse {
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

func (a *iotSyncApplicationClient) ResetDeviceSecret(deviceId, secret string, forceDisconnect bool) *ResetDeviceSecretResponse {
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

func (a *iotSyncApplicationClient) FreezeDevice(deviceId string) bool {
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

func (a *iotSyncApplicationClient) UnFreezeDevice(deviceId string) bool {
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

func (a *iotSyncApplicationClient) DeleteDevice(deviceId string) bool {
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

func (a *iotSyncApplicationClient) UpdateDevice(deviceId string, request UpdateDeviceRequest) *DeviceDetailResponse {
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

func (a *iotSyncApplicationClient) ShowDevice(deviceId string) *DeviceDetailResponse {
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

func (a *iotSyncApplicationClient) CreateDevice(request CreateDeviceRequest) *CreateDeviceResponse {
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

func (a *iotSyncApplicationClient) ListDevices(queryParas map[string]string) *ListDeviceResponse {
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

func (a *iotSyncApplicationClient) UpdateDeviceProperties(deviceId string, services interface{}) bool {
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

func (a *iotSyncApplicationClient) QueryDeviceProperties(deviceId, serviceId string) string {
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

func (a *iotSyncApplicationClient) SendDeviceSyncCommand(deviceId string, request DeviceSyncCommandRequest) *DeviceSyncCommandResponse {
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

func (a *iotSyncApplicationClient) SendDeviceMessage(deviceId string, msg SendDeviceMessageRequest) *SendDeviceMessageResponse {
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

func (a *iotSyncApplicationClient) ListDeviceMessages(deviceId string) *DeviceMessages {
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

func (a *iotSyncApplicationClient) ShowDeviceMessage(deviceId, messageId string) *DeviceMessage {
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

func (a *iotSyncApplicationClient) ListApplications() *Applications {
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

func (a *iotSyncApplicationClient) ShowApplication(appId string) *Application {
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

func (a *iotSyncApplicationClient) DeleteApplication(appId string) bool {
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

func (a *iotSyncApplicationClient) CreateApplication(request ApplicationCreateRequest) *Application {
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

func CreateSyncIotApplicationClient(options ApplicationOptions) *iotSyncApplicationClient {
	c := &iotSyncApplicationClient{

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
