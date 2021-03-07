package main

import (
	"fmt"
	iot "huaweicloud-iot-application-sdk-go"
)

func main() {
	options := iot.ApplicationOptions{
		ServerPort:    443,
		ServerAddress: "iotda.cn-north-4.myhuaweicloud.com",
		InstanceId:    "",
		ProjectId:     "25e1be7c374749e9b6a25bc4ad53393a",

		Credential: &iot.Credentials{
			UseAkSk: false,
			Token:   "xxx",
		},
	}

	client := iot.CreateSyncIotApplicationClient(options)

	apps, _ := client.ListApplications()
	fmt.Println(apps.Applications)

	fmt.Println(client.ShowApplication("a04cafa7d2714e9eaff4fe9b210ccec0"))
	//client.DeleteApplication("586523b4fb7b451691fd7ead979d4eed")

	fmt.Println(client.CreateApplication(iot.ApplicationCreateRequest{
		AppName: "gosdk",
	}))
}
