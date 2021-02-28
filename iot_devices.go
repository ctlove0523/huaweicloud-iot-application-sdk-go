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
