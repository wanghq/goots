// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// testcase for encoder

package coder

import (
	"fmt"
	"testing"

	. "github.com/GiterLab/goots/otstype"
	. "github.com/GiterLab/goots/protobuf"
)

func Test_encode_create_table(t *testing.T) {
	t.Log("testing _encode_create_table...")
	// ----
	table_meta := OTSTableMeta{
		TableName: "myTable",
		SchemaOfPrimaryKey: OTSSchemaOfPrimaryKey{
			"gid": "INTEGER",
			"uid": "INTEGER",
		},
	}

	reserved_throughput := OTSReservedThroughput{
		OTSCapacityUnit{0, 0},
	}

	req, _ := _encode_create_table(&table_meta, &reserved_throughput)
	t.Log("CreateTableRequest:", req)
	// ----
	t.Log("test _encode_create_table ok!")
	// t.Fail()
}

func Test_encode_delete_table(t *testing.T) {
	t.Log("testing _encode_delete_table...")
	// ----
	req, _ := _encode_delete_table("myTable")
	t.Log("DeleteTableRequest:", req)
	// ----
	t.Log("test _encode_delete_table ok!")
	// t.Fail()
}

func Test_encode_list_table(t *testing.T) {
	t.Log("testing _encode_list_table...")
	// ----
	req, _ := _encode_list_table()
	t.Log("ListTableRequest:", req)
	// ----
	t.Log("test _encode_list_table ok!")
	// t.Fail()
}

func Test_encode_update_table(t *testing.T) {
	t.Log("testing _encode_update_table...")
	// ----
	reserved_throughput := OTSReservedThroughput{
		OTSCapacityUnit{0, 0},
	}
	req, _ := _encode_update_table("myTable", &reserved_throughput)
	t.Log("UpdateTableRequest:", req)
	// ----
	t.Log("test _encode_update_table ok!")
	// t.Fail()
}

func Test_EncodeRequest(t *testing.T) {
	t.Log("testing EncodeRequest...")
	// ----
	table_meta := OTSTableMeta{
		TableName: "myTable",
		SchemaOfPrimaryKey: OTSSchemaOfPrimaryKey{
			"gid": "INTEGER",
			"uid": "INTEGER",
		},
	}

	reserved_throughput := OTSReservedThroughput{
		OTSCapacityUnit{0, 0},
	}

	req, err := EncodeRequest("CreateTable", &table_meta, &reserved_throughput)
	if err != nil {
		t.Logf("EncodeRequest error: %s", err)
		t.Fail()
	}

	if len(req) != 2 {
		t.Fail()
	} else {
		if req[1].Interface() != nil {
			err, ok := req[1].Interface().(error)
			if ok {
				if err != nil {
					t.Logf("err: %s", err)
					t.Fail()
				}
			} else {
				t.Log("Illegal data parameters, parse err failed")
				t.Fail()
			}
		}
	}

	if req[0].Interface() != nil {
		v, ok := req[0].Interface().(*CreateTableRequest)
		if ok {
			if v != nil {
				// fmt.Println("TableName", v.GetTableMeta().GetTableName())
				if v.GetTableMeta().GetTableName() != "myTable" {
					t.Log("TableName:", v.GetTableMeta().GetTableName())
					t.Fail()
				}
			} else {
				fmt.Println("CreateTableRequest is nil")
			}
		} else {
			fmt.Println("CreateTableRequest error")
		}
	}

	// ----
	t.Log("test EncodeRequest ok!")
	// t.Fail()
}
