// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// example for ots2
package main

import (
	"fmt"
	"os"
	"time"

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

	// delete_table
	ots_err := ots_client.DeleteTable("myTable")
	if ots_err != nil {
		fmt.Println(ots_err)
		// os.Exit(1)
	}
	fmt.Println("表已删除")

	// create_table
	// 注意：OTS是按设置的ReservedThroughput计量收费，即使没有读写也会产生费用。
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

	ots_err = ots_client.CreateTable(table_meta, reserved_throughput)
	if ots_err != nil {
		fmt.Println(ots_err)
		os.Exit(1)
	}
	fmt.Println("表已创建")
	fmt.Println("表创建后需要等待一会，才能进行继续操作,否则可能会出现操作失败...")
	for cnt := 5; cnt != 0; cnt-- {
		fmt.Println(cnt)
		time.Sleep(1 * time.Second)
	}

	// list_table
	list_tables, ots_err := ots_client.ListTable()
	if ots_err != nil {
		fmt.Println(ots_err)
		os.Exit(1)
	}
	fmt.Println("表的列表如下：")
	fmt.Println("list_tables:", list_tables.TableNames)

	// update_table
	// 只能在高性能实例实例下测试，不要在容量型实例上测试，否则出错
	//
	// 每次调整操作的间隔应大于10分钟
	// 如果是刚创建表，需要2分钟之后才能调整表的预留读写吞吐量。
	// 注意：OTS是按设置的ReservedThroughput计量收费，即使没有读写也会产生费用。
	//
	// Note:
	// 容量型实例: The value of read capacity unit can only be 0
	//             The value of write capacity unit can only be 0.
	//             Your instance is forbidden to update capacity unit.
	// 高性能实例: at least one of read or write of CapacityUnit is required
	update_reserved_throughput := &OTSReservedThroughput{
		OTSCapacityUnit{0, 0},
	}

	fmt.Println("Need to sleep 2 Minute, be patient...")
	time.Sleep(2 * time.Minute)
	time.Sleep(5 * time.Second)
	update_response, ots_err := ots_client.UpdateTable("myTable", update_reserved_throughput)
	if ots_err != nil {
		fmt.Println(ots_err)
		os.Exit(1)
	}
	fmt.Println("表的预留读吞吐量:", update_response.ReservedThroughputDetails.CapacityUnit.Read)
	fmt.Println("表的预留写吞吐量:", update_response.ReservedThroughputDetails.CapacityUnit.Write)
	fmt.Println("最后一次上调预留读写吞吐量时间:", update_response.ReservedThroughputDetails.LastIncreaseTime)
	fmt.Println("最后一次下调预留读写吞吐量时间:", update_response.ReservedThroughputDetails.LastDecreaseTime)
	fmt.Println("UTC自然日内总的下调预留读写吞吐量次数:", update_response.ReservedThroughputDetails.NumberOfDecreasesToday)

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

	// put_row
	primary_key := &OTSPrimaryKey{
		"gid": 1,
		"uid": 101,
	}
	attribute_columns := &OTSAttribute{
		"name":    "张三",
		"mobile":  111111111,
		"address": "中国A地",
		"age":     20,
	}
	condition := OTSCondition_EXPECT_NOT_EXIST
	put_row_response, ots_err := ots_client.PutRow("myTable", condition, primary_key, attribute_columns)
	if ots_err != nil {
		fmt.Println(ots_err)
		os.Exit(1)
	}
	fmt.Println("成功插入数据，消耗的写CapacityUnit为:", put_row_response.GetWriteConsumed())

	// get_row
	primary_key = &OTSPrimaryKey{
		"gid": 1,
		"uid": 101,
	}
	columns_to_get := &OTSColumnsToGet{
		"name", "address", "age",
	}
	// columns_to_get = nil
	// read all
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

	// update_row
	primary_key = &OTSPrimaryKey{
		"gid": 1,
		"uid": 101,
	}
	update_of_attribute_columns := &OTSUpdateOfAttribute{
		OTSOperationType_PUT: OTSColumnsToPut{
			"name":    "张三丰",
			"address": "中国B地",
		},

		OTSOperationType_DELETE: OTSColumnsToDelete{
			"mobile", "age",
		},
	}
	condition = OTSCondition_EXPECT_EXIST
	update_row_response, ots_err := ots_client.UpdateRow("myTable", condition, primary_key, update_of_attribute_columns)
	if ots_err != nil {
		fmt.Println(ots_err)
		os.Exit(1)
	}
	fmt.Println("成功插入数据，消耗的写CapacityUnit为:", update_row_response.GetWriteConsumed())

	// get_row
	primary_key = &OTSPrimaryKey{
		"gid": 1,
		"uid": 101,
	}
	columns_to_get = &OTSColumnsToGet{
		"name", "address", "age",
	}
	columns_to_get = nil // read all
	get_row_response, ots_err = ots_client.GetRow("myTable", primary_key, columns_to_get)
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

	// delete_row
	primary_key = &OTSPrimaryKey{
		"gid": 1,
		"uid": 101,
	}
	condition = OTSCondition_IGNORE
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
	columns_to_get = &OTSColumnsToGet{
		"name", "address", "age",
	}
	columns_to_get = nil // read all
	get_row_response, ots_err = ots_client.GetRow("myTable", primary_key, columns_to_get)
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
	columns_to_get = &OTSColumnsToGet{
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

	// delete_table
	ots_err = ots_client.DeleteTable("myTable")
	if ots_err != nil {
		fmt.Println(ots_err)
	}
	fmt.Println("测试完毕，表已删除")
}
