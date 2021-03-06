package iot

type CreateDeviceGroupRequest struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	SuperGroupID string `json:"super_group_id,omitempty"`
	AppID        string `json:"app_id,omitempty"`
}

type CreateDeviceGroupResponse struct {
	GroupID      string `json:"group_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	SuperGroupID string `json:"super_group_id"`
}

type ShowDeviceGroupResponse struct {
	GroupID      string `json:"group_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	SuperGroupID string `json:"super_group_id"`
}

type UpdateDeviceGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateDeviceGroupResponse struct {
	GroupID      string `json:"group_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	SuperGroupID string `json:"super_group_id"`
}

type ListDeviceGroupRequest struct {
	Limit            int    `json:"limit,omitempty"`
	Marker           string `json:"marker,omitempty"`
	Offset           int    `json:"offset,omitempty"`
	LastModifiedTime string `json:"last_modified_time,omitempty"`
	AppId            string `json:"app_id,omitempty"`
}

type ListDeviceGroupResponse struct {
	DeviceGroups []DeviceGroupResponseDTO `json:"device_groups"`
	Page         Page
}



type DeviceGroupResponseDTO struct {
	GroupId      string `json:"group_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	SuperGroupId string `json:"super_group_id"`
}

type ListDeviceInDeviceGroupRequest struct {
	Limit  int    `json:"limit"`
	Marker string `json:"marker"`
	Offset int    `json:"offset"`
}

type ListDeviceInDeviceGroupResponse struct {
	Devices []SimplifyDevice `json:"devices"`
	Page    Page             `json:"page"`
}

type SimplifyDevice struct {
	DeviceID   string `json:"device_id"`
	NodeID     string `json:"node_id"`
	DeviceName string `json:"device_name"`
	ProductID  string `json:"product_id"`
}
