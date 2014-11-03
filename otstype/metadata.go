// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// metadata for ots2
package otstype

import (
	"time"
)

const (
	// OTSColumnType
	OTSColumnType_INF_MIN = "INF_MIN" // only for GetRange
	OTSColumnType_INF_MAX = "INF_MAX" // only for GetRange
	OTSColumnType_INTEGER = "INTEGER"
	OTSColumnType_STRING  = "STRING"
	OTSColumnType_BOOLEAN = "BOOLEAN"
	OTSColumnType_DOUBLE  = "DOUBLE"
	OTSColumnType_BINARY  = "BINARY"

	// OTSRowExistenceExpectation
	OTSRowExistenceExpectation_IGNORE           = "IGNORE"
	OTSRowExistenceExpectation_EXPECT_EXIST     = "EXPECT_EXIST"
	OTSRowExistenceExpectation_EXPECT_NOT_EXIST = "EXPECT_NOT_EXIST"

	OTSCondition_IGNORE           = "IGNORE"
	OTSCondition_EXPECT_EXIST     = "EXPECT_EXIST"
	OTSCondition_EXPECT_NOT_EXIST = "EXPECT_NOT_EXIST"

	// UpdateRow
	// OTSOperationType
	OTSOperationType_PUT    = "PUT"
	OTSOperationType_DELETE = "DELETE"

	// GetRange
	// OTSDirection
	OTSDirection_FORWARD  = "FORWARD"
	OTSDirection_BACKWARD = "BACKWARD"
)

// 表示一个表的结构信息
type OTSTableMeta struct {
	// 该表的表名
	TableName string
	// 该表全部的主键列描述
	SchemaOfPrimaryKey OTSSchemaOfPrimaryKey //  map[string]string{"PK0": "STRING", "PK1": "INTEGER", ...}
}

// 表主键列描述
type OTSSchemaOfPrimaryKey map[string]string

// 表示一次操作消耗服务能力单元的值或是一个表的预留读写吞吐量的值
type OTSCapacityUnit struct {
	// 本次操作消耗的读服务能力单元或该表的读服务能力单元
	Read int32
	// 本次操作消耗的写服务能力单元或该表的写服务能力单元
	Write int32
}

// 表示一个表设置的预留读写吞吐量数值
type OTSReservedThroughput struct {
	// 表当前的预留读写吞吐量数值
	CapacityUnit OTSCapacityUnit
}

// 表示一个表的预留读写吞吐量信息
type OTSReservedThroughputDetails struct {
	// 该表的预留读写吞吐量的数值
	CapacityUnit OTSCapacityUnit
	// 最近一次上调该表的预留读写吞吐量设置的时间，使用UTC 秒数表示
	LastIncreaseTime time.Time
	// 最近一次下调该表的预留读写吞吐量设置的时间，使用UTC 秒数表示
	LastDecreaseTime time.Time
	// 一个自然日内已下调该表的预留读写吞吐量设置的次数
	NumberOfDecreasesToday int32
}

// 表示一列
type OTSColumn struct {
	// 该列的列名
	Name string
	// 该列的列值
	Value OTSColumnValue
}

// 表示一列的列值
type OTSColumnValue struct {
	// 该列的数据类型
	Type string
	// 该列的数据，只在type 为INTEGER 时有效
	VInt int64
	// 该列的数据，只在type 为STRING 时有效，必须为UTF-8 编码
	// golang默认为UTF-8编码
	VString string
	// 该列的数据，只在type 为BOOLEAN 时有效
	VBool bool
	// 该列的数据，只在type 为DOUBLE 时有效
	VDouble float64
	// 该列的数据，只在type 为BINARY 时有效
	VBinary []byte
}

// 表的主键列值，精简数据模型
type OTSPrimaryKey DictString

// 表的属性列值，精简数据模型
type OTSAttribute DictString

// 表更新属性列值，精简数据模型
type OTSUpdateOfAttribute DictString

// 表的主键列值，复杂数据模型
type OTSPrimaryKeyColumns []OTSColumn

// 表的属性列值，复杂数据模型
type OTSAttributeColumns []OTSColumn

// 表更新属性列值，复杂数据模型
type OTSUpdateOfAttributeColumns []OTSColumn

// 在数据读取时，指定数据行中哪些属性列需要读取
type OTSColumnsToGet []string

// 在数据更新时，指定数据行中哪些属性列需要更新
type OTSColumnsToPut DictString

// 在数据更新时，指定数据行中哪些属性列需要删除
type OTSColumnsToDelete []string

// IsOk can be True or False
// when IsOk is False,
//     ErrorCode & ErrorMessage are available
// when IsOk is True,
//     Consumed & PrimaryKeyColumns & AttributeColumns are available
type OTSRowDataItem struct {
	IsOk              bool
	ErrorCode         int32
	ErrorMessage      string
	Consumed          OTSCapacityUnit
	PrimaryKeyColumns OTSPrimaryKeyColumns
	AttributeColumns  OTSAttributeColumns
}

// 创建行对象
type OTSPutRowItem struct {
	Condition        string
	PrimaryKey       OTSPrimaryKey
	AttributeColumns OTSAttribute
}

// type OTSPutRowItemRaw struct {
// 	Condition        string
// 	PrimaryKey       OTSPrimaryKeyColumns
// 	AttributeColumns OTSAttributeColumns
// }

// 更新行对象
type OTSUpdateRowItem struct {
	Condition                string
	PrimaryKey               OTSPrimaryKey
	UpdateOfAttributeColumns OTSUpdateOfAttribute
}

// type OTSUpdateRowItemRaw struct {
// 	Condition                string
// 	PrimaryKey               OTSPrimaryKeyColumns
// 	UpdateOfAttributeColumns OTSUpdateOfAttributeColumns
// }

// 删除行对象
type OTSDeleteRowItem struct {
	Condition  string
	PrimaryKey OTSPrimaryKey
}

// type OTSDeleteRowItemRaw struct {
// 	Condition  string
// 	PrimaryKey OTSPrimaryKeyColumns
// }

// 表的多行主键列值，精简数据模型
type OTSPrimaryKeyRows []OTSPrimaryKey

// 在BatchGetRow 操作中，表示要读取的一个表的请求信息
type OTSTableInBatchGetRowRequest struct {
	// 该表的表名
	TableName string
	// 该表中需要读取的全部行的信息
	Rows OTSPrimaryKeyRows
	// 该表中需要返回的全部列的列名
	ColumnsToGet OTSColumnsToGet
}

// 在BatchGetRow 操作中，表示要读取的多个表的请求信息
type OTSBatchGetRowRequest []OTSTableInBatchGetRowRequest

// 表的多行操作集合，精简数据模型
type OTSPutRows []OTSPutRowItem
type OTSUpdateRows []OTSUpdateRowItem
type OTSDeleteRows []OTSDeleteRowItem

// 在BatchWriteRow 操作中，表示要写入的一个表的请求信息
type OTSTableInBatchWriteRowRequest struct {
	// 该表的表名
	TableName string
	// 该表中需要写入的全部行的信息
	PutRows OTSPutRows
	// 该表中需要更新的全部行的信息
	UpdateRows OTSUpdateRows
	// 该表中需要删除的全部行的信息
	DeleteRows OTSDeleteRows
}

// 在BatchWriteRow 操作中，表示要写入的多个表的请求信息
type OTSBatchWriteRowRequest []OTSTableInBatchWriteRowRequest

// 在BatchWriteRow 操作中 服务器反馈对象
type OTSBatchWriteRowResponseItem struct {
	IsOk         bool
	ErrorCode    int32
	ErrorMessage string
	Consumed     OTSCapacityUnit
}

// OTS 表中的行按主键进行从小到大排序，GetRange 的读取范围是一个左闭右开的区间。操作
// 会返回主键属于该区间的行数据，区间的起始点是有效的主键或者是由INF_MIN 和INF_MAX
// 类型组成的虚拟点，虚拟点的列数必须与主键相同。其中，INF_MIN 表示无限小，任何类型的
// 值都比它大，INF_MAX 表示无限大，任何类型的值都比它小。

// only for GetRange
type OTS_INF_MIN struct {
}

// only for GetRange
type OTS_INF_MAX struct {
}

// 表示一个OTS实例下的表的列表
type OTSListTableResponse struct {
	TableNames []string
}

// 更新指定表的读服务能力单元或写服务能力单元设置，（新设定将于更新成功一分钟内生效）服务器响应
//
// tip:
//     调整每个表预留读写吞吐量的最小时间间隔为10 分钟，如果本次UpdateTable 操作距上次
//     不到10 分钟将被拒绝。
//     每个自然日(UTC 时间00:00:00 到第二天的00:00:00) 内每个表上调预留读写吞吐量次数不
//     限，但下调预留读写吞吐量次数不能超过4 次。下调写服务能力单元或者读服务能力单元其中
//     之一即视为下调预留读写吞吐量
type OTSUpdateTableResponse struct {
	// 更新后该表的预留读写吞吐量设置信息，除了包含当前的预留读写吞吐量设置值之外，还
	// 包含了最近一次更新该表的预留读写吞吐量设置的时间和当日已下调预留读写吞吐量的次数
	ReservedThroughputDetails *OTSReservedThroughputDetails
}

// 查询指定表的结构信息和预留读写吞吐量设置信息服务器响应
type OTSDescribeTableResponse struct {
	// 该表的Schema，与建表时给出的Schema 相同
	TableMeta *OTSTableMeta
	// 该表的预留读写吞吐量设置信息，除了包含当前的预留读写吞吐量设置值之外，还包含了
	// 最近一次更新该表的预留读写吞吐量设置的时间和当日已下调预留读写吞吐量的次数。
	ReservedThroughputDetails *OTSReservedThroughputDetails
}
