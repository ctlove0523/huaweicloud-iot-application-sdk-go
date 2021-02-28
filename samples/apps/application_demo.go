package main

import (
	"fmt"
	iot "huaweicloud-iot-application-sdk-go"
)

func main() {
	options := iot.ApplicationOptions{
		ServerPort:    443,
		ServerAddress: "iotda.cn-north-4.myhuaweicloud.com",
		InstanceId:    "6797fccc-e700-4d68-ad1a-7e516ddcb0cc",
		ProjectId:     "25e1be7c374749e9b6a25bc4ad53393a",

		Credential: &iot.Credentials{
			Ak:      "xxx",
			Sk:      "xxx",
			UseAkSk: false,
			Token:   "xxx",
		},
	}

	client := iot.CreateIotApplicationClient(options)

	apps := client.ListApplications("", "25e1be7c374749e9b6a25bc4ad53393a")
	fmt.Println(apps.Applications)
}
