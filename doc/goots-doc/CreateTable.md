CreateTable
=========
	
	// 说明：根据表信息创建表。
	//
	// ``table_meta``是``otstype.OTSTableMeta``类的实例，它包含表名和PrimaryKey的schema，
	// 请参考``OTSTableMeta``类的文档。当创建了一个表之后，通常要等待1分钟时间使partition load
	// 完成，才能进行各种操作。
	// ``reserved_throughput``是``otstype.ReservedThroughput``类的实例，表示预留读写吞吐量。
	//
	// 返回：无。
	//       错误信息。
	//
	// 示例：
	//
	// table_meta := &OTSTableMeta{
	// 	TableName: "myTable",
	// 	SchemaOfPrimaryKey: OTSSchemaOfPrimaryKey{
	//		{K: "gid", V: "INTEGER"},
	//		{K: "uid", V: "INTEGER"},
	// 	},
	// }
	//
	// reserved_throughput := &OTSReservedThroughput{
	// 	OTSCapacityUnit{0, 0},
	// }
	//
	// ots_err := ots_client.CreateTable(table_meta, reserved_throughput)
	//
	func (o *OTSClient) CreateTable(table_meta *OTSTableMeta, reserved_throughput *OTSReservedThroughput) (err *OTSError)

Example
=======
[CreateTable.go](https://github.com/GiterLab/goots/blob/master/example/1-CreateTable.go)

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
	
		// create_table
		table_meta := &OTSTableMeta{
			TableName: "myTable",
			SchemaOfPrimaryKey: OTSSchemaOfPrimaryKey{
				{K: "gid", V: "INTEGER"},
				{K: "uid", V: "INTEGER"},
			},
		}
	
		reserved_throughput := &OTSReservedThroughput{
			OTSCapacityUnit{0, 0},
		}
	
		ots_err := ots_client.CreateTable(table_meta, reserved_throughput)
		if ots_err != nil {
			fmt.Println(ots_err)
			os.Exit(1)
		}
		fmt.Println("表已创建")
	}