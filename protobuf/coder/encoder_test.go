// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// testcase for encoder

package coder

import (
	. "github.com/GiterLab/goots/otstype"
	"testing"
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
		OTSCapacityUnit{100, 100},
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
		OTSCapacityUnit{100, 100},
	}
	req, _ := _encode_update_table("myTable", &reserved_throughput)
	t.Log("UpdateTableRequest:", req)
	// ----
	t.Log("test _encode_update_table ok!")
	// t.Fail()
}
