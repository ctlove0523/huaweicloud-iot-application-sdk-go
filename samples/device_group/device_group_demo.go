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
			Token:   "xxx",
		},
	}

	client := iot.CreateSyncIotApplicationClient(options)

	result, _ := client.RemoveDeviceFromDeviceGroup("78fa8117-b39b-4466-8c3a-0920f9ddaf0e", "5fdb75cccbfe2f02ce81d4bf_go-app")

	fmt.Println(result)
	//
	//resp, err := client.ListDeviceGroups(iot.ListDeviceGroupRequest{})
	//if err != nil {
	//	fmt.Println(err)
	//	panic(0)
	//}
	//
	//fmt.Println(resp.DeviceGroups)
	//
	//response, err := client.CreateDeviceGroup(iot.CreateDeviceGroupRequest{
	//	Name:  "first",
	//	AppID: "a04cafa7d2714e9eaff4fe9b210ccec0",
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	panic(0)
	//}
	//
	//groupId := response.GroupID
	//
	//showDeviceResponse, err := client.ShowDeviceGroup(groupId)
	//if err != nil {
	//	fmt.Println(err)
	//	panic(0)
	//}
	//
	//fmt.Printf("group id is %s\n", showDeviceResponse.GroupID)
	//fmt.Printf("group name is %s\n", showDeviceResponse.Name)
	//fmt.Printf("group description is %s\n", showDeviceResponse.Description)
}
