BatchWriteRow
=========
	
	// 说明：批量修改多行数据。
	//
	// ``batch_list``表示获取多行的条件列表，格式如下：
	//
	// batch_list := &OTSBatchWriteRowRequest{
	// 	{
	// 		TableName: "table_name0",
	// 		PutRows: OTSPutRows{
	// 			put_row_item, ...
	// 		},
	// 		UpdateRows: OTSUpdateRows{
	// 			update_row_item, ...
	// 		},
	// 		DeleteRows: OTSDeleteRows{
	// 			delete_row_item, ...
	// 		},
	// 	},
	// 	{
	// 		TableName: "table_name1",
	// 		PutRows: OTSPutRows{
	// 			put_row_item, ...
	// 		},
	// 		UpdateRows: OTSUpdateRows{
	// 			update_row_item, ...
	// 		},
	// 		DeleteRows: OTSDeleteRows{
	// 			delete_row_item, ...
	// 		},
	// 	},
	// 	...
	// }
	//
	// 其中，put_row_item, 是``otstype.OTSPutRows``类的实例；
	//       update_row_item, 是``otstype.OTSUpdateRows``类的实例；
	//       delete_row_item, 是``otstype.OTSDeleteRows``类的实例。
	//
	// 返回：对应行的修改结果列表。
	//       错误信息。
	//
	//
	// ``response_items_list``为``otstype.OTSBatchWriteRowResponse``的实例
	// ``response_items_list.Tables``为返回的结果列表，与请求的顺序一一对应，格式如下：
	// response_items_list.Tables --> []*OTSTableInBatchWriteRowResponseItem{
	// 	{
	// 		TableName: "table_name0", // for table_name0
	// 		PutRows: []*OTSRowInBatchWriteRowResponseItem{
	// 			put_row_resp, ...
	// 		},
	// 		UpdateRows: []*OTSRowInBatchWriteRowResponseItem{
	// 			update_row_resp, ...
	// 		},
	// 		DeleteRows: []*OTSRowInBatchWriteRowResponseItem{
	// 			delete_row_resp, ...
	// 		}
	// 	},
	// 	{
	// 		TableName: "table_name1", // for table_name1
	// 		PutRows: []*OTSRowInBatchWriteRowResponseItem{
	// 			put_row_resp, ...
	// 		},
	// 		UpdateRows: []*OTSRowInBatchWriteRowResponseItem{
	// 			update_row_resp, ...
	// 		},
	// 		DeleteRows: []*OTSRowInBatchWriteRowResponseItem{
	// 			delete_row_resp, ...
	// 		}
	// 	},
	// 	...
	// }
	//
	// 其中put_row_resp，update_row_resp和delete_row_resp都是``*otstype.OTSRowInBatchWriteRowResponseItem``类的实例。
	//
	// 示例：
	//
	// put_row_item := OTSPutRowItem{
	// 	Condition: OTSCondition_EXPECT_NOT_EXIST, // OTSCondition_IGNORE
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 2,
	// 		"uid": 202,
	// 	},
	// 	AttributeColumns: OTSAttribute{
	// 		"name":    "李四",
	// 		"address": "中国某地",
	// 		"age":     20,
	// 	},
	// }
	// // [2] update_row
	// update_row_item := OTSUpdateRowItem{
	// 	Condition: OTSCondition_IGNORE,
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 3,
	// 		"uid": 303,
	// 	},
	// 	UpdateOfAttributeColumns: OTSUpdateOfAttribute{
	// 		OTSOperationType_PUT: OTSColumnsToPut{
	// 			"name":    "李三",
	// 			"address": "中国某地",
	// 		},
	// 		OTSOperationType_DELETE: OTSColumnsToDelete{
	// 			"mobile", "age",
	// 		},
	// 	},
	// }
	// // [3] delete_row
	// delete_row_item := OTSDeleteRowItem{
	// 	Condition: OTSCondition_IGNORE,
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 4,
	// 		"uid": 404,
	// 	},
	// }
	// batch_list := &OTSBatchWriteRowRequest{
	// 	{
	// 		TableName: "myTable",
	// 		PutRows: OTSPutRows{
	// 			put_row_item,
	// 		},
	// 		UpdateRows: OTSUpdateRows{
	// 			update_row_item,
	// 		},
	// 		DeleteRows: OTSDeleteRows{
	// 			delete_row_item,
	// 		},
	// 	},
	// 	{
	// 		TableName: "notExistTable",
	// 		PutRows: OTSPutRows{
	// 			put_row_item,
	// 		},
	// 		UpdateRows: OTSUpdateRows{
	// 			update_row_item,
	// 		},
	// 		DeleteRows: OTSDeleteRows{
	// 			delete_row_item,
	// 		},
	// 	},
	// }
	// batch_write_response, ots_err := ots_client.BatchWriteRow(batch_list)
	//
	func (o *OTSClient) BatchWriteRow(batch_list *OTSBatchWriteRowRequest) (response_item_list *OTSBatchWriteRowResponse, err *OTSError)

Example
=======
[BatchWriteRow.go](https://github.com/GiterLab/goots/blob/master/example/11-BatchWriteRow.go)

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
	
		// batch_write_row
		// [1] put row
		put_row_item := OTSPutRowItem{
			Condition: OTSCondition_EXPECT_NOT_EXIST, // OTSCondition_IGNORE
			PrimaryKey: OTSPrimaryKey{
				"gid": 2,
				"uid": 202,
			},
			AttributeColumns: OTSAttribute{
				"name":    "李四",
				"address": "中国某地",
				"age":     20,
			},
		}
		// [2] update_row
		update_row_item := OTSUpdateRowItem{
			Condition: OTSCondition_IGNORE,
			PrimaryKey: OTSPrimaryKey{
				"gid": 3,
				"uid": 303,
			},
			UpdateOfAttributeColumns: OTSUpdateOfAttribute{
				OTSOperationType_PUT: OTSColumnsToPut{
					"name":    "李三",
					"address": "中国某地",
				},
				OTSOperationType_DELETE: OTSColumnsToDelete{
					"mobile", "age",
				},
			},
		}
		// [3] delete_row
		delete_row_item := OTSDeleteRowItem{
			Condition: OTSCondition_IGNORE,
			PrimaryKey: OTSPrimaryKey{
				"gid": 4,
				"uid": 404,
			},
		}
		batch_list_write := &OTSBatchWriteRowRequest{
			{
				TableName: "myTable",
				PutRows: OTSPutRows{
					put_row_item,
				},
				UpdateRows: OTSUpdateRows{
					update_row_item,
				},
				DeleteRows: OTSDeleteRows{
					delete_row_item,
				},
			},
			{
				TableName: "notExistTable",
				PutRows: OTSPutRows{
					put_row_item,
				},
				UpdateRows: OTSUpdateRows{
					update_row_item,
				},
				DeleteRows: OTSDeleteRows{
					delete_row_item,
				},
			},
		}
		batch_write_response, ots_err := ots_client.BatchWriteRow(batch_list_write)
		if ots_err != nil {
			fmt.Println(ots_err)
			os.Exit(1)
		}
		// NOTE: 实际测试如果部分行操作失败，不消耗写CapacityUnit，而不是说明书写的整体失败
		if batch_write_response != nil {
			var succeed_total, failed_total, consumed_write_total int32
			for _, v := range batch_write_response.Tables {
				fmt.Println("操作的表名:", v.TableName)
				fmt.Println("操作 PUT:")
				if len(v.PutRows) != 0 {
					for i1, v1 := range v.PutRows {
						if v1.IsOk {
							succeed_total = succeed_total + 1
							fmt.Println("   --第", i1, "行操作成功, 消耗写CapacityUnit为", v1.Consumed.GetWrite())
							// NOTE: 为什么这里当条件设置为 OTSCondition_IGNORE, 同时这个put的PK值已经存在时
							// 一个put会消耗2个CapacityUnit呢???
							consumed_write_total = consumed_write_total + v1.Consumed.GetWrite()
						} else {
							failed_total = failed_total + 1
							if v1.Consumed == nil {
								fmt.Println("   --第", i1, "行操作失败, 消耗写CapacityUnit为", 0, "ErrorCode:", v1.ErrorCode, "ErrorMessage:", v1.ErrorMessage)
							} else {
								// 实际测试这里不会执行到
								fmt.Println("   --第", i1, "行操作失败, 消耗写CapacityUnit为", v1.Consumed.GetWrite, "ErrorCode:", v1.ErrorCode, "ErrorMessage:", v1.ErrorMessage)
								consumed_write_total = consumed_write_total + v1.Consumed.GetWrite()
							}
						}
					}
				}
				fmt.Println("操作 Update:")
				if len(v.UpdateRows) != 0 {
					for i1, v1 := range v.UpdateRows {
						if v1.IsOk {
							succeed_total = succeed_total + 1
							fmt.Println("   --第", i1, "行操作成功, 消耗写CapacityUnit为", v1.Consumed.GetWrite())
							consumed_write_total = consumed_write_total + v1.Consumed.GetWrite()
						} else {
							failed_total = failed_total + 1
							if v1.Consumed == nil {
								fmt.Println("   --第", i1, "行操作失败, 消耗写CapacityUnit为", 0, "ErrorCode:", v1.ErrorCode, "ErrorMessage:", v1.ErrorMessage)
							} else {
								// 实际测试这里不会执行到
								fmt.Println("   --第", i1, "行操作失败, 消耗写CapacityUnit为", v1.Consumed.GetWrite, "ErrorCode:", v1.ErrorCode, "ErrorMessage:", v1.ErrorMessage)
								consumed_write_total = consumed_write_total + v1.Consumed.GetWrite()
							}
						}
					}
				}
				fmt.Println("操作 Delete:")
				if len(v.DeleteRows) != 0 {
					for i1, v1 := range v.DeleteRows {
						if v1.IsOk {
							succeed_total = succeed_total + 1
							fmt.Println("   --第", i1, "行操作成功, 消耗写CapacityUnit为", v1.Consumed.GetWrite())
							consumed_write_total = consumed_write_total + v1.Consumed.GetWrite()
						} else {
							failed_total = failed_total + 1
							if v1.Consumed == nil {
								fmt.Println("   --第", i1, "行操作失败, 消耗写CapacityUnit为", 0, "ErrorCode:", v1.ErrorCode, "ErrorMessage:", v1.ErrorMessage)
	
							} else {
								// 实际测试这里不会执行到
								fmt.Println("   --第", i1, "行操作失败, 消耗写CapacityUnit为", v1.Consumed.GetWrite, "ErrorCode:", v1.ErrorCode, "ErrorMessage:", v1.ErrorMessage)
								consumed_write_total = consumed_write_total + v1.Consumed.GetWrite()
							}
						}
					}
				}
			}
			fmt.Printf("本次操作命中 %d 个, 失败 %d 个, 共消耗写CapacityUnit为 %d\n", succeed_total, failed_total, consumed_write_total)
		} else {
			fmt.Println("本次操作都失败，不消耗写CapacityUnit")
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