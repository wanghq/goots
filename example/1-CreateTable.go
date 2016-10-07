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

	// create_table
	// 注意：OTS是按设置的ReservedThroughput计量收费，即使没有读写也会产生费用。
	table_meta := &OTSTableMeta{
		TableName: "myTable",
		SchemaOfPrimaryKey: OTSSchemaOfPrimaryKey{
			"gid": "INTEGER",
			"uid": "INTEGER",
		},
	}

	reserved_throughput := &OTSReservedThroughput{
		OTSCapacityUnit{9, 9},
	}

	ots_err := ots_client.CreateTable(table_meta, reserved_throughput)
	if ots_err != nil {
		fmt.Println(ots_err)
		os.Exit(1)
	}
	fmt.Println("表已创建")
}
