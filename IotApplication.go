package iot

type Application struct {
	AppId      string `json:"app_id"`
	AppName    string `json:"app_name"`
	CreateTime string `json:"create_time"`
	DefaultApp bool   `json:"default_app"`
}

type Applications struct {
	Applications []Application `json:"applications"`
}