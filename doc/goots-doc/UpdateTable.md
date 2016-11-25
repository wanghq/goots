UpdateTable
=========
	
	// 说明：更新表属性，目前只支持修改预留读写吞吐量。
	//
	// ``table_name``是对应的表名。
	// ``reserved_throughput``是``otstype.ReservedThroughput``类的实例，表示预留读写吞吐量。
	//
	// 返回：针对该表的预留读写吞吐量的最近上调时间、最近下调时间和当天下调次数。
	//       错误信息。
	//
	// ``update_table_response``表示更新的结果，是``otstype.OTSUpdateTableResponse``类的实例。
	//
	// 示例：
	// reserved_throughput := &OTSReservedThroughput{
	//  OTSCapacityUnit{0, 0},
	// }
	//
	// // 每次调整操作的间隔应大于10分钟
	// // 如果是刚创建表，需要10分钟之后才能调整表的预留读写吞吐量。
	// update_response, ots_err := ots_client.UpdateTable("myTable", reserved_throughput)
	//
	func (o *OTSClient) UpdateTable(table_name string, reserved_throughput *OTSReservedThroughput) (update_table_response *OTSUpdateTableResponse, err *OTSError)

Example
=======
[UpdateTable.go](https://github.com/GiterLab/goots/blob/master/example/4-UpdateTable.go)

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
		reserved_throughput := &OTSReservedThroughput{
			OTSCapacityUnit{0, 0},
		}
	
		// 每次调整操作的间隔应大于10分钟
		// 如果是刚创建表，需要10分钟之后才能调整表的预留读写吞吐量。
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