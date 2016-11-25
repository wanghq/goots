// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// example for ots2
package main

import (
	"fmt"
	"os"

	ots2 "github.com/GiterLab/goots"
	. "github.com/GiterLab/goots/otstype"
)

// modify it to yours
const (
	ENDPOINT     = "your_instance_address"
	ACCESSID     = "your_accessid"
	ACCESSKEY    = "your_accesskey"
	INSTANCENAME = "your_instance_name"
)

func main() {
	// set running environment
	ots2.OTSDebugEnable = true
	ots2.OTSLoggerEnable = true
	ots2.OTSErrorPanicMode = true // 默认为开启，如果不喜欢panic则设置此为false

	fmt.Println("Test goots start ...")

	ots_client, err := ots2.New(ENDPOINT, ACCESSID, ACCESSKEY, INSTANCENAME)
	if err != nil {
		fmt.Println(err)
	}

	// update_table
	// 容量型实例: 只能设置为0
	// 高性能实例: 可设置为其他值（0-5000）
	reserved_throughput := &OTSReservedThroughput{
		OTSCapacityUnit{0, 0},
	}

	// 每次调整操作的间隔应大于10分钟
	// 如果是刚创建表，需要2分钟之后才能调整表的预留读写吞吐量。
	update_response, ots_err := ots_client.UpdateTable("myTable", reserved_throughput)
	if ots_err != nil {
		fmt.Println(ots_err)
		os.Exit(1)
	}
	fmt.Println("表的预留读吞吐量:", update_response.ReservedThroughputDetails.CapacityUnit.Read)
	fmt.Println("表的预留写吞吐量:", update_response.ReservedThroughputDetails.CapacityUnit.Write)
	fmt.Println("最后一次上调预留读写吞吐量时间:", update_response.ReservedThroughputDetails.LastIncreaseTime)
	fmt.Println("最后一次下调预留读写吞吐量时间:", update_response.ReservedThroughputDetails.LastDecreaseTime)
	fmt.Println("UTC自然日内总的下调预留读写吞吐量次数:", update_response.ReservedThroughputDetails.NumberOfDecreasesToday)
}
