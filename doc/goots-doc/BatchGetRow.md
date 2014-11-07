BatchGetRow
=========
	
	// 说明：批量获取多行数据。
	//
	// ``batch_list``表示获取多行的条件列表，格式如下：
	//
	// batch_list := &OTSBatchGetRowRequest{
	// 	{
	// 		// TableName
	// 		TableName: "table_name0",
	// 		// PrimaryKey
	// 		Rows: OTSPrimaryKeyRows{
	// 			{"gid": 1, "uid": 101},
	// 			{"gid": 2, "uid": 202},
	// 			{"gid": 3, "uid": 303},
	// 		},
	// 		// ColumnsToGet
	// 		ColumnsToGet: OTSColumnsToGet{"name", "address", "mobile", "age"},
	// 	},
	// 	{
	// 		// TableName
	// 		TableName: "table_name1",
	// 		// PrimaryKey
	// 		Rows: OTSPrimaryKeyRows{
	// 			{"gid": 1, "uid": 101},
	// 			{"gid": 2, "uid": 202},
	// 			{"gid": 3, "uid": 303},
	// 		},
	// 		// ColumnsToGet
	// 		ColumnsToGet: OTSColumnsToGet{"name", "address", "mobile", "age"},
	// 	},
	// 	...
	// }
	//
	// 其中，Rows 为主键，类型为``otstype.OTSPrimaryKeyRows``。
	//
	// 返回：对应行的结果列表。
	//       错误信息
	//
	// ``response_rows_list``为``otstype.OTSBatchGetRowResponse``的实例
	// ``response_rows_list.Tables``为返回的结果列表，与请求的顺序一一对应，格式如下：
	// response_rows_list.Tables --> []*OTSTableInBatchGetRowResponseItem{
	// 	{
	// 		TableName: "table_name0",
	// 		Rows : []*OTSRowInBatchGetRowResponseItem{
	// 			row_data_item0, row_data_item1, ...
	// 		},
	// 	},
	// 	{
	// 		TableName: "table_name1",
	// 		Rows : []*OTSRowInBatchGetRowResponseItem{
	// 			row_data_item0, row_data_item1, ...
	// 		},
	// 	},
	// 	...
	// }
	//
	// 其中，row_data_item0, row_data_item1为``otstype.OTSRowInBatchGetRowResponseItem``的实例。
	//
	// 示例：
	//
	// batch_list_get := &OTSBatchGetRowRequest{
	// 	{
	// 		// TableName
	// 		TableName: "myTable",
	// 		// PrimaryKey
	// 		Rows: OTSPrimaryKeyRows{
	// 			{"gid": 1, "uid": 101},
	// 			{"gid": 2, "uid": 202},
	// 			{"gid": 3, "uid": 303},
	// 		},
	// 		// ColumnsToGet
	// 		ColumnsToGet: OTSColumnsToGet{"name", "address", "mobile", "age"},
	// 	},
	// 	{
	// 		// TableName
	// 		TableName: "notExistTable",
	// 		// PrimaryKey
	// 		Rows: OTSPrimaryKeyRows{
	// 			{"gid": 1, "uid": 101},
	// 			{"gid": 2, "uid": 202},
	// 			{"gid": 3, "uid": 303},
	// 		},
	// 		// ColumnsToGet
	// 		ColumnsToGet: OTSColumnsToGet{"name", "address", "mobile", "age"},
	// 	},
	// }
	// batch_get_response, ots_err := ots_client.BatchGetRow(batch_list_get)
	//
	func (o *OTSClient) BatchGetRow(batch_list *OTSBatchGetRowRequest) (response_rows_list *OTSBatchGetRowResponse, err *OTSError)

Example
=======
[BatchGetRow.go](https://github.com/GiterLab/goots/blob/master/example/10-BatchGetRow.go)

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
	
		// batch_get_row
		batch_list_get := &OTSBatchGetRowRequest{
			{
				// TableName
				TableName: "myTable",
				// PrimaryKey
				Rows: OTSPrimaryKeyRows{
					{"gid": 1, "uid": 101},
					{"gid": 2, "uid": 202},
					{"gid": 3, "uid": 303},
				},
				// ColumnsToGet
				ColumnsToGet: OTSColumnsToGet{"name", "address", "mobile", "age"},
			},
			{
				// TableName
				TableName: "notExistTable",
				// PrimaryKey
				Rows: OTSPrimaryKeyRows{
					{"gid": 1, "uid": 101},
					{"gid": 2, "uid": 202},
					{"gid": 3, "uid": 303},
				},
				// ColumnsToGet
				ColumnsToGet: OTSColumnsToGet{"name", "address", "mobile", "age"},
			},
		}
		batch_get_response, ots_err := ots_client.BatchGetRow(batch_list_get)
		if ots_err != nil {
			fmt.Println(ots_err)
			os.Exit(1)
		}
		if batch_get_response != nil {
			var succeed_total, failed_total, consumed_write_total int32
			for _, v := range batch_get_response.Tables {
				fmt.Println("操作的表名:", v.TableName)
				for i1, v1 := range v.Rows {
					if v1.IsOk {
						succeed_total = succeed_total + 1
						fmt.Println("   --第", i1, "行操作成功, 消耗读CapacityUnit为", v1.Consumed.GetRead())
						consumed_write_total = consumed_write_total + v1.Consumed.GetRead()
						// print get value
						fmt.Println(v1.Row)
					} else {
						failed_total = failed_total + 1
						if v1.Consumed == nil {
							fmt.Println("   --第", i1, "行操作失败, 消耗读CapacityUnit为", 0, "ErrorCode:", v1.ErrorCode, "ErrorMessage:", v1.ErrorMessage)
						} else {
							// 实际测试这里不会执行到
							fmt.Println("   --第", i1, "行操作失败, 消耗读CapacityUnit为", v1.Consumed.GetRead, "ErrorCode:", v1.ErrorCode, "ErrorMessage:", v1.ErrorMessage)
							consumed_write_total = consumed_write_total + v1.Consumed.GetRead()
						}
					}
				}
			}
			fmt.Printf("本次操作命中 %d 个, 失败 %d 个, 共消耗读CapacityUnit为 %d\n", succeed_total, failed_total, consumed_write_total)
		} else {
			fmt.Println("本次操作都失败，不消耗读CapacityUnit")
		}
	}