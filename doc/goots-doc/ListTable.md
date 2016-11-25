ListTable
=========
	
	// 说明：获取所有表名的列表。
	//
	// 返回：表名列表。
	//       错误信息。
	//
	// ``table_list``表示获取的表名列表，类型为OTSListTableResponse。
	//
	// 示例：
	//
	//     table_list, ots_err := ots_client.ListTable()
	//
	func (o *OTSClient) ListTable() (table_list *OTSListTableResponse, err *OTSError)

Example
=======
[ListTable.go](https://github.com/GiterLab/goots/blob/master/example/3-ListTable.go)

	package main
	
	import (
		"fmt"
		"os"
	
		ots2 "github.com/GiterLab/goots"
		// . "github.com/GiterLab/goots/otstype"
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
	
		// list_table
		list_tables, ots_err := ots_client.ListTable()
		if ots_err != nil {
			fmt.Println(ots_err)
			os.Exit(1)
		}
		fmt.Println("表的列表如下：")
		fmt.Println("list_tables:", list_tables.TableNames)
	}