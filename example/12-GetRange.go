// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// example for ots2
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
	ENDPOINT     = "your_instance_address"
	ACCESSID     = "your_accessid"
	ACCESSKEY    = "your_accesskey"
	INSTANCENAME = "your_instance_name"
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
