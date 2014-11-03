DeleteTable
=========

	// 说明：根据表名删除表。
	//
	// ``table_name``是对应的表名。
	//
	// 返回：无。
	//       错误信息。
	//
	// 示例：
	//
	// ots_client.DeleteTable("myTable")
	//
	func (o *OTSClient) DeleteTable(table_name string) (err *OTSError)

Example
=======
[DeleteTable.go](https://github.com/GiterLab/goots/blob/master/example/2-DeleteTable.go)

	package main
	
	import (
		"fmt"
		"os"
	
		ots2 "github.com/GiterLab/goots"
		"github.com/GiterLab/goots/log"
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
	
		// delete_table
		ots_err := ots_client.DeleteTable("myTable")
		if ots_err != nil {
			fmt.Println(ots_err)
			os.Exit(1)
		}
		fmt.Println("表已删除")
	}