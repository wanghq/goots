// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// example for ots2
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
