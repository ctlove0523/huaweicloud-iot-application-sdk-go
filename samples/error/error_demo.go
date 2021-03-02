package main

import (
	"fmt"
	iot "huaweicloud-iot-application-sdk-go"
)

func main() {
	e:=iot.ApplicationError{
		Status:    403,
		ErrorCode: "IoTDA>0001",
		ErrorMsg:  "test",
	}

	fmt.Println(e.Error())
}
