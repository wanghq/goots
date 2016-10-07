goots
=====
[![Gitter](https://badges.gitter.im/Join Chat.svg)](https://gitter.im/GiterLab/goots?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Aliyun OTS <del>(Open Table Service)</del> golang SDK.

OTS现更名为表格存储(Table Store)

[![wercker status](https://app.wercker.com/status/08d83208aa0215a6d6a0383b9b77b81d/m "wercker status")](https://app.wercker.com/project/bykey/08d83208aa0215a6d6a0383b9b77b81d)

[![Build Status](https://travis-ci.org/GiterLab/goots.svg?branch=master)](https://travis-ci.org/GiterLab/goots)
[![GoDoc](http://godoc.org/github.com/GiterLab/goots?status.svg)](http://godoc.org/github.com/GiterLab/goots)

[![Build Status](https://drone.io/github.com/GiterLab/goots/status.png)](https://drone.io/github.com/GiterLab/goots/latest)
[![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/GiterLab/goots)
[![](http://gocover.io/_badge/github.com/GiterLab/goots)](http://gocover.io/github.com/GiterLab/goots)

## Support API
- **Table**
	- [CreateTable](https://github.com/GiterLab/goots/blob/master/doc/goots-doc/CreateTable.md) ☑
	- [DeleteTable](https://github.com/GiterLab/goots/blob/master/doc/goots-doc/DeleteTable.md) ☑
	- [ListTable](https://github.com/GiterLab/goots/blob/master/doc/goots-doc/ListTable.md) ☑
	- [UpdateTable](https://github.com/GiterLab/goots/blob/master/doc/goots-doc/UpdateTable.md) ☑
	- [DescribeTable](https://github.com/GiterLab/goots/blob/master/doc/goots-doc/DescribeTable.md) ☑
- **SingleRow**
	- [GetRow](https://github.com/GiterLab/goots/blob/master/doc/goots-doc/GetRow.md) ☑
	- [PutRow](https://github.com/GiterLab/goots/blob/master/doc/goots-doc/PutRow.md) ☑
	- [UpdateRow](https://github.com/GiterLab/goots/blob/master/doc/goots-doc/UpdateRow.md) ☑
	- [DeleteRow](https://github.com/GiterLab/goots/blob/master/doc/goots-doc/DeleteRow.md) ☑
- **BatchRow**
	- [BatchGetRow](https://github.com/GiterLab/goots/blob/master/doc/goots-doc/BatchGetRow.md) ☑
	- [BatchWriteRow](https://github.com/GiterLab/goots/blob/master/doc/goots-doc/BatchWriteRow.md) ☑
	- [GetRange](https://github.com/GiterLab/goots/blob/master/doc/goots-doc/GetRange.md) ☑
	- <del>XGetRange</del>

## Install

	$ go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
	$ go get github.com/GiterLab/goots

## Usage

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

		// delete a table
		ots_err := ots_client.DeleteTable("myTable")
		if ots_err != nil {
			fmt.Println(ots_err)
			// os.Exit(1)
		}
		fmt.Println("表已删除")

		// create a table
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

		ots_err = ots_client.CreateTable(table_meta, reserved_throughput)
		if ots_err != nil {
			fmt.Println(ots_err)
			os.Exit(1)
		}
		fmt.Println("表已创建")

		// list tables
		list_tables, ots_err := ots_client.ListTable()
		if ots_err != nil {
			fmt.Println(ots_err)
			os.Exit(1)
		}
		fmt.Println("表的列表如下：")
		fmt.Println("list_tables:", list_tables.TableNames)

		// insert a row
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

		// get a row
		primary_key = &OTSPrimaryKey{
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

More examples, please see [example/interfaces.go](https://github.com/GiterLab/goots/blob/master/example/interfaces.go).

## Links 
- [Open Table Service，OTS(表格存储)](http://www.aliyun.com/product/ots)
- [OTS介绍](https://help.aliyun.com/document_detail/27280.html?spm=5176.7838592.6.103.Hlwl1P)
- [OTS产品文档](https://help.aliyun.com/product/27278.html?spm=5176.doc27304.3.1.OJe3Hd)
- [使用API开发指南](http://help.aliyun.com/view/11108328_13761831.html?spm=5176.383723.9.6.RYJAsQ)
- [API操作概览](https://help.aliyun.com/document_detail/27304.html?spm=5176.7838567.6.143.BzmR39)
- [Python SDK开发包](http://oss.aliyuncs.com/aliyun_portal_storage/help/ots/ots_python_sdk_2.0.2.zip?spm=5176.383723.9.8.RYJAsQ&file=ots_python_sdk_2.0.2.zip)
- [Java SDK开发包](http://oss.aliyuncs.com/aliyun_portal_storage/help/ots/aliyun-openservices-OTS-2.0.4.zip?spm=5176.383723.9.9.RYJAsQ&file=aliyun-openservices-OTS-2.0.4.zip)
- [nodejs SDK](https://github.com/alibaba/ots)

## License

This project is under the MIT License. See the [LICENSE](https://github.com/GiterLab/goots/blob/master/LICENSE) file for the full license text.
