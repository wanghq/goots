DeleteRow
=========
	
	// 说明：删除一行数据。
	//
	// ``table_name``是对应的表名。
	// ``condition``表示执行操作前做条件检查，满足条件才执行，是string的实例。
	// 目前只支持对行的存在性进行检查，检查条件包括：'IGNORE'，'EXPECT_EXIST'和'EXPECT_NOT_EXIST'。
	// ``primary_key``表示主键，类型为``otstype.OTSPrimaryKey``的实例。
	//
	// 返回：本次操作消耗的CapacityUnit。
	//       错误信息。
	//
	// ``delete_row_response``为``otstype.OTSDeleteRowResponse``类的实例包含了：
	// ``Consumed``表示消耗的CapacityUnit，是``otstype.OTSCapacityUnit``类的实例。
	//
	// 示例：
	//
	// primary_key := &OTSPrimaryKey{
	// 	"gid": 1,
	// 	"uid": 101,
	// }
	// condition := OTSCondition_IGNORE
	// delete_row_response, ots_err := ots_client.DeleteRow("myTable", condition, primary_key)
	//
	func (o *OTSClient) DeleteRow(table_name string, condition string, primary_key *OTSPrimaryKey) (delete_row_response *OTSDeleteRowResponse, err *OTSError)

Example
=======
[DeleteRow.go](https://github.com/GiterLab/goots/blob/master/example/9-DeleteRow.go)

	package main
	
	import (
		"fmt"
		"os"
	
		ots2 "github.com/GiterLab/goots"
		"github.com/GiterLab/goots/log"
		. "github.com/GiterLab/goots/otstype"
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
	
		// delete_row
		primary_key := &OTSPrimaryKey{
			"gid": 1,
			"uid": 101,
		}
		condition := OTSCondition_IGNORE
		delete_row_response, ots_err := ots_client.DeleteRow("myTable", condition, primary_key)
		if ots_err != nil {
			fmt.Println(ots_err)
			os.Exit(1)
		}
		fmt.Println("成功删除数据，消耗的写CapacityUnit为:", delete_row_response.GetWriteConsumed())
	
		// get_row
		primary_key = &OTSPrimaryKey{
			"gid": 1,
			"uid": 101,
		}
		columns_to_get := &OTSColumnsToGet{
			"name", "address", "age",
		}
		columns_to_get = nil // read all
		get_row_response, ots_err := ots_client.GetRow("myTable", primary_key, columns_to_get)
		if ots_err != nil {
			fmt.Println(ots_err)
			os.Exit(1)
		}
		fmt.Println("成功读取数据，消耗的读CapacityUnit为:", get_row_response.GetReadConsumed())
		if get_row_response.Row != nil {
			if attribute_columns := get_row_response.Row.GetAttributeColumns(); attribute_columns != nil {
				fmt.Println("name信息:", attribute_columns.Get("name"))
				fmt.Println("address信息:", attribute_columns.Get("address"))
				fmt.Println("age信息:", attribute_columns.Get("age"))
				fmt.Println("mobile信息:", attribute_columns.Get("mobile"))
			} else {
				fmt.Println("未查询到数据")
			}
		} else {
			fmt.Println("未查询到数据")
		}
	}