package iot

type DeviceBindTagsRequest struct {
	ResourceType string     `json:"resource_type,omitempty"`
	ResourceID   string     `json:"resource_id,omitempty"`
	Tags         []TagV5DTO `json:"tags,omitempty"`
}

type DeviceUnBindTagsRequest struct {
	ResourceType string   `json:"resource_type,omitempty"`
	ResourceID   string   `json:"resource_id,omitempty"`
	TagKeys      []string `json:"tag_keys,omitempty"`
}

type ListDeviceByTagsRequest struct {
	Limit        int        `json:"limit,omitempty"`
	Marker        string     `json:"marker,omitempty"`
	Offset       int        `json:"offset,omitempty"`
	ResourceType string     `json:"resource_type,omitempty"`
	Tags         []TagV5DTO `json:"tags,omitempty"`
}

type ListDeviceByTagsResponse struct {
	Resources []ResourceDTO `json:"resources"`
	Page      Page          `json:"page"`
}

type ResourceDTO struct {
	ResourceID string `json:"resource_id"`
}
