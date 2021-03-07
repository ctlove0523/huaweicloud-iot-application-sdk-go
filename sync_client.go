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
	// 设备管理
	ListDevices(queryParas map[string]string) (*ListDeviceResponse, error)
	CreateDevice(request CreateDeviceRequest) (*CreateDeviceResponse, error)
	ShowDevice(deviceId string) (*DeviceDetailResponse, error)
	UpdateDevice(deviceId string, request UpdateDeviceRequest) (*DeviceDetailResponse, error)
	DeleteDevice(deviceId string) (bool, error)
	FreezeDevice(deviceId string) (bool, error)
	UnFreezeDevice(deviceId string) (bool, error)
	ResetDeviceSecret(deviceId, secret string, forceDisconnect bool) (*ResetDeviceSecretResponse, error)

	// 设备消息
	ListDeviceMessages(deviceId string) (*DeviceMessages, error)
	ShowDeviceMessage(deviceId, messageId string) (*DeviceMessage, error)
	SendDeviceMessage(deviceId string, msg SendDeviceMessageRequest) (*SendDeviceMessageResponse, error)

	// 设备命令
	SendDeviceSyncCommand(deviceId string, request DeviceSyncCommandRequest) (*DeviceSyncCommandResponse, error)

	// 设备属性
	QueryDeviceProperties(deviceId, serviceId string) (interface{}, error)
	UpdateDeviceProperties(deviceId string, services interface{}) (bool, error)

	// AMQP队列管理
	ListAmqpQueues(req ListAmqpQueuesRequest) (*ListAmqpQueuesResponse, error)
	CreateAmqpQueue(queueName string) (*CreateAmqpQueueResponse, error)
	ShowAmqpQueue(queueId string) (*ShowAmqpQueueResponse, error)
	DeleteAmqpQueue(queueId string) (bool, error)

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
	DeviceBindTags(request DeviceBindTagsRequest) (bool, error)
	DeviceUnBindTags(request DeviceUnBindTagsRequest) (bool, error)
	ListDeviceByTags(request ListDeviceByTagsRequest) (*ListDeviceByTagsResponse, error)

	// 资源空间管理
	ListApplications() (*Applications, error)
	ShowApplication(appId string) (*Application, error)
	DeleteApplication(appId string) (bool, error)
	CreateApplication(request ApplicationCreateRequest) (*Application, error)

	// 批量任务

	// 设备CA证书管理
	ListDeviceCertificates(request ListDeviceCertificatesRequest) (*ListDeviceCertificatesResponse, error)
	UploadDeviceCertificates(request UploadDeviceCertificatesRequest) (*UploadDeviceCertificatesResponse, error)
	DeleteDeviceCertificates(certificateId string) (bool, error)
	VerifyDeviceCertificates(certificateId, verifyContent string) (bool, error)
}

type syncClient struct {
	client  *resty.Client
	options ApplicationOptions
}

func (client *syncClient) VerifyDeviceCertificates(certificateId, verifyContent string) (bool, error) {
	requestBody := struct {
		VerifyContent string `json:"verify_content"`
	}{
		VerifyContent: verifyContent,
	}

	binaryRequest, err := json.Marshal(requestBody)
	if err != nil {
		return false, err
	}

	httpResponse, err := client.client.R().
		SetPathParam("certificate_id", certificateId).
		SetQueryParam("action_id", "verify").
		SetBody(binaryRequest).
		Post("/v5/iot/{project_id}/certificates/{certificate_id}/action")
	if err != nil {
		return false, err
	}

	if httpResponse.StatusCode() != 200 {
		return false, convertResponseToApplicationError(httpResponse)
	}

	return true, nil
}

func (client *syncClient) DeleteDeviceCertificates(certificateId string) (bool, error) {
	httpResponse, err := client.client.R().
		SetPathParam("certificate_id", certificateId).
		Delete("/v5/iot/{project_id}/certificates/{certificate_id}")
	if err != nil {
		return false, err
	}

	if httpResponse.StatusCode() != 204 {
		return false, convertResponseToApplicationError(httpResponse)
	}

	return true, nil
}

func (client *syncClient) UploadDeviceCertificates(request UploadDeviceCertificatesRequest) (*UploadDeviceCertificatesResponse, error) {
	binaryRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpResponse, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(binaryRequest).
		Post("/v5/iot/{project_id}/certificates")
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode() != 200 {
		return nil, convertResponseToApplicationError(httpResponse)
	}

	response := &UploadDeviceCertificatesResponse{}
	err = json.Unmarshal(httpResponse.Body(), response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (client *syncClient) ListDeviceCertificates(request ListDeviceCertificatesRequest) (*ListDeviceCertificatesResponse, error) {
	rawRequest := client.client.R().
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

	if len(request.AppId) != 0 {
		rawRequest.SetQueryParam("app_id", request.AppId)
	}

	httpResponse, err := rawRequest.
		Get("/v5/iot/{project_id}/certificates")
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode() != 200 {
		return nil, convertResponseToApplicationError(httpResponse)
	}

	response := &ListDeviceCertificatesResponse{}

	err = json.Unmarshal(httpResponse.Body(), response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *syncClient) ListDeviceByTags(request ListDeviceByTagsRequest) (*ListDeviceByTagsResponse, error) {
	rawRequest := client.client.R().
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

	requestBody := struct {
		ResourceType string     `json:"resource_type,omitempty"`
		Tags         []TagV5DTO `json:"tags,omitempty"`
	}{
		ResourceType: "device",
		Tags:         request.Tags,
	}

	binaryRequest, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	httpResponse, err := rawRequest.
		SetBody(binaryRequest).
		Post("/v5/iot/{project_id}/tags/query-resources")
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode() != 200 {
		return nil, convertResponseToApplicationError(httpResponse)
	}

	response := &ListDeviceByTagsResponse{}

	err = json.Unmarshal(httpResponse.Body(), response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *syncClient) DeviceUnBindTags(request DeviceUnBindTagsRequest) (bool, error) {
	binaryRequest, err := json.Marshal(request)
	if err != nil {
		return false, err
	}

	httpResponse, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(binaryRequest).
		Post("/v5/iot/{project_id}/tags/unbind-resource")
	if err != nil {
		return false, err
	}

	if httpResponse.StatusCode() != 200 {
		return false, convertResponseToApplicationError(httpResponse)
	}

	return true, nil
}

func (client *syncClient) DeviceBindTags(request DeviceBindTagsRequest) (bool, error) {
	binaryRequest, err := json.Marshal(request)
	if err != nil {
		return false, err
	}

	httpResponse, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(binaryRequest).
		Post("/v5/iot/{project_id}/tags/bind-resource")
	if err != nil {
		return false, err
	}

	if httpResponse.StatusCode() != 200 {
		return false, convertResponseToApplicationError(httpResponse)
	}

	return true, nil
}

func (client *syncClient) ListDeviceInDeviceGroup(deviceGroupId string, request ListDeviceInDeviceGroupRequest) (*ListDeviceInDeviceGroupRequest, error) {
	rawRequest := client.client.R().
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

func (client *syncClient) AddDeviceToDeviceGroup(deviceGroupId, deviceId string) (bool, error) {
	return client.manageDeviceGroupDevices(deviceGroupId, "addDevice", deviceId)

}

func (client *syncClient) RemoveDeviceFromDeviceGroup(deviceGroupId, deviceId string) (bool, error) {
	return client.manageDeviceGroupDevices(deviceGroupId, "removeDevice", deviceId)
}

func (client *syncClient) manageDeviceGroupDevices(deviceGroupId, actionId, deviceId string) (bool, error) {
	httpResponse, err := client.client.R().
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
func (client *syncClient) ListDeviceGroups(request ListDeviceGroupRequest) (*ListDeviceGroupResponse, error) {
	rawRequest := client.client.R().
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

func (client *syncClient) DeleteDeviceGroup(deviceGroupId string) (bool, error) {
	httpResponse, err := client.client.R().
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

func (client *syncClient) UpdateDeviceGroup(deviceGroupId string, request UpdateDeviceGroupRequest) (*UpdateDeviceGroupResponse, error) {
	binaryRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpResponse, err := client.client.R().
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

func (client *syncClient) ShowDeviceGroup(deviceGroupId string) (*ShowDeviceGroupResponse, error) {
	httpResponse, err := client.client.R().
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

func (client *syncClient) CreateDeviceGroup(request CreateDeviceGroupRequest) (*CreateDeviceGroupResponse, error) {
	binaryRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpResponse, err := client.client.R().
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

func (client *syncClient) UpdateDeviceShadow(deviceId string, request UpdateDeviceShadowRequest) (*ShowDeviceShadowResponse, error) {
	binaryRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpResponse, err := client.client.R().
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
func (client *syncClient) ShowDeviceShadow(deviceId string) (*ShowDeviceShadowResponse, error) {
	response, err := client.client.R().
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

func (client *syncClient) CreateAccessCode(accessType string) (*CreateAccessCodeResponse, error) {
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

	response, err := client.client.R().
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

func (client *syncClient) DeleteAmqpQueue(queueId string) (bool, error) {
	glog.Infof("begin to delete amqp queue with id %s", queueId)
	response, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParam("queue_id", queueId).
		Delete("v5/iot/{project_id}/amqp-queues/{queue_id}")
	if err != nil {
		return false, err
	}

	if response.StatusCode() != 204 {
		glog.Warningf("delete amqp queue response code is %d", response.StatusCode())
		return false, convertResponseToApplicationError(response)
	}

	return true, nil

}

func (client *syncClient) ShowAmqpQueue(queueId string) (*ShowAmqpQueueResponse, error) {
	response, err := client.client.R().
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

func (client *syncClient) CreateAmqpQueue(queueName string) (*CreateAmqpQueueResponse, error) {
	createAmqpRequest := struct {
		QueueName string `json:"queue_name,omitempty"`
	}{QueueName: queueName}

	requestBytes, err := json.Marshal(createAmqpRequest)
	if err != nil {
		return nil, nil
	}
	response, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBytes).
		Post("/v5/iot/{project_id}/amqp-queues")
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != 201 {
		return nil, convertResponseToApplicationError(response)
	}

	resp := &CreateAmqpQueueResponse{}

	err = json.Unmarshal(response.Body(), resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *syncClient) ListAmqpQueues(req ListAmqpQueuesRequest) (*ListAmqpQueuesResponse, error) {
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

	response, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(queryParas).
		Get("/v5/iot/{project_id}/amqp-queues")

	if err != nil {
		return nil, err
	}

	if response.StatusCode() != 200 {
		return nil, convertResponseToApplicationError(response)
	}

	resp := &ListAmqpQueuesResponse{}

	err = json.Unmarshal(response.Body(), resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *syncClient) ResetDeviceSecret(deviceId, secret string, forceDisconnect bool) (*ResetDeviceSecretResponse, error) {
	resetSecret := struct {
		Secret          string `json:"secret,omitempty"`
		ForceDisconnect bool   `json:"force_disconnect,omitempty"`
	}{Secret: secret, ForceDisconnect: forceDisconnect}

	body, err := json.Marshal(resetSecret)
	if err != nil {
		return nil, err
	}
	response, err := client.client.R().
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
		return nil, err
	}

	resp := &ResetDeviceSecretResponse{}

	err = json.Unmarshal(response.Body(), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (client *syncClient) FreezeDevice(deviceId string) (bool, error) {
	response, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Post("/v5/iot/{project_id}/devices/{device_id}/freeze")
	if err != nil {
		return false, err
	}

	return response.StatusCode() == 204, nil
}

func (client *syncClient) UnFreezeDevice(deviceId string) (bool, error) {
	response, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Post("/v5/iot/{project_id}/devices/{device_id}/unfreeze")
	if err != nil {
		return false, err
	}

	return response.StatusCode() == 204, nil
}

func (client *syncClient) DeleteDevice(deviceId string) (bool, error) {
	response, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Delete("/v5/iot/{project_id}/devices/{device_id}")
	if err != nil {
		return false, nil
	}

	return response.StatusCode() == 204, nil
}

func (client *syncClient) UpdateDevice(deviceId string, request UpdateDeviceRequest) (*DeviceDetailResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Put("/v5/iot/{project_id}/devices/{device_id}")
	if err != nil {
		return nil, err
	}

	device := &DeviceDetailResponse{}
	err = json.Unmarshal(response.Body(), device)
	if err != nil {
		return nil, err
	}

	return device, nil

}

func (client *syncClient) ShowDevice(deviceId string) (*DeviceDetailResponse, error) {
	response, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Get("/v5/iot/{project_id}/devices/{device_id}")
	if err != nil {
		return nil, err
	}

	deviceDetail := &DeviceDetailResponse{}
	err = json.Unmarshal(response.Body(), deviceDetail)
	if err != nil {
		return nil, err
	}

	return deviceDetail, nil
}

func (client *syncClient) CreateDevice(request CreateDeviceRequest) (*CreateDeviceResponse, error) {
	bytesBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	response, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bytesBody).
		Post("/v5/iot/{project_id}/devices")

	if err != nil {
		return nil, err
	}

	resp := &CreateDeviceResponse{}
	err = json.Unmarshal(response.Body(), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (client *syncClient) ListDevices(queryParas map[string]string) (*ListDeviceResponse, error) {
	response, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(queryParas).
		Get("/v5/iot/{project_id}/devices")
	if err != nil {
		fmt.Println("list devices failed")
		return nil, err
	}

	if !successResponse(response) {
		fmt.Println("response failed")
		return nil, err
	}

	devices := &ListDeviceResponse{}

	err = json.Unmarshal(response.Body(), devices)
	if err != nil {
		fmt.Println("un marshal failed")
		return nil, err
	}

	return devices, nil
}

func (client *syncClient) UpdateDeviceProperties(deviceId string, services interface{}) (bool, error) {
	response, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		SetBody(services).
		Put("/v5/iot/{project_id}/devices/{device_id}/properties")

	if err != nil {
		return false, nil
	}

	return response.StatusCode() == 200, nil
}

func (client *syncClient) QueryDeviceProperties(deviceId, serviceId string) (interface{}, error) {
	httpResponse, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("service_id", serviceId).
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Get("/v5/iot/{project_id}/devices/{device_id}/properties")

	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode() != 200 {
		return nil, convertResponseToApplicationError(httpResponse)
	}

	var response interface{}

	err = json.Unmarshal(httpResponse.Body(), response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *syncClient) SendDeviceSyncCommand(deviceId string, request DeviceSyncCommandRequest) (*DeviceSyncCommandResponse, error) {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	response, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Post("/v5/iot/{project_id}/devices/{device_id}/commands")
	if err != nil {
		return nil, err
	}

	resp := &DeviceSyncCommandResponse{}
	err = json.Unmarshal(response.Body(), resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *syncClient) SendDeviceMessage(deviceId string, msg SendDeviceMessageRequest) (*SendDeviceMessageResponse, error) {
	reqBody, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	response, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Post("/v5/iot/{project_id}/devices/{device_id}/messages")
	if err != nil {
		return nil, err
	}

	resp := &SendDeviceMessageResponse{}

	err = json.Unmarshal(response.Body(), resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *syncClient) ListDeviceMessages(deviceId string) (*DeviceMessages, error) {
	response, err := client.client.R().
		SetPathParams(map[string]string{
			"device_id": deviceId,
		}).
		Get("/v5/iot/{project_id}/devices/{device_id}/messages")
	if err != nil {
		return &DeviceMessages{}, err
	}

	messages := &DeviceMessages{}
	err = json.Unmarshal(response.Body(), messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (client *syncClient) ShowDeviceMessage(deviceId, messageId string) (*DeviceMessage, error) {
	response, err := client.client.R().
		SetPathParams(map[string]string{
			"device_id":  deviceId,
			"message_id": messageId,
		}).
		Get("/v5/iot/{project_id}/devices/{device_id}/messages/{message_id}")
	if err != nil {
		return nil, err
	}

	messages := &DeviceMessage{}
	err = json.Unmarshal(response.Body(), messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (client *syncClient) ListApplications() (*Applications, error) {
	response, err := client.client.R().Get("/v5/iot/{project_id}/apps")
	if err != nil {
		return nil, err
	}

	app := &Applications{}
	err = json.Unmarshal(response.Body(), app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (client *syncClient) ShowApplication(appId string) (*Application, error) {
	response, err := client.client.R().
		SetPathParams(map[string]string{
			"app_id": appId,
		}).
		Get("/v5/iot/{project_id}/apps/{app_id}")
	if err != nil {
		return nil, err
	}

	app := &Application{}
	err = json.Unmarshal(response.Body(), app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (client *syncClient) DeleteApplication(appId string) (bool, error) {
	response, err := client.client.R().
		SetPathParams(map[string]string{
			"app_id": appId,
		}).
		Delete("/v5/iot/{project_id}/apps/{app_id}")
	if err != nil {
		return false, err
	}

	if response.StatusCode() != 204 {
		return false, convertResponseToApplicationError(response)
	}

	return true, nil
}

func (client *syncClient) CreateApplication(request ApplicationCreateRequest) (*Application, error) {
	body, err := json.Marshal(request)
	if err != nil {
		fmt.Println("marshal application create request failed")
		return nil, err
	}

	response, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("/v5/iot/{project_id}/apps")
	if err != nil {
		return nil, err
	}

	app := &Application{}
	err = json.Unmarshal(response.Body(), app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func CreateSyncIotApplicationClient(options ApplicationOptions) *syncClient {
	c := &syncClient{

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
