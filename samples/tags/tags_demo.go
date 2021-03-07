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

	tag := iot.TagV5DTO{
		TagValue: "tag-value",
		TagKey:   "tag-key",
	}
	resp, err := client.DeviceBindTags(iot.DeviceBindTagsRequest{
		ResourceType: "device",
		ResourceID:   "5fdb75cccbfe2f02ce81d4bf_go-app",
		Tags:         []iot.TagV5DTO{tag},
	})
	if err != nil {
		fmt.Println(err)
		panic(0)
	}

	fmt.Println(resp)

	result, err := client.ListDeviceByTags(iot.ListDeviceByTagsRequest{
		ResourceType: "device",
		Tags:         []iot.TagV5DTO{tag},
	})

	fmt.Println(err)

	fmt.Println(result)
}
