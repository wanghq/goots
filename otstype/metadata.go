// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// metadata for ots2
package otstype

import (
	"fmt"
	"time"
)

var OTSColumnType_INF_MIN OTS_INF_MIN // only for GetRange
var OTSColumnType_INF_MAX OTS_INF_MAX // only for GetRange
const (
	// OTSColumnType
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
type OTSSchemaOfPrimaryKey DictString

func (o OTSSchemaOfPrimaryKey) Del(key string) {
	DictString(o).Del(key)
}

func (o OTSSchemaOfPrimaryKey) Get(key string) interface{} {
	return DictString(o).Get(key)
}

func (o OTSSchemaOfPrimaryKey) Set(key string, value interface{}) {
	DictString(o).Set(key, value)
}

// 表示一次操作消耗服务能力单元的值或是一个表的预留读写吞吐量的值
type OTSCapacityUnit struct {
	// 本次操作消耗的读服务能力单元或该表的读服务能力单元
	Read int32
	// 本次操作消耗的写服务能力单元或该表的写服务能力单元
	Write int32
}

// 获取本次操作消耗的读服务能力单元或该表的读服务能力单元
func (o *OTSCapacityUnit) GetRead() int32 {
	return o.Read
}

// 获取本次操作消耗的写服务能力单元或该表的写服务能力单元
func (o *OTSCapacityUnit) GetWrite() int32 {
	return o.Write
}

// 表示一个表设置的预留读写吞吐量数值
type OTSReservedThroughput struct {
	// 表当前的预留读写吞吐量数值
	CapacityUnit OTSCapacityUnit
}

// 表示一个表的预留读写吞吐量信息
type OTSReservedThroughputDetails struct {
	// 该表的预留读写吞吐量的数值
	CapacityUnit *OTSCapacityUnit
	// 最近一次上调该表的预留读写吞吐量设置的时间，使用UTC 秒数表示
	LastIncreaseTime time.Time
	// 最近一次下调该表的预留读写吞吐量设置的时间，使用UTC 秒数表示
	LastDecreaseTime time.Time
	// 一个自然日内已下调该表的预留读写吞吐量设置的次数
	NumberOfDecreasesToday int32
}

// 表示一列，复杂数据模型
// type OTSColumn struct {
// 	// 该列的列名
// 	Name string
// 	// 该列的列值
// 	Value OTSColumnValue
// }

// 表示一列的列值，复杂数据模型
// type OTSColumnValue struct {
// 	// 该列的数据类型
// 	Type string
// 	// 该列的数据，只在type 为INTEGER 时有效
// 	VInt int64
// 	// 该列的数据，只在type 为STRING 时有效，必须为UTF-8 编码
// 	// golang默认为UTF-8编码
// 	VString string
// 	// 该列的数据，只在type 为BOOLEAN 时有效
// 	VBool bool
// 	// 该列的数据，只在type 为DOUBLE 时有效
// 	VDouble float64
// 	// 该列的数据，只在type 为BINARY 时有效
// 	VBinary []byte
// }

// 表的主键列值，精简数据模型
type OTSPrimaryKey DictString

func (o OTSPrimaryKey) String() string {
	r := ""
	if o == nil {
		return "None"
	}

	for k, v := range o {
		r = r + fmt.Sprintf("(%s:%v)", k, v)
	}

	return r
}

func (o OTSPrimaryKey) Del(key string) {
	DictString(o).Del(key)
}

func (o OTSPrimaryKey) Get(key string) interface{} {
	return DictString(o).Get(key)
}

func (o OTSPrimaryKey) Set(key string, value interface{}) {
	DictString(o).Set(key, value)
}

// 表的属性列值，精简数据模型
type OTSAttribute DictString

func (o OTSAttribute) String() string {
	r := ""
	if o == nil {
		return "None"
	}

	for k, v := range o {
		r = r + fmt.Sprintf("(%s:%v)", k, v)
	}

	return r
}

func (o OTSAttribute) Del(key string) {
	DictString(o).Del(key)
}

func (o OTSAttribute) Get(key string) interface{} {
	return DictString(o).Get(key)
}

func (o OTSAttribute) Set(key string, value interface{}) {
	DictString(o).Set(key, value)
}

// 表更新属性列值，精简数据模型
type OTSUpdateOfAttribute DictString

func (o OTSUpdateOfAttribute) String() string {
	r := ""
	if o == nil {
		return "None"
	}

	for k, v := range o {
		r = r + fmt.Sprintf("(%s:%v)", k, v)
	}

	return r
}

func (o OTSUpdateOfAttribute) Del(key string) {
	DictString(o).Del(key)
}

func (o OTSUpdateOfAttribute) Get(key string) interface{} {
	return DictString(o).Get(key)
}

func (o OTSUpdateOfAttribute) Set(key string, value interface{}) {
	DictString(o).Set(key, value)
}

// 表的主键列值，复杂数据模型
// type OTSPrimaryKeyColumns []OTSColumn

// 表的属性列值，复杂数据模型
// type OTSAttributeColumns []OTSColumn

// 表更新属性列值，复杂数据模型
// type OTSUpdateOfAttributeColumns []OTSColumn

// 在数据读取时，指定数据行中哪些属性列需要读取
type OTSColumnsToGet []string

// 在数据更新时，指定数据行中哪些属性列需要更新
type OTSColumnsToPut DictString

// 在数据更新时，指定数据行中哪些属性列需要删除
type OTSColumnsToDelete []string

//////////////////////////////////////////
/// Request
//////////////////////////////////////////

// 创建行对象
type OTSPutRowItem struct {
	Condition        string
	PrimaryKey       OTSPrimaryKey
	AttributeColumns OTSAttribute
}

// 更新行对象
type OTSUpdateRowItem struct {
	Condition                string
	PrimaryKey               OTSPrimaryKey
	UpdateOfAttributeColumns OTSUpdateOfAttribute
}

// 删除行对象
type OTSDeleteRowItem struct {
	Condition  string
	PrimaryKey OTSPrimaryKey
}

// 创建行对象，复杂数据模型
// type OTSPutRowItemRaw struct {
// 	Condition        string
// 	PrimaryKey       OTSPrimaryKeyColumns
// 	AttributeColumns OTSAttributeColumns
// }

// 更新行对象，复杂数据模型
// type OTSUpdateRowItemRaw struct {
// 	Condition                string
// 	PrimaryKey               OTSPrimaryKeyColumns
// 	UpdateOfAttributeColumns OTSUpdateOfAttributeColumns
// }

// 删除行对象，复杂数据模型
// type OTSDeleteRowItemRaw struct {
// 	Condition  string
// 	PrimaryKey OTSPrimaryKeyColumns
// }

// 表的多行主键列值，精简数据模型
type OTSPrimaryKeyRows []OTSPrimaryKey

// 在BatchGetRow 操作中，表示要读取的一个表的请求信息
type OTSTableInBatchGetRowRequestItem struct {
	// 该表的表名
	TableName string
	// 该表中需要读取的全部行的信息
	Rows OTSPrimaryKeyRows
	// 该表中需要返回的全部列的列名
	ColumnsToGet OTSColumnsToGet
}

// 在BatchGetRow 操作中，表示要读取的多个表的请求信息
type OTSBatchGetRowRequest []OTSTableInBatchGetRowRequestItem

// 表的多行操作集合，精简数据模型
type OTSPutRows []OTSPutRowItem
type OTSUpdateRows []OTSUpdateRowItem
type OTSDeleteRows []OTSDeleteRowItem

// 在BatchWriteRow 操作中，表示要写入的一个表的请求信息
type OTSTableInBatchWriteRowRequestItem struct {
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
type OTSBatchWriteRowRequest []OTSTableInBatchWriteRowRequestItem

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

//////////////////////////////////////////
/// Response
//////////////////////////////////////////

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

// 一行数据的主键列和属性列
type OTSRow struct {
	// 主键列
	PrimaryKeyColumns OTSPrimaryKey
	// 属性列
	AttributeColumns OTSAttribute
}

func (o *OTSRow) String() string {
	r := "PrimaryKeyColumns: " + o.GetPrimaryKeyColumns().String() + "\n"
	r = r + "AttributeColumns: " + o.GetAttributeColumns().String() + "\n"

	return r
}

func (o *OTSRow) GetPrimaryKeyColumns() OTSPrimaryKey {
	if o.PrimaryKeyColumns == nil {
		return nil
	} else {
		return o.PrimaryKeyColumns
	}
}

func (o *OTSRow) GetAttributeColumns() OTSAttribute {
	if o.AttributeColumns == nil {
		return nil
	} else {
		return o.AttributeColumns
	}
}

// 多行数据
type OTSRows []*OTSRow

// 获取一行数据
type OTSGetRowResponse struct {
	// 消耗的读服务能力单元或该表的读服务能力单元
	Consumed *OTSCapacityUnit
	// 行数据，包含了主键列和属性列
	Row *OTSRow
}

func (o *OTSGetRowResponse) GetReadConsumed() int32 {
	if o.Consumed != nil {
		return o.Consumed.GetRead()
	}

	return 0
}

func (o *OTSGetRowResponse) GetAttributeColumns() OTSAttribute {
	if o.Row != nil {
		return o.Row.GetAttributeColumns()
	}

	return nil
}

// 插入一行数据
type OTSPutRowResponse struct {
	// 消耗的读服务能力单元或该表的读服务能力单元
	Consumed *OTSCapacityUnit
}

func (o *OTSPutRowResponse) GetWriteConsumed() int32 {
	if o.Consumed != nil {
		return o.Consumed.GetWrite()
	}

	return 0
}

// 更新一行数据
type OTSUpdateRowResponse struct {
	// 消耗的读服务能力单元或该表的读服务能力单元
	Consumed *OTSCapacityUnit
}

func (o *OTSUpdateRowResponse) GetWriteConsumed() int32 {
	if o.Consumed != nil {
		return o.Consumed.GetWrite()
	}

	return 0
}

// 删除一行数据
type OTSDeleteRowResponse struct {
	// 消耗的读服务能力单元或该表的读服务能力单元
	Consumed *OTSCapacityUnit
}

func (o *OTSDeleteRowResponse) GetWriteConsumed() int32 {
	if o.Consumed != nil {
		return o.Consumed.GetWrite()
	}

	return 0
}

// 在BatchGetRow 操作的返回消息中，表示一行数据。
type OTSRowInBatchGetRowResponseItem struct {
	// 该行操作是否成功。若为true，则该行读取成功，error 无效；若为false，则该行读取失败，row 无效
	IsOk bool
	// 该行操作的错误号
	ErrorCode string
	// 该行操作的错误信息
	ErrorMessage string
	// 该行操作消耗的服务能力单元
	Consumed *OTSCapacityUnit
	// 行数据，包含了主键列和属性列
	Row *OTSRow
}

func (o *OTSRowInBatchGetRowResponseItem) GetErrorCode() string {
	return o.ErrorCode
}

func (o *OTSRowInBatchGetRowResponseItem) GetErrorMessage() string {
	return o.ErrorMessage
}

func (o *OTSRowInBatchGetRowResponseItem) GetReadConsumed() int32 {
	if o.Consumed != nil {
		return o.Consumed.GetRead()
	}

	return 0
}

func (o *OTSRowInBatchGetRowResponseItem) GetRow() *OTSRow {
	if o.Row != nil {
		return o.Row
	}

	return nil
}

// 在 BatchGetRow 操作的返回消息中，表示一个表的数据。
type OTSTableInBatchGetRowResponseItem struct {
	// 该表的表名
	TableName string
	// 该表中读取到的全部行数据
	Rows []*OTSRowInBatchGetRowResponseItem
}

func (o *OTSTableInBatchGetRowResponseItem) GetTableName() string {
	return o.TableName
}

func (o *OTSTableInBatchGetRowResponseItem) GetRows() []*OTSRowInBatchGetRowResponseItem {
	return o.Rows
}

// 对应了每个 table 下读取到的数据。
// 响应消息中 OTSTableInBatchGetRowResponseItem 对象的顺序与 OTSBatchGetRowRequest 中的
// OTSTableInBatchGetRowRequestItem 对象的顺序相同；每个 OTSTableInBatchGetRowResponseItem 下面的
// OTSRowInBatchGetRowResponseItem 对象的顺序与 OTSTableInBatchGetRowRequestItem 下面的 Rows 相同。
// 如果某行不存在或者某行在指定的 ColumnsToGet 下没有数据，仍然会在 OTSTableInBatchGetRowResponseItem
// 中有一条对应的 OTSRowInBatchGetRowResponseItem，但其 Row 下面的 PrimaryKeyColumns 和
// AttributeColumns 将为空。
//
// 若某行读取失败，该行所对应的 OTSRowInBatchGetRowResponseItem 中 IsOk 将为 false，此时 Row
// 将为空。
//
// 注意: BatchGetRow 操作可能会在行级别部分失败，此时返回的 HTTP 状态码仍为200。应用
// 程序必须对 OTSRowInBatchGetRowResponseItem 中的error 进行检查确认每一行的执行结果，并进行相
// 应的处理。
//
// 服务能力单元消耗:
// 如果本次操作整体失败，不消耗任何服务能力单元。
// 如果请求超时，结果未定义，服务能力单元有可能被消耗，也可能未被消耗。
// 其他情况将每个 OTSRowInBatchGetRowResponseItem 视为一个 OTSGetRow 操作独立计算读服务能力单
// 元。
type OTSBatchGetRowResponse struct {
	Tables []*OTSTableInBatchGetRowResponseItem
}

func (o *OTSBatchGetRowResponse) GetTables() []*OTSTableInBatchGetRowResponseItem {
	return o.Tables
}

// 在 BatchWriteRow 操作的返回消息中，表示一行写入操作的结果。
type OTSRowInBatchWriteRowResponseItem struct {
	// 该行操作是否成功。若为true，则该行写入成功，error 无效；若为false，则该行写入失败。
	IsOk bool
	// 该行操作的错误号
	ErrorCode string
	// 该行操作的错误信息
	ErrorMessage string
	// 该行操作消耗的服务能力单元
	Consumed *OTSCapacityUnit
}

func (o *OTSRowInBatchWriteRowResponseItem) GetErrorCode() string {
	return o.ErrorCode
}

func (o *OTSRowInBatchWriteRowResponseItem) GetErrorMessage() string {
	return o.ErrorMessage
}

func (o *OTSRowInBatchWriteRowResponseItem) GetWriteConsumed() int32 {
	if o.Consumed != nil {
		return o.Consumed.GetWrite()
	}

	return 0
}

// 在 BatchWriteRow 操作中，表示对一个表进行写入的结果。
type OTSTableInBatchWriteRowResponseItem struct {
	// 该表的表名
	TableName string
	// 该表中PutRow 操作的结果
	PutRows []*OTSRowInBatchWriteRowResponseItem
	// 该表中UpdateRow 操作的结果
	UpdateRows []*OTSRowInBatchWriteRowResponseItem
	// 该表中DeleteRow 操作的结果
	DeleteRows []*OTSRowInBatchWriteRowResponseItem
}

func (o *OTSTableInBatchWriteRowResponseItem) GetTableName() string {
	return o.TableName
}

func (o *OTSTableInBatchWriteRowResponseItem) GetPutRows() []*OTSRowInBatchWriteRowResponseItem {
	return o.PutRows
}

func (o *OTSTableInBatchWriteRowResponseItem) GetUpdateRows() []*OTSRowInBatchWriteRowResponseItem {
	return o.UpdateRows
}

func (o *OTSTableInBatchWriteRowResponseItem) GetDeleteRows() []*OTSRowInBatchWriteRowResponseItem {
	return o.DeleteRows
}

// 对应了每个 table 下各操作的响应信息，包括是否成功执行，错误码和消耗的服务能力单元。
// 响应消息中 OTSTableInBatchWriteRowResponseItem 对象的顺序与 OTSBatchWriteRowRequest 中的
// OTSTableInBatchWriteRowRequestItem 对象的顺序相同；每个 OTSTableInBatchWriteRowRequestItem
// 中 PutRows、UpdateRows、DeleteRows 包含的OTSRowInBatchWriteRowResponseItem 对象的顺序分别与
// OTSTableInBatchWriteRowRequestItem 中 PutRows、UpdateRows、DeleteRows 包含的 OTSPutRowItem，
// OTSUpdateRowItem 和 OTSDeleteRowItem 对象的顺序相同。
//
// 若某行读取失败，该行所对应的 OTSRowInBatchWriteRowResponseItem 中 IsOk 将为false。
//
// 注意:BatchWriteRow 操作可能会在行级别部分失败，此时返回的HTTP 状态码仍为200。应
// 用程序必须对 OTSRowInBatchWriteRowResponseItem 中的 error 进行检查，确认每一行的执行结果并进
// 行相应的处理。
//
// 服务能力单元消耗:
// 如果本次操作整体失败，不消耗任何服务能力单元。
// 如果请求超时，结果未定义，服务能力单元有可能被消耗，也可能未被消耗。
// 其他情况将每个 OTSPutRowItem、OTSUpdateRowItem、OTSDeleteRowItem
// 依次视作相对应的写操作独立计算读服务能力单元。
type OTSBatchWriteRowResponse struct {
	Tables []*OTSTableInBatchWriteRowResponseItem
}

func (o *OTSBatchWriteRowResponse) GetTables() []*OTSTableInBatchWriteRowResponseItem {
	return o.Tables
}

// 本次GetRange 操作的服务器反馈
//
// 服务能力单元消耗:
// GetRange 操作消耗读服务能力单元的数值为查询范围内所有行数据大小除以1KB 向上取
// 整。关于数据大小的计算请参见相关章节。
// 如果请求超时，结果未定义，服务能力单元有可能被消耗，也可能未被消耗。
// 如果返回内部错误（HTTP 状态码:5XX），则此次操作不消耗服务能力单元，其他错误情况
// 消耗1 读服务能力单元。
type OTSGetRangeResponse struct {
	// 该行操作消耗的服务能力单元
	Consumed *OTSCapacityUnit
	// 本次GetRange 操作的断点信息
	// 若为空，则本次GetRange 的响应消息中已包含了请求范围内的所有数据
	// 若不为空， 则表示本次GetRange 的响应消息中只包含了[inclusive_start_primary_key,
	// next_start_primary_key) 间的数据，若需要剩下的数据，需要将next_start_primary_key 作为inclusive_
	// start_primary_key，原始请求中的exclusive_end_primary_key 作为exclusive_end_primary_key
	// 继续执行GetRange 操作。
	// 注意:OTS 系统中限制了GetRange 操作的响应消息中数据不超过5000 行，大小不超过1M。
	// 即使在GetRange 请求中未设定limit，在响应中仍可能出现next_start_primary_key。因此在使用
	// GetRange 时一定要对响应中是否有next_start_primary_key 进行处理。
	NextStartPrimaryKey OTSPrimaryKey
	// 读取到的所有数据，若请求中direction 为FORWARD，则所有行按照主键由小到大进行排
	// 序；若请求中direction 为BACKWARD，则所有行按照主键由大到小进行排序
	// 其中每行的primary_key_columns 和attribute_columns 均只包含在columns_to_get 中指定的
	// 列，其顺序不保证与request 中的columns_to_get 一致；primary_key_columns 的顺序亦不保证与
	// 建表时指定的顺序一致。
	// 如果请求中指定的columns_to_get 不含有任何主键列，那么其主键在查询范围内，但没有任
	// 何一个属性列在columns_to_get 中的行将不会出现在响应消息里。
	Rows OTSRows
}

func (o *OTSGetRangeResponse) GetReadConsumed() int32 {
	if o.Consumed != nil {
		return o.Consumed.GetRead()
	}

	return 0
}

func (o *OTSGetRangeResponse) GetNextStartPrimaryKey() OTSPrimaryKey {
	if o.NextStartPrimaryKey != nil {
		return o.NextStartPrimaryKey
	}

	return nil
}

func (o *OTSGetRangeResponse) GetRows() OTSRows {
	if o.Rows != nil {
		return o.Rows
	}

	return nil
}
