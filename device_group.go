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
	Name         string `json:"name"`
	Description  string `json:"description"`
}

type UpdateDeviceGroupResponse struct {
	GroupID      string `json:"group_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	SuperGroupID string `json:"super_group_id"`
}

