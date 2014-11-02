ListTable
=========

	func (o *OTSClient) ListTable() (table_list *OTSListTableResponse, err *OTSError)

Example
=======

	package main
	
	import (
		"fmt"
		"os"
	
		ots2 "github.com/GiterLab/goots"
		"github.com/GiterLab/goots/log"
	)
	
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
	
		// list_table
		list_tables, ots_err := ots_client.ListTable()
		if ots_err != nil {
			fmt.Println(ots_err)
			os.Exit(1)
		}
		fmt.Println("表的列表如下：")
		fmt.Println("list_tables:", list_tables.TableNames)
	}