GetRow
=========
	
	// 说明：获取一行数据。
	//
	// ``table_name``是对应的表名。
	// ``primary_key``是主键，类型为``otstype.OTSPrimaryKey``。
	// ``columns_to_get``是可选参数，表示要获取的列的名称列表，类型为``otstype.OTSColumnsToGet``；如果填nil，表示获取所有列。
	//
	// 返回：本次操作消耗的CapacityUnit、行数据（包含主键列和属性列）。
	//       错误信息。
	//
	// ``get_row_response``为``otstype.OTSGetRowResponse``类的实例包含了：
	// ``Consumed``表示消耗的CapacityUnit，是``otstype.OTSCapacityUnit``类的实例。
	// ``Row``表示一行的数据，是``otstype.OTSRow``的实例,也包含了:
	// ``PrimaryKeyColumns``表示主键列，类型为``otstype.OTSPrimaryKey``，如：{"PK0":value0, "PK1":value1}。
	// ``AttributeColumns``表示属性列，类型为``otstype.OTSAttribute``，如：{"COL0":value0, "COL1":value1}。
	//
	// 示例：
	//
	// primary_key := &OTSPrimaryKey{
	// 	"gid": 1,
	// 	"uid": 101,
	// }
	// columns_to_get := &OTSColumnsToGet{
	// 	"name", "address", "age",
	// }
	// // columns_to_get = nil // read all
	// get_row_response, ots_err := ots_client.GetRow("myTable", primary_key, columns_to_get)
	//
	func (o *OTSClient) GetRow(table_name string, primary_key *OTSPrimaryKey, columns_to_get *OTSColumnsToGet) (get_row_response *OTSGetRowResponse, err *OTSError)

Example
=======
[GetRow.go](https://github.com/GiterLab/goots/blob/master/example/6-GetRow.go)

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
	
		// get_row
		primary_key := &OTSPrimaryKey{
			"gid": 1,
			"uid": 101,
		}
		columns_to_get := &OTSColumnsToGet{
			"name", "address", "age",
		}
		// columns_to_get = nil // read all
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