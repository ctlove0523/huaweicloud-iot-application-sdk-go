package iot

type ListAmqpQueuesRequest struct {
	QueueName string `json:"queue_name,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Marker    string `json:"marker,omitempty"`
	Offset    string `json:"offset,omitempty"`
}

type ListAmqpQueuesResponse struct {
	Queues []QueryQueueBase `json:"queues"`
	Page   Page             `json:"page"`
}

type QueryQueueBase struct {
	QueueName      string `json:"queue_name"`
	CreateTime     string `json:"create_time"`
	LastModifyTime string `json:"last_modify_time"`
	QueueID        string `json:"queue_id"`
}

type CreateAmqpQueueResponse struct {
	QueueID        string `json:"queue_id"`
	QueueName      string `json:"queue_name"`
	CreateTime     string `json:"create_time"`
	LastModifyTime string `json:"last_modify_time"`
}

type ShowAmqpQueueResponse struct {
	QueueID        string `json:"queue_id"`
	QueueName      string `json:"queue_name"`
	CreateTime     string `json:"create_time"`
	LastModifyTime string `json:"last_modify_time"`
}
