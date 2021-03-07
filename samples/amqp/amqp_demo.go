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
			Ak:      "S4QUJL4COTKPPR2VIFTF",
			Sk:      "hRsE5wFm31FpjCmQjxx9vqcodn7eFgDuE8q6eq5W",
			UseAkSk: true,
		},
	}

	client := iot.CreateSyncIotApplicationClient(options)

	queues, err := client.ListAmqpQueues(iot.ListAmqpQueuesRequest{})
	if err != nil {
		fmt.Println(err)
		panic(1)
	}

	fmt.Println(queues.Queues)

}
