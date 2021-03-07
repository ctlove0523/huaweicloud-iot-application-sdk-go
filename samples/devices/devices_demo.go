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
			Ak:      "xxx",
			Sk:      "xxx",
			UseAkSk: false,
			Token:   "",
		},
	}

	client := iot.CreateSyncIotApplicationClient(options)

	fmt.Println(client.ResetDeviceSecret("5fdb75cccbfe2f02ce81d4bf_1234567", "987654321", true))

}

func showDevice(client iot.ApplicationClient) {
	fmt.Println(client.ShowDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt"))
}
