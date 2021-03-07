package iot

type DeviceMessages struct {
	DeviceId string          `json:"device_id"`
	Messages []DeviceMessage `json:"messages"`
}

type DeviceMessage struct {
	MessageId    string `json:"message_id"`
	Name         string `json:"name"`
	Message      string `json:"message"`
	Topic        string `json:"topic"`
	Status       string `json:"status"`
	CreatedTime  string `json:"created_time"`
	FinishedTime string `json:"finished_time"`
}

type SendDeviceMessageRequest struct {
	MessageId     string `json:"message_id"`
	Name          string `json:"name"`
	Message       string `json:"message"`
	Topic         string `json:"topic"`
	TopicFullName string `json:"topic_full_name"`
}

type SendDeviceMessageResponse struct {
	MessageId string        `json:"message_id"`
	Result    MessageResult `json:"result"`
}

type MessageResult struct {
	Status       string `json:"status"`
	CreatedTime  string `json:"created_time"`
	FinishedTime string `json:"finished_time"`
}

type DeviceSyncCommandRequest struct {
	ServiceId   string      `json:"service_id"`
	CommandName string      `json:"command_name"`
	Paras       interface{} `json:"paras"`
}

type DeviceSyncCommandResponse struct {
	CommandId string      `json:"command_id"`
	Response  interface{} `json:"response"`
}

// 设备管理

type ListDeviceResponse struct {
	Devices []QueryDeviceSimplify `json:"devices"`
	Page    Page                  `json:"page"`
}

type QueryDeviceSimplify struct {
	AppID       string     `json:"app_id"`
	AppName     string     `json:"app_name"`
	DeviceID    string     `json:"device_id"`
	NodeID      string     `json:"node_id"`
	GatewayID   string     `json:"gateway_id"`
	DeviceName  string     `json:"device_name"`
	NodeType    string     `json:"node_type"`
	Description string     `json:"description"`
	FwVersion   string     `json:"fw_version"`
	SwVersion   string     `json:"sw_version"`
	ProductID   string     `json:"product_id"`
	ProductName string     `json:"product_name"`
	Status      string     `json:"status"`
	Tags        []TagV5DTO `json:"tags"`
}

type TagV5DTO struct {
	TagKey   string `json:"tag_key"`
	TagValue string `json:"tag_value"`
}

type Page struct {
	Count  int    `json:"count"`
	Marker string `json:"marker"`
}

// 设备管理-创建设备
type CreateDeviceRequest struct {
	DeviceID      string           `json:"device_id,omitempty"`
	NodeID        string           `json:"node_id"`
	DeviceName    string           `json:"device_name,omitempty"`
	ProductID     string           `json:"product_id"`
	AuthInfo      AuthInfo         `json:"auth_info,omitempty"`
	Description   string           `json:"description,omitempty"`
	GatewayID     string           `json:"gateway_id,omitempty"`
	AppID         string           `json:"app_id,omitempty"`
	ExtensionInfo interface{}      `json:"extension_info,omitempty"`
	Shadow        []InitialDesired `json:"shadow,omitempty"`
}

type CreateDeviceResponse struct {
	AppID         string      `json:"app_id"`
	AppName       string      `json:"app_name"`
	DeviceID      string      `json:"device_id"`
	NodeID        string      `json:"node_id"`
	GatewayID     string      `json:"gateway_id"`
	DeviceName    string      `json:"device_name"`
	NodeType      string      `json:"node_type"`
	Description   string      `json:"description"`
	FwVersion     string      `json:"fw_version"`
	SwVersion     string      `json:"sw_version"`
	AuthInfo      AuthInfo    `json:"auth_info"`
	ProductID     string      `json:"product_id"`
	ProductName   string      `json:"product_name"`
	Status        string      `json:"status"`
	CreateTime    string      `json:"create_time"`
	Tags          []TagV5DTO  `json:"tags"`
	ExtensionInfo interface{} `json:"extension_info"`
}

type AuthInfo struct {
	AuthType     string `json:"auth_type,omitempty"`
	SecureAccess bool   `json:"secure_access,omitempty"`
	Fingerprint  string `json:"fingerprint,omitempty"`
	Secret       string `json:"secret,omitempty"`
	Timeout      int    `json:"timeout,omitempty"`
}

type InitialDesired struct {
	Desired   interface{} `json:"desired"`
	ServiceID string      `json:"service_id"`
}

type DeviceDetailResponse struct {
	AppID         string      `json:"app_id"`
	AppName       string      `json:"app_name"`
	DeviceID      string      `json:"device_id"`
	NodeID        string      `json:"node_id"`
	GatewayID     string      `json:"gateway_id"`
	DeviceName    string      `json:"device_name"`
	NodeType      string      `json:"node_type"`
	Description   string      `json:"description"`
	FwVersion     string      `json:"fw_version"`
	SwVersion     string      `json:"sw_version"`
	AuthInfo      AuthInfo    `json:"auth_info"`
	ProductID     string      `json:"product_id"`
	ProductName   string      `json:"product_name"`
	Status        string      `json:"status"`
	CreateTime    string      `json:"create_time"`
	Tags          []TagV5DTO  `json:"tags"`
	ExtensionInfo interface{} `json:"extension_info"`
}

type UpdateDeviceRequest struct {
	DeviceName    string                `json:"device_name,omitempty"`
	Description   string                `json:"description,omitempty"`
	ExtensionInfo interface{}           `json:"extension_info,omitempty"`
	AuthInfo      AuthInfoWithoutSecret `json:"auth_info,omitempty"`
}

type AuthInfoWithoutSecret struct {
	SecureAccess bool `json:"secure_access,omitempty"`
	Timeout      int  `json:"timeout,omitempty"`
}

type ResetDeviceSecretResponse struct {
	DeviceId string `json:"device_id"`
	Secret   string `json:"secret"`
}
