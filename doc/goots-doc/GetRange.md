GetRange
=========

	// 说明：根据范围条件获取多行数据。
	//
	// ``table_name``是对应的表名。
	// ``direction``表示范围的方向，字符串格式，取值包括'FORWARD'和'BACKWARD'。
	// ``inclusive_start_primary_key``表示范围的起始主键（在范围内）。
	// ``exclusive_end_primary_key``表示范围的结束主键（不在范围内）。
	// ``columns_to_get``是可选参数，表示要获取的列的名称列表，类型为``otstype.OTSColumnsToGet``；如果为nil，表示获取所有列。
	// ``limit``是可选参数，表示最多读取多少行；如果为0，则没有限制。
	//
	// 返回：符合条件的结果列表。
	//       错误信息。
	//
	// ``response_row_list``为``otstype.OTSGetRangeResponse``类的实例包含了：
	// ``Consumed``表示消耗的CapacityUnit，是``otstype.OTSCapacityUnit``类的实例。
	// ``NextStartPrimaryKey``表示下次get_range操作的起始点的主健列，类型为``otstype.OTSPrimaryKey``。
	// ``Rows``表示本次操作返回的行数据列表，是``otstype.OTSRows``类的实例。
	//
	// 示例：
	//
	// // get_range
	// // 查询区间：[(1, INF_MIN), (4, INF_MAX))，左闭右开。
	// inclusive_start_primary_key := &OTSPrimaryKey{
	// 	"gid": 1,
	// 	"uid": OTSColumnType_INF_MIN,
	// }
	// exclusive_end_primary_key := &OTSPrimaryKey{
	// 	"gid": 4,
	// 	"uid": OTSColumnType_INF_MAX,
	// }
	// columns_to_get := &OTSColumnsToGet{
	// 	"gid", "uid", "name", "address", "mobile", "age",
	// }
	//
	// // 选择方向
	// // OTSDirection_FORWARD
	// // OTSDirection_BACKWARD
	// response_row_list, ots_err := ots_client.GetRange("myTable", OTSDirection_FORWARD,
	// 	inclusive_start_primary_key, exclusive_end_primary_key, columns_to_get, 100)
	//
	func (o *OTSClient) GetRange(table_name string, direction string,
		inclusive_start_primary_key *OTSPrimaryKey,
		exclusive_end_primary_key *OTSPrimaryKey,
		columns_to_get *OTSColumnsToGet,
		limit int32) (response_row_list *OTSGetRangeResponse, err *OTSError)

Example
=======
[GetRange.go](https://github.com/GiterLab/goots/blob/master/example/12-GetRange.go)

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
	
		// get_range
		// 查询区间：[(1, INF_MIN), (4, INF_MAX))，左闭右开。
		inclusive_start_primary_key := &OTSPrimaryKey{
			"gid": 1,
			"uid": OTSColumnType_INF_MIN,
		}
		exclusive_end_primary_key := &OTSPrimaryKey{
			"gid": 4,
			"uid": OTSColumnType_INF_MAX,
		}
		columns_to_get := &OTSColumnsToGet{
			"gid", "uid", "name", "address", "mobile", "age",
		}
	
		// 选择方向
		// OTSDirection_FORWARD
		// OTSDirection_BACKWARD
		response_row_list, ots_err := ots_client.GetRange("myTable", OTSDirection_FORWARD,
			inclusive_start_primary_key, exclusive_end_primary_key, columns_to_get, 100)
		if ots_err != nil {
			fmt.Println(ots_err)
			os.Exit(1)
		}
		if response_row_list.GetRows() != nil {
			for i, v := range response_row_list.GetRows() {
				if v.GetPrimaryKeyColumns() != nil {
					fmt.Println("第 ", i, " 行:",
						"gid:", v.GetPrimaryKeyColumns().Get("gid"),
						"uid:", v.GetPrimaryKeyColumns().Get("uid"))
				} else {
					fmt.Println("第 ", i, " 行:")
				}
	
				fmt.Println("    - name信息:", v.GetAttributeColumns().Get("name"))
				fmt.Println("    - address信息:", v.GetAttributeColumns().Get("address"))
				fmt.Println("    - age信息:", v.GetAttributeColumns().Get("age"))
				fmt.Println("    - mobile信息:", v.GetAttributeColumns().Get("mobile"))
			}
		}
		fmt.Println("成功读取数据，消耗的读CapacityUnit为:", response_row_list.Consumed.GetRead())
		if response_row_list.GetNextStartPrimaryKey() != nil {
			fmt.Println("还有数据未读取完毕，用户可以继续调用GetRange()进行读取")
			fmt.Println("下次开始的主键:", response_row_list.GetNextStartPrimaryKey())
		}
	}