package iot

type ShowDeviceShadowResponse struct {
	DeviceID string             `json:"device_id"`
	Shadow   []DeviceShadowData `json:"shadow"`
}

type DeviceShadowData struct {
	ServiceID string                 `json:"service_id"`
	Desired   DeviceShadowProperties `json:"desired"`
	Reported  DeviceShadowProperties `json:"reported"`
	Version   int                    `json:"version"`
}

type DeviceShadowProperties struct {
	Properties interface{} `json:"properties"`
	EventTime  string      `json:"event_time"`
}

type UpdateDeviceShadowRequest struct {
	Shadow []UpdateDeviceShadowDesired `json:"shadow"`
}

type UpdateDeviceShadowDesired struct {
	ServiceId string      `json:"service_id"`
	Desired   interface{} `json:"desired"`
	Version   int         `json:"version"`
}
