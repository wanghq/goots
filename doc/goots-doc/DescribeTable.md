DescribeTable
=========
	
	// 说明：获取表的描述信息。
	//
	// ``table_name``是对应的表名。
	//
	// 返回：表的描述信息。
	//
	// ``describe_table_response``表示表的描述信息，是``otstype.OTSDescribeTableResponse``类的实例。
	//
	// 示例：
	//
	// describe_response, ots_err := ots_client.DescribeTable("myTable")
	//
	func (o *OTSClient) DescribeTable(table_name string) (describe_table_response *OTSDescribeTableResponse, err *OTSError)

Example
=======
[DescribeTable.go](https://github.com/GiterLab/goots/blob/master/example/5-DescribeTable.go)

	package main
	
	import (
		"fmt"
		"os"
	
		ots2 "github.com/GiterLab/goots"
		"github.com/GiterLab/goots/log"
		// . "github.com/GiterLab/goots/otstype"
	)
	
	// modify it to yours
	const (
		ENDPOINT     = "http://127.0.0.1:8800"
		ACCESSID     = "OTSMultiUser177_accessid"
		ACCESSKEY    = "OTSMultiUser177_accesskey"
		INSTANCENAME = "TestInstance177"
	)
	
	func main() {
		// set running environment
		ots2.OTSDebugEnable = true
		ots2.OTSLoggerEnable = true
		log.OTSErrorPanicMode = true // 默认为开启，如果不喜欢panic则设置此为false
	
		fmt.Println("Test goots start ...")
	
		ots_client, err := ots2.New(ENDPOINT, ACCESSID, ACCESSKEY, INSTANCENAME)
		if err != nil {
			fmt.Println(err)
		}
	
		// describe_table
		describe_response, ots_err := ots_client.DescribeTable("myTable")
		if ots_err != nil {
			fmt.Println(ots_err)
			os.Exit(1)
		}
		fmt.Println("表的名称:", describe_response.TableMeta.TableName)
		fmt.Println("表的主键:", describe_response.TableMeta.SchemaOfPrimaryKey)
		fmt.Println("表的预留读吞吐量:", describe_response.ReservedThroughputDetails.CapacityUnit.Read)
		fmt.Println("表的预留写吞吐量:", describe_response.ReservedThroughputDetails.CapacityUnit.Write)
		fmt.Println("最后一次上调预留读写吞吐量时间:", describe_response.ReservedThroughputDetails.LastIncreaseTime)
		fmt.Println("最后一次下调预留读写吞吐量时间:", describe_response.ReservedThroughputDetails.LastDecreaseTime)
		fmt.Println("UTC自然日内总的下调预留读写吞吐量次数:", describe_response.ReservedThroughputDetails.NumberOfDecreasesToday)
	}